package amqp

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/triasbrata/adios/pkgs/messagebroker/manager"
	"github.com/triasbrata/adios/pkgs/messagebroker/manager/connections"
	"github.com/triasbrata/adios/pkgs/messagebroker/publisher"
	"github.com/triasbrata/adios/pkgs/messagebroker/publisher/envelop"
	"github.com/triasbrata/adios/pkgs/utils"
	"golang.org/x/sync/errgroup"
)

type cPub struct {
	ctx         context.Context
	envelopOpt  envelop.Envelope
	deliveryTag string
	err         chan error
}
type ampqPub struct {
	conHolder   []manager.Manager[connections.ConnectionAMQP]
	middleware  []publisher.PublisherMiddleware
	midMut      sync.Mutex
	envelopChan chan cPub
	numWorker   int64
}

func (a *ampqPub) startCentralPub() {
	numPubWorker := slices.Max([]int64{slices.Min([]int64{2047, a.numWorker}), 10})
	eg, gctx := errgroup.WithContext(context.TODO())
	for _, conHolder := range a.conHolder {
		for range numPubWorker {
			eg.Go(func() error {
				for {
					select {
					case evlp, ok := <-a.envelopChan:
						if !ok {
							return nil
						}
						a.consume(evlp, conHolder.GetCon())
					case <-gctx.Done():
						return nil
					}
				}
			})
		}
	}
	eg.Wait()
}

func (a *ampqPub) consume(evlp cPub, con connections.ConnectionAMQP) {
	ch, err := con.Channel()
	if err != nil {
		evlp.err <- fmt.Errorf("failed when open channel %w", err)
	}
	defer ch.Close()
	evlp.err <- a.safePublish(evlp.ctx, ch, evlp.envelopOpt)
}

// Use implements publisher.Publisher.
func (a *ampqPub) Use(middleware publisher.PublisherMiddleware) {
	a.midMut.Lock()
	defer a.midMut.Unlock()
	a.middleware = append(a.middleware, middleware)
}

// Publish implements publisher.Publisher.
func (a *ampqPub) Publish(ctx context.Context, envelopOption envelop.EnvelopeOption) error {
	envelopOpt, err := a.buildEnvelop(ctx, envelopOption)
	if err != nil {
		return err
	}
	return a.pickupEnvelop(ctx, envelopOpt)
}

func (a *ampqPub) pickupEnvelop(ctx context.Context, envelopOpt envelop.Envelope) error {
	errChan := make(chan error, 1)
	defer close(errChan)
	a.envelopChan <- cPub{
		ctx:        ctx,
		envelopOpt: envelopOpt,
		err:        errChan,
	}
	return <-errChan
}

func (*ampqPub) safePublish(ctx context.Context, ch connections.ChannelAMQP, envelopOpt envelop.Envelope) error {
	if err := ch.Confirm(false); err != nil {
		return fmt.Errorf("got error when change to  confirm mode: %w", err)
	}
	returnMsg := ch.NotifyReturn(make(chan amqp091.Return, 1))
	confirm := ch.NotifyPublish(make(chan amqp091.Confirmation, 1))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, gctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		err := ch.PublishWithContext(gctx,
			envelopOpt.AMQP.Exchange.ExchangeName,
			envelopOpt.AMQP.Exchange.RoutingKey,
			envelopOpt.AMQP.Mandatory.Value(),
			false,
			envelopOpt.AMQP.Payload)
		if err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		//cancel publisher if already timeout
		defer cancel()
		select {
		case msg, ok := <-returnMsg:
			if ok {
				return fmt.Errorf("publish to %s failed, because returned", msg.RoutingKey)
			}
		case cfrm, ok := <-confirm:
			if ok && !cfrm.Ack {
				return fmt.Errorf("publish to %s failed, because nacked", envelopOpt.AMQP.Exchange.RoutingKey)
			}
		case <-time.After(envelopOpt.Timeout):
			return nil
		}
		return nil
	})
	return eg.Wait()
}

func (a *ampqPub) buildEnvelop(ctx context.Context, envelopOption envelop.EnvelopeOption) (envelop.Envelope, error) {
	var err error
	envelopOpt := envelopOption()
	for _, mfx := range a.middleware {
		err, envelopOpt = mfx(ctx, envelopOpt)
		if err != nil {
			return envelop.Envelope{}, fmt.Errorf("got error when build envelop from middleware: %w", err)
		}
	}
	return envelopOpt, nil
}

// PublishToQueue implements publisher.Publisher.
func (a *ampqPub) PublishToQueue(ctx context.Context, queueName string, Payload publisher.PublishPayload) error {
	return a.Publish(ctx, envelop.WithAMQPEnvelope(
		envelop.AMQPEnvelope{
			Exchange: envelop.AMQPEnvelopeExchange{
				RoutingKey:   queueName,
				ExchangeName: "",
			},
			Mandatory: utils.OptionBool{HasValue: true, Val: true},
			Payload: amqp091.Publishing{
				Headers: Payload.Header,
				Body:    Payload.Body,
			},
		},
	))
}

func NewPublisher(conHolder []manager.Manager[connections.ConnectionAMQP]) publisher.Publisher {
	inst := &ampqPub{
		conHolder:   conHolder,
		envelopChan: make(chan cPub, 1000),
	}
	go inst.startCentralPub()
	return inst
}

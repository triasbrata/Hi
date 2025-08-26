package amqp

import (
	"context"
	"fmt"
	"sync"

	"github.com/rabbitmq/amqp091-go"
	"github.com/triasbrata/adios/pkgs/messagebroker/manager"
	"github.com/triasbrata/adios/pkgs/messagebroker/manager/connections"
	"github.com/triasbrata/adios/pkgs/messagebroker/publisher"
	"github.com/triasbrata/adios/pkgs/messagebroker/publisher/envelop"
	"golang.org/x/sync/errgroup"
)

type ampqPub struct {
	conHolder  manager.Manager[connections.ConnectionAMQP]
	middleware []publisher.PublisherMiddleware
	midMut     sync.Mutex
}

// Use implements publisher.Publisher.
func (a *ampqPub) Use(middleware publisher.PublisherMiddleware) {
	a.midMut.Lock()
	defer a.midMut.Unlock()
	a.middleware = append(a.middleware, middleware)
}

// Publish implements publisher.Publisher.
func (a *ampqPub) Publish(ctx context.Context, envelopOption envelop.EnvelopeOption) error {
	ch, err := a.conHolder.GetCon().Channel()
	if err != nil {
		return fmt.Errorf("failed when open channel %w", err)
	}
	envelopOpt, err := a.buildEnvelop(ctx, envelopOption, err)
	if err != nil {
		return err
	}
	defer ch.Close()
	return a.safePublish(ctx, ch, err, envelopOpt)
}

func (*ampqPub) safePublish(ctx context.Context, ch connections.ChannelAMQP, err error, envelopOpt envelop.Envelope) error {
	if err := ch.Confirm(false); err != nil {
		return fmt.Errorf("got error when change to  confirm mode: %w", err)
	}
	returnMsg := ch.NotifyReturn(make(chan amqp091.Return, 1))
	confirm := ch.NotifyPublish(make(chan amqp091.Confirmation, 1))
	ctx, cancel := context.WithTimeout(ctx, envelopOpt.Timeout)
	defer cancel()
	eg, gctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		err = ch.PublishWithContext(gctx,
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
		select {
		case <-returnMsg:
			return fmt.Errorf("publish to %s failed, because returned", envelopOpt.AMQP.Exchange.RoutingKey)
		case cfrm := <-confirm:
			if !cfrm.Ack {
				return fmt.Errorf("publish to %s failed, because nacked", envelopOpt.AMQP.Exchange.RoutingKey)
			}
		case <-gctx.Done():
		}
		return nil
	})
	return eg.Wait()
}

func (a *ampqPub) buildEnvelop(ctx context.Context, envelopOption envelop.EnvelopeOption, err error) (envelop.Envelope, error) {
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
				ExchangeName: amqp091.ExchangeDirect,
			},
			Payload: amqp091.Publishing{
				Headers: Payload.Header,
				Body:    Payload.Body,
			},
		},
	))
}

func NewPublisher(conHolder manager.Manager[connections.ConnectionAMQP]) publisher.Publisher {
	return &ampqPub{conHolder: conHolder}
}

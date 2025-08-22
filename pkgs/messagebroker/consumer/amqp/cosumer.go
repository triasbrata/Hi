package amqp

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/triasbrata/adios/pkgs/messagebroker/consumer"
	"github.com/triasbrata/adios/pkgs/messagebroker/manager"
	"golang.org/x/sync/errgroup"
)

type amqpStack struct {
	queuName      string
	consumerName  string
	prefetchCount int64
	topology      consumer.AmqpTopologyConsumer
	handlers      []consumer.ConsumeHandler
}
type csmr struct {
	man              manager.Manager[amqp091.Connection]
	isStart          bool
	restartTime      time.Duration
	mut              sync.Mutex
	stack            []amqpStack
	globalMiddleware []consumer.ConsumeHandler
	ctxPool          sync.Pool
}

// Consume implements consumer.Consumer.
func (c *csmr) Consume(queueName string, topology consumer.ConsumerTopology, handlers ...consumer.ConsumeHandler) consumer.Consumer {
	if len(handlers) == 0 {
		slog.Warn("cant register queue", slog.String("queueName", queueName))
		return c
	}
	c.mut.Lock()
	c.stack = append(c.stack, amqpStack{
		queuName: queueName,
		topology: topology().Amqp,
		handlers: handlers,
	})
	c.mut.Unlock()
	return c
}

// SimpleConsume implements consumer.Consumer.
func (c *csmr) SimpleConsume(queueName string, handlers ...consumer.ConsumeHandler) consumer.Consumer {
	if len(handlers) == 0 {
		slog.Warn("cant register queue", slog.String("queueName", queueName))
		c.stack = append(c.stack, amqpStack{
			queuName: queueName,
			topology: consumer.AmqpTopologyConsumer{},
			handlers: handlers,
		})
		c.mut.Unlock()
		return c
	}
	return c
}

// Use implements consumer.Consumer.
func (c *csmr) Use(handlers ...consumer.ConsumeHandler) consumer.Consumer {
	if len(handlers) == 0 {
		slog.Warn("no global middleware was registered")
		return c
	}
	c.globalMiddleware = append(c.globalMiddleware, c.globalMiddleware...)
	return c
}
func (c *csmr) Start(ctx context.Context) error {
	eg, gctx := errgroup.WithContext(ctx)
	chanStack := make(chan amqpStack, len(c.stack))
	eg.Go(func() error {
		defer close(chanStack)
		for range c.man.Ready() {
			slog.InfoContext(ctx, "create consumer topology")
			for _, stack := range c.stack {
				chanStack <- stack
			}
		}
		return nil
	})
	for stack := range chanStack {
		eg.Go(func() error {
			err := c.buildTopology(gctx, stack)
			if err != nil && errors.Is(err, amqp091.ErrClosed) {
				slog.ErrorContext(gctx, "got error when define the topology", slog.Any("err", err))
				time.Sleep(c.restartTime)
				chanStack <- stack // restart topology and consumer
				return nil
			}
			err = c.startConsuming(gctx, stack)
			if err != nil && errors.Is(err, amqp091.ErrClosed) {
				slog.ErrorContext(gctx, "got error when define the topology", slog.Any("err", err))
				time.Sleep(c.restartTime)
				chanStack <- stack // restart topology and consumer
				return nil
			}

			return err
		})
	}

	return nil
}

func (c *csmr) startConsuming(gctx context.Context, stack amqpStack) error {
	ch, err := c.man.GetCon().Channel()
	if err != nil {
		return err
	}
	err = ch.Qos(int(stack.prefetchCount), 0, false)
	if err != nil {
		return err
	}
	del, err := ch.ConsumeWithContext(gctx,
		stack.queuName,
		stack.consumerName,
		stack.topology.AutoAck.Value(),
		stack.topology.Exclusive.Value(),
		stack.topology.NoLocal.Value(),
		stack.topology.NoWait.Value(),
		stack.topology.Args)
	if err != nil {
		return err
	}

	for msgDelivery := range del {
		attrs := []any{
			slog.String("msg", string(msgDelivery.Body)),
			slog.String("tag", msgDelivery.ConsumerTag),
		}
		go func() {
			ctx := c.ctxPool.Get().(*consumer.CtxConsumer)
			var errH error
			defer func() {
				c.ctxPool.Put(ctx)
				if errH != nil {
					slog.ErrorContext(gctx, "Failed when consume", append(attrs, slog.Any("err", errH))...)
				}
				if errH != nil && stack.topology.AutoAck.Value() && msgDelivery.Acknowledger != nil {
					errAck := msgDelivery.Acknowledger.Reject(msgDelivery.DeliveryTag, false)
					if errAck != nil {
						slog.ErrorContext(gctx, "Fail to reject", append(attrs, slog.Any("err", errAck))...)
					}
					return
				}
				if stack.topology.AutoAck.Value() && msgDelivery.Acknowledger != nil {
					errAck := msgDelivery.Acknowledger.Ack(msgDelivery.DeliveryTag, false)
					if errAck != nil {
						slog.ErrorContext(gctx, "Fail to reject", append(attrs, slog.Any("err", errAck))...)
					}
				}
			}()
			ctx.Context = gctx
			ctx.Body = msgDelivery.Body
			ctx.Header = msgDelivery.Headers
			for _, fx := range c.globalMiddleware {
				errH = fx(ctx)
				if errH != nil {
					return
				}
			}
			for _, fx := range stack.handlers {
				errH = fx(ctx)
				if errH != nil {
					return
				}
			}
		}()
	}
	return nil
}

func (c *csmr) buildTopology(gctx context.Context, stack amqpStack) (err error) {
	var ch *amqp091.Channel
	ch, err = c.man.GetCon().Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(stack.queuName,
		stack.topology.Durable.Value(),
		stack.topology.AutoDelete.Value(),
		stack.topology.Exclusive.Value(),
		stack.topology.NoWait.Value(),
		stack.topology.Args)
	if err != nil {
		return err
	}
	if stack.topology.BindExchange != nil {
		if stack.topology.BindExchange.Exchange != nil {
			err = ch.ExchangeDeclare(
				stack.topology.BindExchange.ExchangeName,
				stack.topology.BindExchange.Exchange.Kind,
				stack.topology.BindExchange.Exchange.Durable.Value(),
				stack.topology.BindExchange.Exchange.AutoDelete.Value(),
				stack.topology.BindExchange.Exchange.Internal.Value(),
				stack.topology.BindExchange.NoWait.Value(),
				stack.topology.BindExchange.Exchange.Args)
			if err != nil {
				return err
			}
		}
		err = ch.QueueBind(stack.queuName,
			stack.topology.BindExchange.RoutingKey,
			stack.topology.BindExchange.ExchangeName,
			stack.topology.BindExchange.NoWait.Value(),
			stack.topology.BindExchange.Args)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewConsumer(conManager manager.Manager[amqp091.Connection]) consumer.Consumer {
	return &csmr{
		man:              conManager,
		mut:              sync.Mutex{},
		globalMiddleware: make([]consumer.ConsumeHandler, 0),
		ctxPool: sync.Pool{
			New: func() any {
				return &consumer.CtxConsumer{}
			},
		},
	}
}

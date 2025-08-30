package bootstrap

import (
	"context"
	"fmt"

	"github.com/triasbrata/adios/internals/config"
	"github.com/triasbrata/adios/internals/delivery"
	"github.com/triasbrata/adios/pkgs/log"
	"github.com/triasbrata/adios/pkgs/messagebroker"
	"github.com/triasbrata/adios/pkgs/messagebroker/broker"
	"go.uber.org/fx"
)

type InvokeParam struct {
	fx.In
	Bk      broker.Broker
	Lc      fx.Lifecycle
	Conf    *config.Config
	Routing delivery.ConsumerRouting
}

func BootConsumerAmqp() fx.Option {
	return fx.Module("bootstrap/BootConsumerAmqp",
		log.LoadLoggerSlog(),
		config.LoadConfig(),
		messagebroker.LoadMessageBrokerAmqp(),
		delivery.ModuleConsumer(),
		fx.Invoke(func(param InvokeParam) {
			invCtx, cancel := context.WithCancel(context.Background())
			param.Lc.Append(fx.Hook{OnStart: func(ctx context.Context) error {
				consumer, err := param.Bk.Consumer(invCtx, broker.ConsumeWithAmqp(broker.AmqpConsumerConfig{
					RestartTime: param.Conf.Consumer.Amqp.RestartTime,
				}))

				if err != nil {
					return fmt.Errorf("error when try to consume %w", err)
				}
				param.Routing(consumer)
				go consumer.Start(invCtx)
				ok, errChan := consumer.Status()
				select {
				case err := <-errChan:
					return err
				case <-ctx.Done():
					cancel()
				case <-ok:
				}
				return nil
			}, OnStop: func(ctx context.Context) error {
				cancel()
				return nil
			}})
		}),
	)
}

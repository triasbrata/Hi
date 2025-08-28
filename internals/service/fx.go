package service

import (
	"context"

	"github.com/triasbrata/adios/internals/repositories"
	"github.com/triasbrata/adios/internals/service/hello/impl"
	"github.com/triasbrata/adios/pkgs/messagebroker/broker"
	"github.com/triasbrata/adios/pkgs/messagebroker/publisher"
	"go.uber.org/fx"
)

func LoadHelloService() fx.Option {
	return fx.Module("service/hello",
		fx.Provide(impl.NewServiceHello),
		repositories.LoadWordRepository(),
		fx.Provide(func(brk broker.Broker, lc fx.Lifecycle) (publisher.Publisher, error) {
			ctx, cancel := context.WithCancel(context.Background())
			// hook to lifecycle when app close then close the context
			lc.Append(fx.Hook{OnStop: func(ctx context.Context) error {
				cancel()
				return nil
			}})
			return brk.Publisher(ctx, broker.PublishWithAmqp())
		}),
	)
}

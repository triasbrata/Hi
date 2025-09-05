package service

import (
	"context"

	"github.com/triasbrata/adios/internals/repositories"
	"github.com/triasbrata/adios/internals/service/weather/impl"
	"github.com/triasbrata/adios/pkgs/messagebroker/broker"
	"github.com/triasbrata/adios/pkgs/messagebroker/publisher"
	"github.com/triasbrata/adios/pkgs/middleware"
	"go.uber.org/fx"
)

func LoadHelloService() fx.Option {
	return fx.Module("service/hello",
		fx.Provide(impl.NewServiceHello),
		repositories.LoadWeatherRepository(),
		fx.Provide(func(brk broker.Broker, lc fx.Lifecycle) (publisher.Publisher, error) {
			ctx, cancel := context.WithCancel(context.Background())
			// hook to lifecycle when app close then close the context
			lc.Append(fx.Hook{OnStop: func(ctx context.Context) error {
				cancel()
				return nil
			}})
			publisher, err := brk.Publisher(ctx, broker.PublishWithAmqp())
			if err != nil {
				return publisher, err

			}
			publisher.Use(middleware.OtelPublisherInject())
			return publisher, nil
		}),
	)
}

package delivery

import (
	"github.com/triasbrata/adios/internals/delivery/consumer"
	"github.com/triasbrata/adios/internals/delivery/consumer/impl"
	"github.com/triasbrata/adios/internals/delivery/http"
	implHttp "github.com/triasbrata/adios/internals/delivery/http/impl"
	"github.com/triasbrata/adios/internals/service"
	cmr "github.com/triasbrata/adios/pkgs/messagebroker/consumer"
	"go.uber.org/fx"
)

type ConsumerRouting func(builder cmr.ConsumerBuilder)

func ModuleHttp() fx.Option {
	return fx.Module("delivery/http",
		service.LoadHelloService(),
		fx.Provide(fx.Private, implHttp.NewHandler),
		fx.Invoke(http.NewRouter),
	)
}

func ModuleConsumer() fx.Option {
	return fx.Module("delivery/consumer",
		fx.Provide(impl.NewHandlerConsumer),
		fx.Provide(func(handler consumer.ConsumerHandler) ConsumerRouting {
			return func(builder cmr.ConsumerBuilder) {

				consumer.NewRoutingConsumer(handler, builder)
			}
		}))
}

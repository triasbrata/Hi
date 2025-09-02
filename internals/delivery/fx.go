package delivery

import (
	v1 "github.com/triasbrata/adios/gen/proto_go/weather/v1"
	"github.com/triasbrata/adios/internals/delivery/consumer"
	implConsumer "github.com/triasbrata/adios/internals/delivery/consumer/impl"
	hGrpc "github.com/triasbrata/adios/internals/delivery/grpc"
	"github.com/triasbrata/adios/internals/delivery/grpc/impl"
	"github.com/triasbrata/adios/internals/delivery/http"
	implHttp "github.com/triasbrata/adios/internals/delivery/http/impl"
	"github.com/triasbrata/adios/internals/service"
	cmr "github.com/triasbrata/adios/pkgs/messagebroker/consumer"
	serverGrpc "github.com/triasbrata/adios/pkgs/server/grpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type ConsumerRouting func(builder cmr.ConsumerBuilder)

func ModuleHttp() fx.Option {
	return fx.Module("delivery/http",
		service.LoadHelloService(),
		fx.Provide(fx.Private, implHttp.NewHandler),
		fx.Provide(http.NewRouter),
	)
}

func ModuleConsumer() fx.Option {
	return fx.Module("delivery/consumer",
		fx.Provide(implConsumer.NewHandlerConsumer),
		fx.Provide(func(handler consumer.ConsumerHandler) ConsumerRouting {
			return func(builder cmr.ConsumerBuilder) {

				consumer.NewRoutingConsumer(handler, builder)
			}
		}))
}

func ModuleGrpc() fx.Option {
	return fx.Module("delivery/grpc",
		fx.Provide(impl.NewWeatherHandler),
		fx.Provide(func(h hGrpc.Handler) serverGrpc.GrpcServerBinding {
			return serverGrpc.GrpcServerBinding(func(s *grpc.Server) {
				v1.RegisterWeatherServiceServer(s, h)
			})
		}),
	)
}

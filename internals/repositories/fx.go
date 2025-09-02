package repositories

import (
	"github.com/triasbrata/adios/internals/config"
	implWeather "github.com/triasbrata/adios/internals/repositories/weather/impl"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func LoadWeatherRepository() fx.Option {
	return fx.Module("repositories/weather",
		fx.Provide(fx.Private, func(cfg *config.Config, tp *trace.TracerProvider) (grpc.ClientConnInterface, error) {
			cc, err := grpc.NewClient(cfg.GrpcClientServices.WeatherService.Target,
				grpc.WithStatsHandler(otelgrpc.NewClientHandler(
					otelgrpc.WithTracerProvider(tp),
					otelgrpc.WithPropagators(otel.GetTextMapPropagator()))),
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, err
			}
			return cc, nil
		}), fx.Provide(implWeather.NewRepository))

}

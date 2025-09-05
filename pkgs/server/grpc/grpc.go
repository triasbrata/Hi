package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/grafana/pyroscope-go"
	"github.com/triasbrata/adios/internals/config"
	"github.com/triasbrata/adios/pkgs/instrumentation"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServerBinding func(s *grpc.Server)
type RoutingBind func() error

func LoadGrpcServer() fx.Option {
	return fx.Module("pkg/server/grpc",
		fx.Provide(func(cfg *config.Config, tProvider trace.TracerProvider, mProvider metric.MeterProvider) *grpc.Server {

			server := grpc.NewServer(grpc.StatsHandler(
				otelgrpc.NewServerHandler(
					otelgrpc.WithTracerProvider(tProvider),
					otelgrpc.WithMeterProvider(mProvider),
					otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
				)))
			if cfg.GrpcServer.EnableReflection {
				reflection.Register(server)
			}
			return server
		}),
		fx.Invoke(func(cfg *config.Config, lc fx.Lifecycle, server *grpc.Server, binder GrpcServerBinding, py *pyroscope.Profiler, tProvider trace.TracerProvider, mProvider metric.MeterProvider) error {
			otel.SetMeterProvider(mProvider)
			otel.SetTracerProvider(tProvider)
			instrumentation.SetTrace(tProvider.Tracer(cfg.AppName))
			lc.Append(fx.Hook{OnStart: func(ctx context.Context) error {
				address := fmt.Sprintf("%s:%s", cfg.GrpcServer.Address, cfg.GrpcServer.Port)
				lis, err := net.Listen("tcp", address)
				if err != nil {
					return fmt.Errorf("got error when listen to network: %w", err)
				}
				binder(server)
				go server.Serve(lis)
				return nil
			}, OnStop: func(ctx context.Context) error {
				return py.Stop()
			}})
			return nil
		}))
}

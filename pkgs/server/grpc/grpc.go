package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/triasbrata/adios/internals/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServerBinding func(s *grpc.Server)

func LoadGrpcServer() fx.Option {
	return fx.Module("pkg/server/grpc",
		fx.Provide(func(cfg *config.Config, tProvider *trace.TracerProvider, mProvider *metric.MeterProvider) *grpc.Server {
			fmt.Println("say hello")
			server := grpc.NewServer(grpc.StatsHandler(
				otelgrpc.NewServerHandler(
					otelgrpc.WithTracerProvider(tProvider),
					otelgrpc.WithMeterProvider(mProvider),
				)))
			if cfg.GrpcServer.EnableReflection {
				reflection.Register(server)
			}
			return server
		}), fx.Invoke(func(cfg *config.Config, lc fx.Lifecycle, server *grpc.Server, binder GrpcServerBinding) error {
			lc.Append(fx.Hook{OnStart: func(ctx context.Context) error {
				lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.GrpcServer.Address, cfg.GrpcServer.Port))
				if err != nil {
					return fmt.Errorf("got error when listen to network: %w", err)
				}
				binder(server)
				go server.Serve(lis)
				return nil
			}})
			return nil
		}))
}

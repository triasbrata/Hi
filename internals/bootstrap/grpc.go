package bootstrap

import (
	"github.com/triasbrata/adios/internals/config"
	"github.com/triasbrata/adios/internals/delivery"
	"github.com/triasbrata/adios/pkgs/instrumentation"
	"github.com/triasbrata/adios/pkgs/log"
	"github.com/triasbrata/adios/pkgs/secrets"
	pkgGrpc "github.com/triasbrata/adios/pkgs/server/grpc"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.31.0"
	"go.uber.org/fx"
)

func BootGRPC() fx.Option {
	return fx.Module("bootstrap/BootConsumerAmqp",
		log.LoadLoggerSlog(),
		config.LoadConfig(),
		instrumentation.OtelModule(
			func(sec secrets.Secret) instrumentation.InstrumentationIn {
				return instrumentation.InstrumentationInFunc(func() []attribute.KeyValue {
					return []attribute.KeyValue{
						semconv.VCSChangeID(sec.GetSecretAsString("GIT_COMMIT_ID", "1234")),
					}
				})
			},
			semconv.VCSRepositoryName("olla"),
		),
		pkgGrpc.LoadGrpcServer(),
		delivery.ModuleGrpc(),
	)
}

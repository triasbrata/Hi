package bootstrap

import (
	"github.com/triasbrata/adios/internals/config"
	"github.com/triasbrata/adios/internals/delivery"
	"github.com/triasbrata/adios/pkgs/instrumentation"
	plog "github.com/triasbrata/adios/pkgs/log"
	"github.com/triasbrata/adios/pkgs/messagebroker"
	"github.com/triasbrata/adios/pkgs/pyroscope"
	routersfx "github.com/triasbrata/adios/pkgs/routers/fx"
	"github.com/triasbrata/adios/pkgs/secrets"
	"github.com/triasbrata/adios/pkgs/server/http"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.31.0"
	"go.uber.org/fx"
)

func BootHttpServer() fx.Option {
	return fx.Options(plog.LoadLoggerSlog(),
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
		delivery.ModuleHttp(),
		routersfx.LoadModuleRouter(),
		messagebroker.LoadMessageBrokerAmqp(),
		pyroscope.LoadPyroscope(),
		http.LoadHttpServer())
}

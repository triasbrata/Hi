package http

import (
	"fmt"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/triasbrata/adios/pkgs/routers"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"

	"go.uber.org/fx"
)

type Param struct {
	fx.In
	Router    routers.Router
	TraceProv *trace.TracerProvider
	MeterProv *metric.MeterProvider
	Handler   Handler
}

func NewRouter(p Param) error {
	globalMiddleware(p)

	p.Router.Get("/hello-:word", p.Handler.HelloWorld)
	return nil
}

func globalMiddleware(p Param) {
	p.Router.GlobalMiddleware(otelfiber.Middleware(
		otelfiber.WithCollectClientIP(true),
		otelfiber.WithTracerProvider(p.TraceProv),
		otelfiber.WithMeterProvider(p.MeterProv),
		otelfiber.WithPropagators(propagation.NewCompositeTextMapPropagator()),
		otelfiber.WithSpanNameFormatter(func(ctx *fiber.Ctx) string {
			pattern := ctx.Route().Path
			if pattern == "" {
				pattern = ctx.OriginalURL()
			}

			return fmt.Sprintf("%s %s", ctx.Method(), pattern)
		}),
	))
}

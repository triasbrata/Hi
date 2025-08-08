package http

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/triasbrata/adios/internals/config"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

type InvokeParam struct {
	fx.In
	Lc        fx.Lifecycle
	App       *fiber.App
	Cfg       *config.Config
	TraceProv *trace.TracerProvider
	MeterProv *metric.MeterProvider
}
type NewFiberParam struct {
	fx.In
	Cfg       *config.Config
	TraceProv *trace.TracerProvider
	MeterProv *metric.MeterProvider
}

func NewFiberServer(p NewFiberParam) *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: func(v interface{}) ([]byte, error) {
			return sonic.Marshal(v)
		},
		JSONDecoder: func(data []byte, v interface{}) error {
			return sonic.Unmarshal(data, v)
		},
	})
	return app
}
func InvokeFiberServer(p InvokeParam) {
	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := p.App.Listen(fmt.Sprintf("%s:%s", p.Cfg.HttpServer.Address, p.Cfg.HttpServer.Port))
				if err != nil {
					slog.ErrorContext(ctx, "fiber server failed to start", slog.Any("err", err))
					os.Exit(1)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return p.App.ShutdownWithContext(ctx)
		},
	})
}
func LoadHttpServer() fx.Option {
	return fx.Module("bootstrap/http",
		fx.Provide(NewFiberServer),
		fx.Invoke(InvokeFiberServer),
	)
}

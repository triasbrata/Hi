package main

import (
	"log/slog"

	"github.com/triasbrata/adios/internals/bootstrap/http"
	"github.com/triasbrata/adios/internals/bootstrap/instrumentation"
	"github.com/triasbrata/adios/internals/bootstrap/log"
	"github.com/triasbrata/adios/internals/config"
	"github.com/triasbrata/adios/internals/delivery"
	routersfx "github.com/triasbrata/adios/pkgs/routers/fx"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	fx.New(
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{
				Logger: logger,
			}
		}),
		log.LoadLogger(),
		config.LoadConfig(),
		instrumentation.OtelModule(),
		delivery.ModuleHttp(),
		routersfx.LoadModuleRouter(),
		http.LoadHttpServer(),
	).Run()
}

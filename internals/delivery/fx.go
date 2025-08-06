package delivery

import (
	"github.com/triasbrata/adios/internals/delivery/http"
	"github.com/triasbrata/adios/internals/delivery/http/impl"
	"github.com/triasbrata/adios/internals/service"
	"go.uber.org/fx"
)

func ModuleHttp() fx.Option {
	return fx.Module("delivery/http",
		service.LoadHelloService(),
		fx.Provide(fx.Private, impl.NewHandler),
		fx.Invoke(http.NewRouter),
	)
}

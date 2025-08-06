package service

import (
	"github.com/triasbrata/adios/internals/repositories"
	"github.com/triasbrata/adios/internals/service/hello/impl"
	"go.uber.org/fx"
)

func LoadHelloService() fx.Option {
	return fx.Module("service/hello", fx.Provide(impl.NewServiceHello), repositories.LoadWordRepository())
}

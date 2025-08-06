package repositories

import (
	"github.com/triasbrata/adios/internals/repositories/words/impl"
	"go.uber.org/fx"
)

func LoadWordRepository() fx.Option {
	return fx.Module("repositories/words", fx.Provide(impl.NewWordRepository))
}

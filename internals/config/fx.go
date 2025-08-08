package config

import (
	"github.com/triasbrata/adios/pkgs/secrets/secretEnv"
	"go.uber.org/fx"
)

func LoadConfig() fx.Option {
	return fx.Module("config",
		fx.Provide(secretEnv.NewSecretFromEnv),
		fx.Provide(NewConfigEnv),
	)
}

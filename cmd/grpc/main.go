package main

import (
	"github.com/triasbrata/adios/internals/bootstrap"
	"go.uber.org/fx"
)

func main() {
	fx.New(bootstrap.BootGRPC()).Run()
}

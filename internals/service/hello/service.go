package hello

import (
	"context"

	"github.com/triasbrata/adios/internals/entities"
)

type HelloService interface {
	FetchHelloWorld(ctx context.Context, param entities.FetchHelloWorldParam) (res entities.FetchHelloWorldRes, err error)
}

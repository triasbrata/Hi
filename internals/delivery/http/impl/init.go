package impl

import (
	"github.com/triasbrata/adios/internals/config"
	httpDelivery "github.com/triasbrata/adios/internals/delivery/http"
	"github.com/triasbrata/adios/internals/service/hello"
)

type httpHandler struct {
	cfg     *config.Config
	service hello.HelloService
}

func NewHandler(cfg *config.Config, service hello.HelloService) (httpDelivery.Handler, error) {
	return &httpHandler{cfg: cfg, service: service}, nil
}

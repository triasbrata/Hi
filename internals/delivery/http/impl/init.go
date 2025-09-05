package impl

import (
	"github.com/triasbrata/adios/internals/config"
	httpDelivery "github.com/triasbrata/adios/internals/delivery/http"
	"github.com/triasbrata/adios/internals/service/weather"
)

type httpHandler struct {
	cfg     *config.Config
	service weather.WeatherService
}

func NewHandler(cfg *config.Config, service weather.WeatherService) (httpDelivery.Handler, error) {
	return &httpHandler{cfg: cfg, service: service}, nil
}

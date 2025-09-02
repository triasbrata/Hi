package impl

import (
	"github.com/grafana/pyroscope-go"
	"github.com/triasbrata/adios/internals/config"
	httpDelivery "github.com/triasbrata/adios/internals/delivery/http"
	"github.com/triasbrata/adios/internals/service/weather"
)

type httpHandler struct {
	cfg     *config.Config
	service weather.WeatherService
	py      *pyroscope.Profiler
}

func NewHandler(cfg *config.Config, service weather.WeatherService, py *pyroscope.Profiler) (httpDelivery.Handler, error) {
	return &httpHandler{cfg: cfg, service: service, py: py}, nil
}

package impl

import (
	"context"

	"github.com/bytedance/sonic"
	v1 "github.com/triasbrata/adios/gen/proto_go/weather/v1"
	"github.com/triasbrata/adios/internals/config"
	"github.com/triasbrata/adios/internals/entities"
	"github.com/triasbrata/adios/internals/repositories/weather"
	sWeather "github.com/triasbrata/adios/internals/service/weather"
	"github.com/triasbrata/adios/pkgs/instrumentation"
	"github.com/triasbrata/adios/pkgs/messagebroker/publisher"
)

type srv struct {
	cfg  *config.Config
	repo weather.WeatherServiceRepo
	pub  publisher.Publisher
}

// FetchCurrentWeather implements hello.HelloService.
func (s *srv) FetchCurrentWeather(ctx context.Context, param entities.FetchCurrentWeatherParam) (res entities.FetchCurrentWeatherRes, err error) {
	ctx, span := instrumentation.Tracer().Start(ctx, "internals:service:weather:impl:FetchCurrentWeather")
	defer span.End()
	currentWether, err := s.repo.GetWeather(ctx, &v1.GetWeatherRequest{
		Latitude:  param.Latitude,
		Longitude: param.Longitude,
	})
	if err != nil {
		return res, err
	}
	body, err := sonic.Marshal(currentWether)
	if err != nil {
		return res, err
	}
	err = s.pub.PublishToQueue(ctx, "latest_weather", publisher.PublishPayload{
		Body: body,
	})
	if err != nil {
		return res, err
	}
	res.Condition = currentWether.GetCondition().String()
	res.Temperature = currentWether.GetTemperature()
	return res, nil
}

func NewServiceHello(cfg *config.Config, repo weather.WeatherServiceRepo, pub publisher.Publisher) sWeather.WeatherService {
	return &srv{repo: repo, cfg: cfg, pub: pub}
}

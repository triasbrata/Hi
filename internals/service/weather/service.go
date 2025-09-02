package weather

import (
	"context"

	"github.com/triasbrata/adios/internals/entities"
)

type WeatherService interface {
	FetchCurrentWeather(ctx context.Context, param entities.FetchCurrentWeatherParam) (res entities.FetchCurrentWeatherRes, err error)
}

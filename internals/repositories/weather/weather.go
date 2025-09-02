package weather

import (
	"context"

	v1 "github.com/triasbrata/adios/gen/proto_go/weather/v1"
)

type WeatherServiceRepo interface {
	GetWeather(ctx context.Context, in *v1.GetWeatherRequest) (*v1.GetWeatherResponse, error)
}

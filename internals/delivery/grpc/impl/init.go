package impl

import (
	"context"
	"fmt"

	v1 "github.com/triasbrata/adios/gen/proto_go/weather/v1"
	"github.com/triasbrata/adios/internals/delivery/grpc"
)

type weatherHandler struct {
	v1.UnimplementedWeatherServiceServer
}

// GetWeather implements v1.WeatherServiceServer.
func (w *weatherHandler) GetWeather(context.Context, *v1.GetWeatherRequest) (*v1.GetWeatherResponse, error) {
	panic("unimplemented")
}

func NewWeatherHandler() grpc.Handler {
	fmt.Println("sayHello2")
	handler := &weatherHandler{}
	return handler
}

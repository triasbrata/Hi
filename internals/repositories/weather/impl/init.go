package impl

import (
	"context"
	"log"

	v1 "github.com/triasbrata/adios/gen/proto_go/weather/v1"
	"github.com/triasbrata/adios/internals/repositories/weather"
	"github.com/triasbrata/adios/pkgs/instrumentation"
	"google.golang.org/grpc"
)

type client struct {
	client v1.WeatherServiceClient
}

// GetWeather implements weather.WeatherServiceRepo.
func (c *client) GetWeather(ctx context.Context, in *v1.GetWeatherRequest) (*v1.GetWeatherResponse, error) {
	ctx, span := instrumentation.Tracer.Start(ctx, "internals:repositories:weather:impl:GetWeather")
	log.Printf("GetWeather trace=%s span=%s", span.SpanContext().TraceID().String(), span.SpanContext().SpanID())
	defer span.End()
	return c.client.GetWeather(ctx, in)
}

func NewRepository(cc grpc.ClientConnInterface) weather.WeatherServiceRepo {
	return &client{
		client: v1.NewWeatherServiceClient(cc),
	}
}

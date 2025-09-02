package impl

import (
	"context"
	"log"
	"math/rand"

	v1 "github.com/triasbrata/adios/gen/proto_go/weather/v1"
	"github.com/triasbrata/adios/internals/delivery/grpc"
	"github.com/triasbrata/adios/pkgs/instrumentation"
)

type weatherHandler struct {
	v1.UnimplementedWeatherServiceServer
}

// GetWeather implements v1.WeatherServiceServer.
func (w *weatherHandler) GetWeather(ctx context.Context, req *v1.GetWeatherRequest) (*v1.GetWeatherResponse, error) {
	ctx, span := instrumentation.Tracer.Start(ctx, "internals:delivery:grpc:impl:GetWeather")
	defer span.End()
	log.Printf("grpc trace=%s span=%s", span.SpanContext().TraceID().String(), span.SpanContext().SpanID())
	r := rand.Int31n(2)
	if r > 2 {
		r = 2
	}
	return &v1.GetWeatherResponse{
		Temperature: float32(rand.Int31n(100)) + rand.Float32(),
		Condition:   v1.Condition(r),
	}, nil
}

func NewWeatherHandler() grpc.Handler {
	handler := &weatherHandler{}
	return handler
}

package grpc

import v1 "github.com/triasbrata/adios/gen/proto_go/weather/v1"

type Handler interface {
	v1.WeatherServiceServer
}

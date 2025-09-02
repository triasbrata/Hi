package config

import "time"

type HttpServerConfig struct {
	Port    string
	Address string
}
type InstrumentationSecureConfig struct {
	CaFile []byte
}
type InstrumentationConfig struct {
	Secure       InstrumentationSecureConfig
	Endpoint     string
	UseGRPC      bool
	UseHttp      bool
	URL          string
	PyroscopeUrl string
}
type ConsumerConfig struct {
	Amqp AmqpConsumerConfig
}
type AmqpConsumerConfig struct {
	ConnectionName string
	URI            string
	RestartTime    time.Duration
}
type GrpcServerConfig struct {
	EnableReflection bool
	Port             string
	Address          string
}
type GrpcClientServicesConfig struct {
	WeatherService GrpcClientServiceConfig
}
type GrpcClientServiceConfig struct {
	Target string
}
type Config struct {
	AppName            string
	HttpServer         HttpServerConfig
	GrpcServer         GrpcServerConfig
	Instrumentation    InstrumentationConfig
	GrpcClientServices GrpcClientServicesConfig
	Consumer           ConsumerConfig
}

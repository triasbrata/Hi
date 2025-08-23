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
	Secure   InstrumentationSecureConfig
	Endpoint string
	UseGRPC  bool
	UseHttp  bool
	URL      string
}
type ConsumerConfig struct {
	Amqp AmqpConsumerConfig
}
type AmqpConsumerConfig struct {
	ConnectionName string
	URI            string
	RestartTime    time.Duration
}
type Config struct {
	AppName         string
	HttpServer      HttpServerConfig
	Instrumentation InstrumentationConfig
	Consumer        ConsumerConfig
}

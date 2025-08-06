package config

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
type Config struct {
	AppName         string
	HttpServer      HttpServerConfig
	Instrumentation InstrumentationConfig
}

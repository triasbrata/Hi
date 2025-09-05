package config

import (
	"fmt"
	"os"
	"time"

	"github.com/triasbrata/adios/pkgs/secrets"
)

func NewConfigEnv(secret secrets.Secret) (*Config, error) {
	cafile, err := os.ReadFile(secret.GetSecretAsString("INS_CA_PATH", ""))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("cant read ca file for instrumentation")
	}
	caStr := secret.GetSecretAsString("INS_CA_CONTENT", "")
	if caStr != "" {
		cafile = []byte(caStr)
	}
	secInst := InstrumentationSecureConfig{
		CaFile: cafile,
	}
	return &Config{
		HttpServer: HttpServerConfig{
			Port: secret.GetSecretAsString("HTTP_PORT", "8000"),
		},
		GrpcServer: GrpcServerConfig{
			EnableReflection: true,
			Port:             secret.GetSecretAsString("GRPC_PORT", "8001"),
		},
		AppName: secret.GetSecretAsString("APP_NAME", "hello"),
		Instrumentation: InstrumentationConfig{
			Secure:       secInst,
			Endpoint:     secret.GetSecretAsString("INS_ENDPOINT", "localhost:4317"),
			UseGRPC:      secret.GetSecretAsBool("INS_USE_GRPC", true),
			PyroscopeUrl: secret.GetSecretAsString("PYROSCOPE_SERVER_ADDRESS", "http://localhost:9999"),
		},
		GrpcClientServices: GrpcClientServicesConfig{
			WeatherService: GrpcClientServiceConfig{
				Target: secret.GetSecretAsString("GRPC_SERVICES_WEATHER_TARGET", "localhost:8001"),
			},
		},
		Consumer: ConsumerConfig{
			Amqp: AmqpConsumerConfig{
				ConnectionName: "consumer-test",
				URI:            "amqp://guest:guest@localhost:5672",
				RestartTime:    5 * time.Second,
			},
		},
	}, nil
}

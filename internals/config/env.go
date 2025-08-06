package config

import (
	"fmt"
	"os"

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
		AppName: secret.GetSecretAsString("APP_NAME", "hello"),
		Instrumentation: InstrumentationConfig{
			Secure:   secInst,
			Endpoint: secret.GetSecretAsString("INS_ENDPOINT", "localhost:4317"),
			UseGRPC:  secret.GetSecretAsBool("INS_USE_GRPC", true),
		},
	}, nil
}

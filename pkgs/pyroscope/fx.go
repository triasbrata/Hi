package pyroscope

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"

	py "github.com/grafana/pyroscope-go"
	"github.com/triasbrata/adios/internals/config"
	"go.uber.org/fx"
)

type pyLog struct {
}

func (p *pyLog) Debugf(msg string, args ...interface{}) {
	slog.Debug(fmt.Sprintf(msg, args...))
}
func (p *pyLog) Errorf(msg string, args ...interface{}) {
	slog.Error(fmt.Sprintf(msg, args...))
}
func (p *pyLog) Infof(msg string, args ...interface{}) {
	slog.Info(fmt.Sprintf(msg, args...))
}

func LoadPyroscope() fx.Option {
	return fx.Module("pkg/pyroscope", fx.Provide(func(cf *config.Config) (*py.Profiler, error) {
		runtime.SetMutexProfileFraction(5)
		runtime.SetBlockProfileRate(5)

		hostName, err := os.Hostname()
		if err != nil {
			return nil, err
		}
		p, err := py.Start(py.Config{
			ApplicationName: cf.AppName,
			ServerAddress:   cf.Instrumentation.PyroscopeUrl,
			Logger:          &pyLog{},
			Tags: map[string]string{
				"instance": hostName,
			},
			ProfileTypes: []py.ProfileType{
				py.ProfileCPU,
				py.ProfileInuseObjects,
				py.ProfileAllocObjects,
				py.ProfileInuseSpace,
				py.ProfileAllocSpace,
				py.ProfileGoroutines,
				py.ProfileMutexCount,
				py.ProfileMutexDuration,
				py.ProfileBlockCount,
				py.ProfileBlockDuration,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("error when init pyroscope : %w", err)
		}
		return p, nil
	}))
}

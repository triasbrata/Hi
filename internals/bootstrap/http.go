package bootstrap

import (
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/triasbrata/adios/internals/config"
	"github.com/triasbrata/adios/internals/delivery"
	"github.com/triasbrata/adios/pkgs/instrumentation"
	plog "github.com/triasbrata/adios/pkgs/log"
	"github.com/triasbrata/adios/pkgs/messagebroker"
	"github.com/triasbrata/adios/pkgs/pyroscope"
	routersfx "github.com/triasbrata/adios/pkgs/routers/fx"
	"github.com/triasbrata/adios/pkgs/secrets"
	"github.com/triasbrata/adios/pkgs/server/http"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.31.0"
	"go.uber.org/fx"
)

func BootHttpServer() fx.Option {
	return fx.Options(plog.LoadLoggerSlog(),
		config.LoadConfig(),
		instrumentation.OtelModule(
			func(sec secrets.Secret) instrumentation.InstrumentationIn {
				return instrumentation.InstrumentationInFunc(func() []attribute.KeyValue {
					return []attribute.KeyValue{
						semconv.VCSChangeID(sec.GetSecretAsString("GIT_COMMIT_ID", "1234")),
					}
				})
			},
			semconv.VCSRepositoryName("olla"),
		),
		delivery.ModuleHttp(),
		routersfx.LoadModuleRouter(),
		messagebroker.LoadMessageBrokerAmqp(),
		pyroscope.LoadPyroscope(),
		http.LoadHttpServer())
}

func dumpAllProfiles(outDir string) any {
	// Names available via runtime/pprof.Lookup
	// Common: goroutine, heap, allocs, threadcreate, block, mutex
	profiles := []string{
		"heap",
		"goroutine",
		"allocs",
		"threadcreate",
		"block",
		"mutex",
	}

	// Heap: WriteHeapProfile is equivalent to Lookup("heap").WriteTo(f, 0) after a GC.
	if err := writeHeap(filepath.Join(outDir, "heap.pprof")); err != nil {
		return err
	}

	for _, name := range profiles {
		// heap already written above
		if name == "heap" {
			continue
		}
		if err := writeProfile(name, filepath.Join(outDir, name+".pprof")); err != nil {
			log.Printf("write %s: %v", name, err)
		}
	}
	return nil
}

func writeHeap(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := pprof.WriteHeapProfile(f); err != nil {
		return err
	}
	log.Printf("wrote %s", path)
	return nil
}

func writeProfile(name, path string) error {
	p := pprof.Lookup(name)
	if p == nil {
		return nil // profile not available; skip
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	// debug=0 -> compressed proto for `go tool pprof`
	// (debug=1/2 are human-readable text; not what you want here)
	if err := p.WriteTo(f, 0); err != nil {
		return err
	}
	log.Printf("wrote %s", path)
	return nil
}

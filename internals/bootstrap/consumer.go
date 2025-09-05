package bootstrap

import (
	"context"
	"fmt"

	"github.com/grafana/pyroscope-go"
	"github.com/triasbrata/adios/internals/config"
	"github.com/triasbrata/adios/internals/delivery"
	"github.com/triasbrata/adios/pkgs/instrumentation"
	"github.com/triasbrata/adios/pkgs/log"
	"github.com/triasbrata/adios/pkgs/messagebroker"
	"github.com/triasbrata/adios/pkgs/messagebroker/broker"
	pyroscopePkg "github.com/triasbrata/adios/pkgs/pyroscope"
	"github.com/triasbrata/adios/pkgs/secrets"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.31.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

type InvokeParam struct {
	fx.In
	Bk        broker.Broker
	Lc        fx.Lifecycle
	Conf      *config.Config
	Routing   delivery.ConsumerRouting
	Tp        trace.TracerProvider
	Mp        metric.MeterProvider
	Pyroscope *pyroscope.Profiler
}

func BootConsumerAmqp() fx.Option {
	return fx.Options(
		pyroscopePkg.LoadPyroscope(),
		log.LoadLoggerSlog(),
		config.LoadConfig(),
		messagebroker.LoadMessageBrokerAmqp(),
		delivery.ModuleConsumer(),
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
		fx.Invoke(func(param InvokeParam) {
			otel.SetMeterProvider(param.Mp)
			otel.SetTracerProvider(param.Tp)
			instrumentation.SetTrace(param.Tp.Tracer(param.Conf.AppName))
			invCtx, cancel := context.WithCancel(context.Background())
			param.Lc.Append(fx.Hook{OnStart: func(ctx context.Context) error {
				consumer, err := param.Bk.Consumer(invCtx, broker.ConsumeWithAmqp(broker.AmqpConsumerConfig{
					RestartTime: param.Conf.Consumer.Amqp.RestartTime,
				}, broker.WithOtel(param.Tp, param.Mp)))

				if err != nil {
					return fmt.Errorf("error when try to consume %w", err)
				}
				param.Routing(consumer)
				go consumer.Start(invCtx)
				ok, errChan := consumer.Status()
				select {
				case err := <-errChan:
					return err
				case <-ctx.Done():
					cancel()
				case <-ok:
				}
				return nil
			}, OnStop: func(ctx context.Context) error {
				return param.Pyroscope.Stop()
				cancel()
				return nil
			}})
		}),
	)
}

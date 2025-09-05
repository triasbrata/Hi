package consumer

import (
	cmr "github.com/triasbrata/adios/pkgs/messagebroker/consumer"
	"github.com/triasbrata/adios/pkgs/middleware"
)

func NewRoutingConsumer(chandler ConsumerHandler, builder cmr.ConsumerBuilder) {
	builder.Use(middleware.OtelConsumerExtract())
	builder.Consume("latest_weather", func() cmr.TopologyConsumer {
		return cmr.TopologyConsumer{
			Amqp: cmr.AmqpTopologyConsumer{PrefetchCount: 200},
		}
	}, chandler.HandleTestConsumer)
}

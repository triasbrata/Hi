package consumer

import (
	cmr "github.com/triasbrata/adios/pkgs/messagebroker/consumer"
)

func NewRoutingConsumer(chandler ConsumerHandler, builder cmr.ConsumerBuilder) {
	builder.Consume("test_consumer", func() cmr.TopologyConsumer {
		return cmr.TopologyConsumer{
			Amqp: cmr.AmqpTopologyConsumer{PrefetchCount: 200},
		}
	}, chandler.HandleTestConsumer)
}

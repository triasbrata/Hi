package broker

import (
	"context"

	"github.com/triasbrata/adios/pkgs/messagebroker/consumer"
	"github.com/triasbrata/adios/pkgs/messagebroker/publisher"
)

type ConsumerConfig struct {
	Amqp bool
}

func WithAmqp() ConBuilder {
	return func() ConsumerConfig {
		return ConsumerConfig{Amqp: true}
	}
}

type ConBuilder func() ConsumerConfig
type PubBuilder func() ConsumerConfig

type Broker interface {
	Publisher(ctx context.Context, builder PubBuilder) (publisher.Publisher, error)
	Consumer(ctx context.Context, builder ConBuilder) (consumer.Consumer, error)
}

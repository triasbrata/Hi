package broker

import (
	"context"
	"time"

	"github.com/triasbrata/adios/pkgs/messagebroker/consumer"
	"github.com/triasbrata/adios/pkgs/messagebroker/publisher"
)

type ConsumerConfig struct {
	Amqp bool
	AmqpConConfig
}
type AmqpConConfig struct {
	RestartTime time.Duration
}

func WithAmqp(config AmqpConConfig) ConBuilder {
	return func() ConsumerConfig {
		return ConsumerConfig{Amqp: true, AmqpConConfig: config}
	}
}

type ConBuilder func() ConsumerConfig
type PubBuilder func() ConsumerConfig

type Broker interface {
	Publisher(ctx context.Context, builder PubBuilder) (publisher.Publisher, error)
	Consumer(ctx context.Context, builder ConBuilder) (consumer.Consumer, error)
}

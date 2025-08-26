package broker

import (
	"context"
	"time"

	"github.com/triasbrata/adios/pkgs/messagebroker/consumer"
	"github.com/triasbrata/adios/pkgs/messagebroker/publisher"
)

type BrokerDestination struct {
	Amqp bool
	AmqpConsumerConfig
}
type AmqpConsumerConfig struct {
	RestartTime time.Duration
}

func ConsumeWithAmqp(config AmqpConsumerConfig) ConBuilder {
	return func() BrokerDestination {
		return BrokerDestination{Amqp: true, AmqpConsumerConfig: config}
	}
}
func PublishWithAmqp() PubBuilder {
	return func() BrokerDestination {
		return BrokerDestination{Amqp: true}
	}
}

type ConBuilder func() BrokerDestination
type PubBuilder func() BrokerDestination

type Broker interface {
	Publisher(ctx context.Context, destination PubBuilder) (publisher.Publisher, error)
	Consumer(ctx context.Context, builder ConBuilder) (consumer.Consumer, error)
}

package consumer

import (
	"context"
)

type CtxConsumer struct {
	Body   []byte
	Header map[string]interface{}
	context.Context
}
type TopologyConsumer struct {
	Amqp AmqpTopologyConsumer
}
type ConsumerTopology func() TopologyConsumer
type ConsumeHandler func(c *CtxConsumer) error
type ConsumerBuilder interface {
	Consume(queueName string, topology ConsumerTopology, handlers ...ConsumeHandler) ConsumerBuilder
	SimpleConsume(queueName string, handlers ...ConsumeHandler) ConsumerBuilder
	Use(handlers ...ConsumeHandler) ConsumerBuilder
}
type Consumer interface {
	ConsumerBuilder
	Start(ctx context.Context) error
	Status() (ok chan struct{}, err chan error)
}

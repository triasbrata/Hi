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
type Consumer interface {
	Consume(queueName string, topology ConsumerTopology, handlers ...ConsumeHandler) Consumer
	SimpleConsume(queueName string, handlers ...ConsumeHandler) Consumer
	Use(handlers ...ConsumeHandler) Consumer
	Start(ctx context.Context) error
}

package amqp

import (
	"context"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

type contextAmqp struct {
	ctx             context.Context
	body            []byte
	header          map[string]interface{}
	stack           *amqpStack
	ciStackHandler  int64
	lenStackHandler int64
}

// SetUserContext implements consumer.CtxConsumer.
func (c *contextAmqp) SetUserContext(ctx context.Context) context.Context {
	c.ctx = ctx
	return c.ctx
}

// UpdateBody implements consumer.CtxConsumer.
func (c *contextAmqp) UpdateBody(body []byte) {
	c.body = body
}

// UpdateHeader implements consumer.CtxConsumer.
func (c *contextAmqp) UpdateHeader(key string, value any) {
	c.header[key] = value
}

// UserContext implements consumer.CtxConsumer.
func (c *contextAmqp) UserContext() context.Context {
	return c.ctx
}

func (c *contextAmqp) populateContext(ctx context.Context, msgDelivery amqp091.Delivery, stack amqpStack) {
	header := make(map[string]any)
	if len(msgDelivery.Headers) == 0 {
		header = msgDelivery.Headers
	}
	c.body = msgDelivery.Body
	c.header = header
	c.stack = &stack
	c.ctx = ctx
	c.ciStackHandler = 0
	c.lenStackHandler = int64(len(stack.handlers))
}

// Body implements consumer.CtxConsumer.
func (c *contextAmqp) Body() []byte {
	return c.body
}

// Header implements consumer.CtxConsumer.
func (c *contextAmqp) Header() map[string]interface{} {
	return c.header
}

// Next implements consumer.CtxConsumer.
func (c *contextAmqp) Next() error {
	defer func() {
		c.ciStackHandler++
	}()
	if c.ciStackHandler < c.lenStackHandler {
		return c.stack.handlers[c.ciStackHandler](c)
	}
	return nil
}

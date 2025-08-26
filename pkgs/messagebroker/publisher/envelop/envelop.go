package envelop

import (
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/triasbrata/adios/pkgs/utils"
)

type EnvelopeOption func() Envelope
type Envelope struct {
	AMQP    AMQPEnvelope
	Timeout time.Duration
}
type AMQPEnvelope struct {
	Payload   amqp091.Publishing
	Exchange  AMQPEnvelopeExchange
	Mandatory utils.OptionBool
}
type AMQPEnvelopeExchange struct {
	RoutingKey   string
	ExchangeName string
}

func WithAMQPEnvelope(envelope AMQPEnvelope, timeout ...time.Duration) EnvelopeOption {
	return func() Envelope {
		tOut := time.Second
		if len(timeout) == 1 {
			tOut = timeout[0]
		}
		return Envelope{
			AMQP:    envelope,
			Timeout: tOut,
		}
	}
}

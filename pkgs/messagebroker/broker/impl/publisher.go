package impl

import (
	"context"
	"fmt"

	"github.com/triasbrata/adios/pkgs/messagebroker/broker"
	"github.com/triasbrata/adios/pkgs/messagebroker/publisher/amqp"

	"github.com/triasbrata/adios/pkgs/messagebroker/publisher"
)

// Publisher implements broker.Broker.
func (b *brk) Publisher(ctx context.Context, builder broker.PubBuilder) (publisher.Publisher, error) {
	config := builder()
	switch {
	case config.Amqp:
		if b.cfg.amqp == nil {
			return nil, fmt.Errorf("configuration amqp was not found")
		}
		conHolder, err := b.openConnectionAmqp(ctx)
		if err != nil {
			return nil, err
		}
		return amqp.NewPublisher(conHolder), nil
	}
	return nil, fmt.Errorf("Consumer cant open")
}

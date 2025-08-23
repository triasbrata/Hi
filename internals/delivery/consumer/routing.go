package consumer

import (
	"fmt"

	cmr "github.com/triasbrata/adios/pkgs/messagebroker/consumer"
)

func NewRoutingConsumer(chandler ConsumerHandler, builder cmr.ConsumerBuilder) {
	builder.SimpleConsume("test_consumer", chandler.HandleTestConsumer)
	fmt.Printf("builder: %v\n", builder)
}

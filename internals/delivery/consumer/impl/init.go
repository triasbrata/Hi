package impl

import (
	"fmt"

	cmr "github.com/triasbrata/adios/internals/delivery/consumer"
	"github.com/triasbrata/adios/pkgs/messagebroker/consumer"
)

type handler struct {
}

// HandleTestConsumer implements consumer.ConsumerHandler.
func (h *handler) HandleTestConsumer(c *consumer.CtxConsumer) error {
	fmt.Printf("c.Body: %s\n", c.Body)
	return nil
}

func NewHandlerConsumer() cmr.ConsumerHandler {
	return &handler{}
}

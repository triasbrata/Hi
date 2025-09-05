package impl

import (
	"fmt"
	"sync/atomic"
	"time"

	cmr "github.com/triasbrata/adios/internals/delivery/consumer"
	"github.com/triasbrata/adios/pkgs/instrumentation"
	"github.com/triasbrata/adios/pkgs/messagebroker/consumer"
)

type handler struct {
	at *atomic.Int64
	st time.Time
}

// HandleTestConsumer implements consumer.ConsumerHandler.
func (h *handler) HandleTestConsumer(c consumer.CtxConsumer) error {
	ctx, span := instrumentation.Tracer().Start(c.UserContext(), "internals:delivery:consumer:impl:HandleTestConsumer")
	defer span.End()
	fmt.Printf("c.Body: %s\n", c.Body())
	h.at.Add(1)
	fmt.Printf("h.at.Load(): %v %v\n", h.at.Load(), time.Since(h.st))
	return ctx.Err()
}

func NewHandlerConsumer() cmr.ConsumerHandler {
	return &handler{at: &atomic.Int64{}, st: time.Now()}
}

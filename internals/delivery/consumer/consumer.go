package consumer

import "github.com/triasbrata/adios/pkgs/messagebroker/consumer"

type ConsumerHandler interface {
	HandleTestConsumer(c consumer.CtxConsumer) error
}

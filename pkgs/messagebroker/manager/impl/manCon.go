package impl

import (
	"sync"

	"github.com/rabbitmq/amqp091-go"
	"github.com/triasbrata/adios/pkgs/messagebroker/manager"
)

type manCon[T any] struct {
	mutex *sync.RWMutex
	con   *T
}

// GetConAMQP implements manager.Manager.
func (m *manCon[T]) GetCon() *T {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.con
}

// SetConAMQP implements manager.Manager.
func (m *manCon[T]) SetCon(con *T) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.con = con
}

func NewManager() manager.Manager[amqp091.Connection] {
	return &manCon[amqp091.Connection]{mutex: &sync.RWMutex{}}
}

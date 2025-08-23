package manager

type Manager[T any] interface {
	SetCon(con T)
	GetCon() T
	Ready() <-chan struct{}
	Release() error
}

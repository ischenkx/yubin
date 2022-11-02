package channel

type Handle[T any] interface {
	Chan() <-chan T
	Close()
}

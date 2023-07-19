package pubsub

import "sync/atomic"

type Handle[T any] struct {
	id      int64
	closed  int32
	pubsub  *PubSub[T]
	channel chan T
}

func newHandle[T any](pubsub *PubSub[T], id int64, size int) *Handle[T] {
	return &Handle[T]{
		id:      id,
		closed:  0,
		pubsub:  pubsub,
		channel: make(chan T, size),
	}
}

func (h *Handle[T]) Chan() <-chan T {
	return h.channel
}

func (h *Handle[T]) Close() {
	if !atomic.CompareAndSwapInt32(&h.closed, 0, 1) {
		return
	}
	close(h.channel)
	h.pubsub.remove(h.id)
}

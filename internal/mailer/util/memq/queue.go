package memq

import (
	"context"
	"smtp-client/pkg/channel"
)

type Queue[T any] chan T

func New[T any](capacity int) Queue[T] {
	return make(chan T, capacity)
}

func (q Queue[T]) Run(_ context.Context) {}

func (q Queue[T]) Push(t T) error {
	q <- t
	return nil
}

func (q Queue[T]) Handle() channel.Handle[T] {
	return newHandle(q)
}

type Handle[T any] struct {
	sender chan T
	closer chan struct{}
}

func newHandle[T any](src <-chan T) Handle[T] {
	sender := make(chan T, 128)
	closer := make(chan struct{})

	go func() {
		select {
		case data := <-src:
			sender <- data
		case <-closer:
			return
		}
	}()

	return Handle[T]{
		sender: sender,
		closer: closer,
	}
}

func (h Handle[T]) Chan() <-chan T {
	return h.sender
}

func (h Handle[T]) Close() {
	close(h.closer)
	close(h.sender)
}

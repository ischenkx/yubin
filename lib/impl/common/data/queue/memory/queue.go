package memq

import (
	"context"
	"yubin/common/data/queue"
)

var _ queue.Queue[int] = (Queue[int])(nil)

type Queue[T any] chan T

func New[T any](capacity int) Queue[T] {
	return make(chan T, capacity)
}

func (q Queue[T]) Push(_ context.Context, t T) error {
	q <- t
	return nil
}

func (q Queue[T]) Read(ctx context.Context) (<-chan T, error) {
	channel := make(chan T)
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case item := <-q:
				channel <- item
			}
		}
	}()

	return channel, nil
}

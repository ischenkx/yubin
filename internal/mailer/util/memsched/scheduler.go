package memsched

import (
	"context"
	"smtp-client/pkg/channel"
	"smtp-client/pkg/pubsub"
	"sync"
	"time"
)

type Scheduler[T any] struct {
	pubsub *pubsub.PubSub[T]
	mu     sync.RWMutex
}

func New[T any]() *Scheduler[T] {
	return &Scheduler[T]{
		pubsub: pubsub.New[T](),
	}
}

func (s *Scheduler[T]) Run(ctx context.Context) {}

func (s *Scheduler[T]) Schedule(data T, timestamp time.Time) error {
	go time.AfterFunc(timestamp.Sub(time.Now()), func() {
		s.pubsub.Publish(data)
	})
	return nil
}

func (s *Scheduler[T]) Handle() channel.Handle[T] {
	return s.pubsub.Subscribe()
}

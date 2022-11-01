package memsched

import (
	"context"
	"log"
	"smtp-client/internal/mailer"
	"time"
)

type TimeStamp[T any] struct {
	ID   string
	Data T
	Time time.Time
}

type Scheduler[T any] struct{}

func (s *Scheduler[T]) Run(ctx context.Context) {
	log.Println("mock scheduler")
}

func (s *Scheduler[T]) Schedule(data T, timestamp time.Time) error {
	return nil
}

func (s *Scheduler[T]) Handle() mailer.Handle[T] {
	return nil
}

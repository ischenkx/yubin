package mailer

import (
	"context"
	"time"
)

type Scheduler[T any] interface {
	Run(ctx context.Context)
	Schedule(T, time.Time) error
	Handle() Handle[T]
}

package mailer

import (
	"context"
	"smtp-client/pkg/channel"
	"time"
)

type Scheduler[T any] interface {
	Run(ctx context.Context)
	Schedule(T, time.Time) error
	Handle() channel.Handle[T]
}

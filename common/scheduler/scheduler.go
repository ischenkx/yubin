package scheduler

import (
	"context"
	"time"
)

type Scheduler[Item any] interface {
	Schedule(ctx context.Context, item Item, at time.Time) error
	Ready(ctx context.Context) <-chan Item
}

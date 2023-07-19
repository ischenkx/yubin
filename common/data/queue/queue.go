package queue

import (
	"context"
)

type Queue[T any] interface {
	Read(ctx context.Context) (<-chan T, error)
	Push(ctx context.Context, item T) error
}

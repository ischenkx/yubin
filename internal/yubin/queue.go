package yubin

import (
	"context"
	"smtp-client/pkg/channel"
)

type Queue[T any] interface {
	Run(ctx context.Context)
	Push(T) error
	Handle() channel.Handle[T]
}

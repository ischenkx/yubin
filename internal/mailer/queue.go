package mailer

import "context"

type Queue[T any] interface {
	Run(ctx context.Context)
	Push(T) error
	Handle() Handle[T]
}

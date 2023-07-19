package kv

import "context"

type Getter[V any] interface {
	Get(ctx context.Context, key string) (V, error)
	Range(ctx context.Context, order Order, offset int, limit int) ([]V, error)
}

type Setter[V any] interface {
	Set(ctx context.Context, key string, value V) error
}

type Deleter interface {
	Delete(ctx context.Context, key string) error
}

type Storage[V any] interface {
	Getter[V]
	Setter[V]
	Deleter
}

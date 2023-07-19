package redikv

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"yubin/common/codec"
	"yubin/common/data"
	"yubin/common/data/kv"
	"yubin/common/util"
)

var _ kv.Storage[int] = (*DB[int])(nil)

type DB[T any] struct {
	client redis.UniversalClient
	codec  codec.Codec[T, []byte]
	prefix string
	keySet string
}

func New[T any](client redis.UniversalClient, codec codec.Codec[T, []byte], prefix string, keySet string) *DB[T] {
	return &DB[T]{
		client: client,
		codec:  codec,
		prefix: prefix,
		keySet: keySet,
	}
}

func (db *DB[T]) Set(ctx context.Context, id string, item T) error {
	payload, err := db.codec.Encode(item)
	if err != nil {
		return err
	}

	return db.client.Watch(ctx, func(tx *redis.Tx) error {
		if cmd := tx.SAdd(ctx, db.keySet, db.key(id)); cmd.Err() != nil {
			return cmd.Err()
		}

		if cmd := tx.Set(ctx, db.key(id), payload, 0); cmd.Err() != nil {
			return cmd.Err()
		}

		return nil
	})
}

func (db *DB[T]) Range(ctx context.Context, order kv.Order, offset int, limit int) ([]T, error) {
	cfg := redis.Sort{
		Get:   []string{"*"},
		Alpha: true,
	}

	switch order {
	case kv.ASC:
		cfg.Order = "ASC"
	case kv.DESC:
		cfg.Order = "DESC"
	}

	if offset >= 0 {
		cfg.Offset = int64(offset)
	}

	if limit >= 0 {
		cfg.Count = int64(limit)
	}

	values, err := db.client.Sort(ctx, db.keySet, &cfg).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to sort: %s", err)
	}

	var result []T
	for _, value := range values {
		obj, err := db.codec.Decode([]byte(value))
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %s", err)
		}
		result = append(result, obj)
	}

	return result, nil
}

func (db *DB[T]) Get(ctx context.Context, id string) (T, error) {
	cmd := db.client.Get(ctx, id)
	if cmd.Err() != nil {
		err := cmd.Err()
		if err == redis.Nil {
			err = data.NotFoundErr
		}
		return util.Default[T](), err
	}

	raw, err := cmd.Bytes()
	if err != nil {
		return util.Default[T](), err
	}

	val, err := db.codec.Decode(raw)
	if err != nil {
		return util.Default[T](), err
	}

	return val, nil
}

func (db *DB[T]) Delete(ctx context.Context, id string) error {
	return db.client.Watch(ctx, func(tx *redis.Tx) error {
		if cmd := tx.SRem(ctx, db.keySet, db.key(id)); cmd.Err() != nil {
			return cmd.Err()
		}

		if cmd := tx.Del(ctx, db.key(id)); cmd.Err() != nil {
			return cmd.Err()
		}

		return nil
	})
}

func (db *DB[T]) key(id string) string {
	return db.prefix + id
}

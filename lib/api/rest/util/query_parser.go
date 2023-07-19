package util

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"yubin/common/data/kv"
	"yubin/common/data/record"
)

type RangeQuery struct {
	Offset *int    `json:"offset" form:"offset"`
	Limit  *int    `json:"limit" form:"limit"`
	Order  *string `json:"order" form:"order"`
}

func ParseRangeQuery(ctx *gin.Context) RangeQuery {
	var query RangeQuery
	if err := ctx.BindQuery(&query); err != nil {
		return RangeQuery{}
	}
	return query
}

func KvRangeQuery[T any](ctx context.Context, storage kv.Storage[T], query RangeQuery) ([]T, error) {
	offset, limit, order := prepareQuery(query)

	return storage.Range(ctx, order, offset, limit)
}

func RecordsQuery(ctx context.Context, set record.Set, sortKey string, query RangeQuery) ([]record.R, error) {
	offset, limit, order := prepareQuery(query)

	cursor := set.Cursor()
	// TODO move "ordering" to a mututal package
	cursor = cursor.Sort(record.Sorter{sortKey, record.Ordering(order)})

	if offset > 0 {
		cursor = cursor.Skip(offset)
	}

	if limit > 0 {
		cursor = cursor.Limit(limit)
	}

	iter := cursor.Iter()
	var result []record.R
	for {
		item, err := iter.Next(ctx)
		if err == record.ErrIterEmpty {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get a record from the iterator: %s", err)
		}
		result = append(result, item)
	}

	return result, nil
}

func prepareQuery(query RangeQuery) (offset int, limit int, order kv.Order) {
	offset = -1
	limit = -1
	order = kv.ASC

	if query.Offset != nil {
		offset = *query.Offset
	}

	if query.Limit != nil {
		limit = *query.Limit
	}

	if query.Order != nil {
		order = kv.Order(strings.ToUpper(strings.TrimSpace(*query.Order)))
	}

	return
}

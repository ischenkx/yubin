package crud

type Query struct {
	Offset, Limit *int
}

func NewQuery() *Query {
	return &Query{}
}

func (q *Query) WithOffset(offset int) *Query {
	q.Offset = &offset
	return q
}

func (q *Query) WithLimit(limit int) *Query {
	q.Limit = &limit
	return q
}

type MultiRead[T any] interface {
	Query(query *Query) ([]T, error)
}

type Read[ID, T any] interface {
	Get(ID) (T, bool, error)
}

type Create[T any] interface {
	Create(T) (T, error)
}

type Update[T any] interface {
	Update(T) error
}

type Delete[ID any] interface {
	Delete(ID) error
}

type CRUD[ID, T any] interface {
	Read[ID, T]
	MultiRead[T]
	Create[T]
	Update[T]
	Delete[ID]
}

type PairKey[T1, T2 any] struct {
	Item1 T1
	Item2 T2
}

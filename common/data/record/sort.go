package record

import "fmt"

type Ordering string

const (
	ASC        Ordering = "asc"
	DESC                = "desc"
	NoOrdering          = "no_ordering"
)

type Sorter struct {
	Key      string
	Ordering Ordering
}

func (s Sorter) String() string {
	return fmt.Sprintf("'%s': %s", s.Key, s.Ordering)
}

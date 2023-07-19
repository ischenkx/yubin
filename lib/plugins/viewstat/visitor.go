package viewstat

import (
	"context"
)

type Identifier struct {
	Publication string
	User        string
}

type Visitor interface {
	GenerateLink(identifier Identifier) (string, error)
	Visits(ctx context.Context) <-chan Identifier
}

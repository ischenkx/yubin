package viewstat

import "smtp-client/pkg/channel"

type Identifier struct {
	Publication string
	User        string
}

type Visitor interface {
	GenerateLink(identifier Identifier) (string, error)
	Visits() channel.Handle[Identifier]
}

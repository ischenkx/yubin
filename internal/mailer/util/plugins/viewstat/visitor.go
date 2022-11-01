package viewstat

type Identifier struct {
	Publication string
	User        string
}

type Handle interface {
	Chan() <-chan Identifier
	Close()
}

type Visitor interface {
	GenerateLink(identifier Identifier) (string, error)
	Visits() Handle
}

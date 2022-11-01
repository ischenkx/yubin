package mail

type Source struct {
	Address  string
	Password string
	Host     string
	Port     int
}

type Package[PayloadFormat any] struct {
	Source      Source
	Destination []string
	Payload     PayloadFormat
}

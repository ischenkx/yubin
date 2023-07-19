package mail

type Source struct {
	Address  string
	Password string
	Host     string
	Port     int
}

type Package[Payload any] struct {
	Source     Source
	Recipients []string
	Payload    Payload
}

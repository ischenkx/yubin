package mailer

type Report struct {
	PublicationID string
	Status        string
	Failed        []string
	OK            []string
}

type PersonalReport struct {
	PublicationID string
	UserID        string
	Status        string
	Meta          map[string]any
}

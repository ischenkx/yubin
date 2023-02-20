package yubin

import "time"

type SendOptions struct {
	Topics   []string
	Users    []string
	SourceID string
	Template string
}

type PublishOptions struct {
	SendOptions SendOptions
	At          *time.Time
	Meta        map[string]any
}

type Publication struct {
	ID   string
	Info PublishOptions
}

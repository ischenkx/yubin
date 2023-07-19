package pubsub

import (
	"time"
	"yubin/common/data/record"
	"yubin/lib/api/rest/util"
	yubin "yubin/src"
	"yubin/src/publication"
	"yubin/src/subscription"
)

type UpdateSubscriptionDto struct {
	Meta map[string]any `json:"meta"`
}

type PublicationDto struct {
	ID         string         `json:"id"`
	Topics     []string       `json:"topics,omitempty"`
	Users      []string       `json:"users,omitempty"`
	Template   string         `json:"template,omitempty"`
	Source     string         `json:"source,omitempty"`
	Properties map[string]any `json:"meta,omitempty"`
}

func publication2dto(p publication.Publication) PublicationDto {
	return PublicationDto{
		ID:         p.ID,
		Topics:     p.Destination.Topics,
		Users:      p.Destination.Users,
		Template:   p.Template,
		Source:     p.Source,
		Properties: p.Properties,
	}
}

type ReportDto struct {
	Publication string         `json:"publication,omitempty"`
	User        string         `json:"user,omitempty"`
	Status      string         `json:"status,omitempty"`
	Meta        map[string]any `json:"meta,omitempty"`
}

func recordReport2dto(r record.R) ReportDto {
	return ReportDto{
		Publication: r["publication"].(string),
		User:        r["user"].(string),
		Status:      r["status"].(string),
		Meta:        r["meta"].(map[string]any),
	}
}

func report2dto(r yubin.Report) ReportDto {
	return ReportDto{
		Publication: r.Publication,
		User:        r.User,
		Status:      r.Status,
		Meta:        r.Meta,
	}
}

type SubscriptionDto struct {
	Subscriber string         `json:"subscriber,omitempty"`
	Topic      string         `json:"topic,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	Meta       map[string]any `json:"meta,omitempty"`
}

func subscription2dto(s subscription.Subscription) SubscriptionDto {
	return SubscriptionDto{
		Subscriber: s.Subscriber,
		Topic:      s.Topic,
		CreatedAt:  s.CreatedAt,
		Meta:       util.ValidateEmptyMap(s.Meta),
	}
}

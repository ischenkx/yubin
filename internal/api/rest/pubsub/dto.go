package pubsub

import (
	"smtp-client/internal/api/rest/util"
	"smtp-client/internal/mailer"
	"smtp-client/internal/mailer/subscription"
	"time"
)

type UpdateSubscriptionDto struct {
	Meta map[string]any `json:"meta"`
}

type PublicationDto struct {
	ID       string         `json:"id"`
	Topics   []string       `json:"topics,omitempty"`
	Users    []string       `json:"users,omitempty"`
	Template string         `json:"template,omitempty"`
	Source   string         `json:"source,omitempty"`
	At       *time.Time     `json:"at,omitempty"`
	Meta     map[string]any `json:"meta,omitempty"`
}

type ReportDto struct {
	PublicationID string   `json:"publication_id,omitempty"`
	Status        string   `json:"status,omitempty"`
	Failed        []string `json:"failed,omitempty"`
	OK            []string `json:"ok,omitempty"`
}

func report2dto(r mailer.Report) ReportDto {
	return ReportDto{
		PublicationID: r.PublicationID,
		Status:        r.Status,
		Failed:        util.ValidateEmptySlice(r.Failed),
		OK:            util.ValidateEmptySlice(r.OK),
	}
}

type PersonalReportDto struct {
	PublicationID string         `json:"publication_id,omitempty"`
	UserID        string         `json:"user_id,omitempty"`
	Status        string         `json:"status,omitempty"`
	Meta          map[string]any `json:"meta,omitempty"`
}

func personalReport2dto(r mailer.PersonalReport) PersonalReportDto {
	return PersonalReportDto{
		PublicationID: r.PublicationID,
		UserID:        r.UserID,
		Status:        r.Status,
		Meta:          util.ValidateEmptyMap(r.Meta),
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

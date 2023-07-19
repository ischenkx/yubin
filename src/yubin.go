package yubin

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"yubin/common/data/kv"
	"yubin/common/data/record"
	"yubin/src/mail"
	"yubin/src/publication"
	"yubin/src/subscription"
	"yubin/src/template"
	"yubin/src/user"
)

type Yubin struct {
	// System
	transport Transport[mail.Package[template.ParametrizedTemplate]]
	plugins   []Plugin

	// Repositories
	sources       kv.Storage[NamedSource]
	publications  kv.Storage[publication.Publication]
	users         kv.Storage[user.User]
	reports       record.Storage
	templates     template.Repo
	subscriptions subscription.Repo
}

func (m *Yubin) Sources() kv.Storage[NamedSource] {
	return m.sources
}

func (m *Yubin) Templates() template.Repo {
	return m.templates
}

func (m *Yubin) Publications() kv.Storage[publication.Publication] {
	return m.publications
}

func (m *Yubin) Reports() record.Storage {
	return m.reports
}

func (m *Yubin) Users() kv.Storage[user.User] {
	return m.users
}

func (m *Yubin) Subscriptions() subscription.Repo {
	return m.subscriptions
}

func (m *Yubin) New(ctx context.Context, options ...PublishOption) (publication.Publication, error) {
	var spec publication.Spec
	for _, option := range options {
		option(&spec)
	}

	pub := publication.Publication{
		ID:   uuid.New().String(),
		Spec: spec,
	}

	if err := m.Publications().Set(ctx, pub.ID, pub); err != nil {
		return publication.Publication{}, fmt.Errorf("failed to create the publication: %s", err)
	}

	return pub, nil
}

func (m *Yubin) Deliver(ctx context.Context, id string) ([]Report, error) {
	pub, err := m.publications.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get the given publication (id='%s'): %s", id, err)
	}

	tem, err := m.Templates().Get(ctx, pub.Template)

	if err != nil {
		return nil, fmt.Errorf("failed to get the given template (name='%s'): %s", pub.Template)
	}

	source, err := m.Sources().Get(ctx, pub.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to get the given source (name='%s'): %s", pub.Source, err)
	}

	recipientIDs := map[string]struct{}{}
	for _, id := range pub.Destination.Users {
		recipientIDs[id] = struct{}{}
	}
	for _, topic := range pub.Destination.Topics {
		subscribers, err := m.Subscriptions().Subscribers(ctx, topic)
		if err != nil {
			return nil, err
		}
		for _, subscriber := range subscribers {
			recipientIDs[subscriber] = struct{}{}
		}
	}

	recipients := make(map[string]user.User, len(recipientIDs))
	for id := range recipientIDs {
		u, err := m.Users().Get(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get a user (id=''): %s", id, err)
		}
		recipients[u.Email] = u
	}

	var reports []Report

	for email, u := range recipients {
		report := Report{
			Publication: pub.ID,
			User:        u.ID,
			Saved:       true,
			Status:      OK,
			Meta:        map[string]any{},
		}

		parametrization := map[string]any{
			"user": u,
			"meta": pub.Properties,
		}

		pack := mail.Package[template.ParametrizedTemplate]{
			Source:     source.Source,
			Recipients: []string{email},
			Payload: template.ParametrizedTemplate{
				Template:  tem,
				Parameter: parametrization,
			},
		}

		var interceptionError error

		for _, plugin := range m.plugins {
			if err := plugin.Intercept(m, pub, u, &pack); err != nil {
				interceptionError = err
				break
			}
		}

		if interceptionError != nil {
			report.Status = Failed
			report.Meta["interception_error"] = interceptionError.Error()
		} else if deliveryError := m.transport.Send(pack); deliveryError != nil {
			report.Status = Failed
			report.Meta["delivery_failure"] = err.Error()
		}

		if err := m.Reports().Insert(ctx, m.report2record(report)); err != nil {
			report.Saved = false
			report.Meta["saving_error"] = err.Error()
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func (m *Yubin) report2record(report Report) record.Record {
	return record.R{
		"publication": report.Publication,
		"user":        report.User,
		"status":      report.Status,
		"meta":        report.Meta,
	}
}

func (m *Yubin) use(p Plugin) error {
	if p == nil {
		return nil
	}

	if err := p.Init(m); err != nil {
		return fmt.Errorf("failed to init plugin: %s", err)
	}
	m.plugins = append(m.plugins, p)

	return nil
}

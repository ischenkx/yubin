package mailer

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
	"runtime"
	"smtp-client/internal/mailer/mail"
	"smtp-client/internal/mailer/subscription"
	"smtp-client/internal/mailer/template"
	"smtp-client/internal/mailer/user"
	"smtp-client/pkg/data/crud"
	"time"
)

type Mailer struct {
	scheduler Scheduler[string]
	queue     Queue[string]
	post      Delivery[mail.Package[template.ParametrizedTemplate]]
	repo      Repo
	plugins   []Plugin
}

func New(
	scheduler Scheduler[string],
	queue Queue[string],
	delivery Delivery[mail.Package[template.ParametrizedTemplate]],
	repo Repo) *Mailer {
	return &Mailer{
		scheduler: scheduler,
		queue:     queue,
		post:      delivery,
		repo:      repo,
	}
}

func (m *Mailer) UsePlugin(p Plugin) {
	if p == nil {
		return
	}
	p.Init(m)
	m.plugins = append(m.plugins, p)
}

func (m *Mailer) Sources() crud.CRUD[string, NamedSource] {
	return m.repo.Sources()
}

func (m *Mailer) Templates() template.Repo {
	return m.repo.Templates()
}

func (m *Mailer) Publications() crud.CRUD[string, Publication] {
	return m.repo.Publications()
}

func (m *Mailer) Reports() crud.CRUD[string, Report] {
	return m.repo.Reports()
}

func (m *Mailer) PersonalReports() crud.CRUD[crud.PairKey[string, string], PersonalReport] {
	return m.repo.PersonalReports()
}

func (m *Mailer) Users() user.Repo {
	return m.repo.Users()
}

func (m *Mailer) Subscriptions() subscription.Repo {
	return m.repo.Subscriptions()
}

func (m *Mailer) Publish(options ...PublishOption) (string, error) {
	var opts PublishOptions
	for _, opt := range options {
		opt(&opts)
	}

	publication := Publication{
		ID:   uuid.New().String(),
		Info: opts,
	}

	publication, err := m.repo.Publications().Create(publication)
	if err != nil {
		return "", err
	}

	if opts.At != nil {
		return publication.ID, m.pushToScheduler(*opts.At, publication.ID)
	}
	return publication.ID, m.pushToQueue(publication.ID)
}

func (m *Mailer) pushToScheduler(at time.Time, id string) error {
	if m.scheduler == nil {
		return errors.New("no scheduler provided")
	}
	return m.scheduler.Schedule(id, at)
}

func (m *Mailer) pushToQueue(id string) error {
	return m.queue.Push(id)
}

func (m *Mailer) deliver(id string) (Report, error) {
	publication, ok, err := m.repo.Publications().Get(id)
	if err != nil {
		return Report{}, err
	}
	if !ok {
		return Report{}, errors.New("publication not found")
	}

	opts := publication.Info.SendOptions

	temp, ok, err := m.Templates().Get(opts.Template)
	if err != nil {
		return Report{}, err
	}
	if !ok {
		return Report{}, errors.New("template not found")
	}

	source, ok, err := m.repo.Sources().Get(opts.SourceID)
	if err != nil {
		return Report{}, err
	}
	if !ok {
		return Report{}, errors.New("source not found")
	}

	recipientIDS := map[string]struct{}{}
	for _, id := range opts.Users {
		recipientIDS[id] = struct{}{}
	}
	for _, topic := range opts.Topics {
		subscribers, err := m.Subscriptions().Subscribers(topic, nil)
		if err != nil {
			return Report{}, err
		}
		for _, subscriber := range subscribers {
			recipientIDS[subscriber] = struct{}{}
		}
	}

	recipients := make(map[string]user.User, len(recipientIDS))
	for id := range recipientIDS {
		u, ok, err := m.Users().Get(id)
		if err != nil {
			return Report{}, err
		}
		if !ok {
			continue
		}
		recipients[u.Email] = u
	}

	var report Report
	report.PublicationID = publication.ID
	for email, u := range recipients {
		parametrization := map[string]any{
			"user": u,
			"meta": publication.Info.Meta,
		}

		pack := mail.Package[template.ParametrizedTemplate]{
			Source:      source.Source,
			Destination: []string{email},
			Payload: template.ParametrizedTemplate{
				Template:  temp,
				Parameter: parametrization,
			},
		}

		for _, plugin := range m.plugins {
			plugin.Intercept(m, publication, u, &pack)
		}

		if err := m.post.Deliver(pack); err != nil {
			report.Failed = append(report.Failed, u.ID)
		} else {
			report.OK = append(report.OK, u.ID)
		}
	}
	return report, nil
}

func (m *Mailer) processQueue(ctx context.Context) {
	handle := m.queue.Handle()
	defer handle.Close()

processor:
	for {
		select {
		case <-ctx.Done():
			break processor
		case id := <-handle.Chan():
			report, err := m.deliver(id)
			if err != nil {
				report = Report{
					PublicationID: id,
					Status:        "failed",
				}
			} else {
				report.Status = "ok"
			}

			if _, err := m.Reports().Create(report); err != nil {
				log.Println("failed to create a report:", err)
			}
			for _, userID := range report.Failed {
				personalReport := PersonalReport{
					PublicationID: report.PublicationID,
					UserID:        userID,
					Status:        "failed",
					Meta:          nil,
				}
				if _, err := m.PersonalReports().Create(personalReport); err != nil {
					log.Println("failed to create a personal report:", err)
				}
			}

			for _, userID := range report.OK {
				personalReport := PersonalReport{
					PublicationID: report.PublicationID,
					UserID:        userID,
					Status:        "ok",
					Meta:          nil,
				}
				if _, err := m.PersonalReports().Create(personalReport); err != nil {
					log.Println("failed to create a personal report:", err)
				}
			}
		}
	}
}

func (m *Mailer) Run(ctx context.Context) {
	for i := 0; i < runtime.NumCPU(); i++ {
		go m.processQueue(ctx)
	}
}

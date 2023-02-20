package yubin

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
	"runtime"
	"smtp-client/internal/yubin/mail"
	"smtp-client/internal/yubin/subscription"
	"smtp-client/internal/yubin/template"
	"smtp-client/internal/yubin/user"
	"smtp-client/pkg/data/crud"
	"time"
)

type Yubin struct {
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
	repo Repo) *Yubin {
	return &Yubin{
		scheduler: scheduler,
		queue:     queue,
		post:      delivery,
		repo:      repo,
	}
}

func (m *Yubin) UsePlugin(p Plugin) {
	if p == nil {
		return
	}
	p.Init(m)
	m.plugins = append(m.plugins, p)
}

func (m *Yubin) Sources() crud.CRUD[string, NamedSource] {
	return m.repo.Sources()
}

func (m *Yubin) Templates() template.Repo {
	return m.repo.Templates()
}

func (m *Yubin) Publications() crud.CRUD[string, Publication] {
	return m.repo.Publications()
}

func (m *Yubin) Reports() crud.CRUD[string, Report] {
	return m.repo.Reports()
}

func (m *Yubin) PersonalReports() crud.CRUD[crud.PairKey[string, string], PersonalReport] {
	return m.repo.PersonalReports()
}

func (m *Yubin) Users() user.Repo {
	return m.repo.Users()
}

func (m *Yubin) Subscriptions() subscription.Repo {
	return m.repo.Subscriptions()
}

func (m *Yubin) Publish(options ...PublishOption) (string, error) {
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

func (m *Yubin) pushToScheduler(at time.Time, id string) error {
	if m.scheduler == nil {
		return errors.New("no scheduler provided")
	}
	return m.scheduler.Schedule(id, at)
}

func (m *Yubin) pushToQueue(id string) error {
	return m.queue.Push(id)
}

func (m *Yubin) deliver(id string) (Report, error) {
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
			log.Println("failed to deliver:", err)
			report.Failed = append(report.Failed, u.ID)
		} else {
			report.OK = append(report.OK, u.ID)
		}
	}
	return report, nil
}

func (m *Yubin) processQueue(ctx context.Context) {
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

func (m *Yubin) processScheduler(ctx context.Context) {
	sub := m.scheduler.Handle()
	defer sub.Close()
processor:
	for {
		select {
		case id := <-sub.Chan():
			if err := m.pushToQueue(id); err != nil {
				log.Println("failed to push to queue a scheduled publication:", err)
			}
		case <-ctx.Done():
			break processor
		}
	}
}

func (m *Yubin) Run(ctx context.Context) {
	for i := 0; i < runtime.NumCPU(); i++ {
		go m.processQueue(ctx)
	}

	go m.processScheduler(ctx)
}

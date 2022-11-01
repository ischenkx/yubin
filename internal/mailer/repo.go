package mailer

import (
	"smtp-client/internal/mailer/mail"
	"smtp-client/internal/mailer/subscription"
	"smtp-client/internal/mailer/template"
	"smtp-client/internal/mailer/user"
	"smtp-client/pkg/data/crud"
)

type NamedSource struct {
	Name string
	mail.Source
}

type Repo interface {
	Users() user.Repo
	Subscriptions() subscription.Repo
	Templates() template.Repo
	Sources() crud.CRUD[string, NamedSource]
	Publications() crud.CRUD[string, Publication]
	Reports() crud.CRUD[string, Report]
	PersonalReports() crud.CRUD[crud.PairKey[string, string], PersonalReport]
}

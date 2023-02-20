package yubin

import (
	"smtp-client/internal/yubin/mail"
	"smtp-client/internal/yubin/subscription"
	"smtp-client/internal/yubin/template"
	"smtp-client/internal/yubin/user"
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

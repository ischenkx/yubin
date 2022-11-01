package postgrepo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"smtp-client/internal/mailer"
	"smtp-client/internal/mailer/subscription"
	"smtp-client/internal/mailer/template"
	"smtp-client/internal/mailer/user"
	"smtp-client/pkg/data/crud"
)

type Repo struct {
	users           *Users
	subscriptions   *Subscriptions
	templates       *Templates
	sources         *Sources
	publications    *Publications
	reports         *Reports
	personalReports *PersonalReports
}

func New(conn *pgx.Conn, templateEngine TemplateEngine) *Repo {
	return &Repo{
		sources:         &Sources{conn: conn},
		users:           &Users{conn: conn},
		subscriptions:   &Subscriptions{conn: conn},
		templates:       &Templates{conn: conn, engine: templateEngine},
		publications:    &Publications{conn: conn},
		reports:         &Reports{conn: conn},
		personalReports: &PersonalReports{conn: conn},
	}
}

func (r *Repo) Sources() crud.CRUD[string, mailer.NamedSource] {
	return r.sources
}

func (r *Repo) Users() user.Repo {
	return r.users
}

func (r *Repo) Subscriptions() subscription.Repo {
	return r.subscriptions
}

func (r *Repo) Templates() template.Repo {
	return r.templates
}

func (r *Repo) Publications() crud.CRUD[string, mailer.Publication] {
	return r.publications
}

func (r *Repo) Reports() crud.CRUD[string, mailer.Report] {
	return r.reports
}

func (r *Repo) PersonalReports() crud.CRUD[crud.PairKey[string, string], mailer.PersonalReport] {
	return r.personalReports
}

func (r *Repo) Init(ctx context.Context) error {
	if err := r.users.InitTable(ctx); err != nil {
		return err
	}
	if err := r.templates.InitTable(ctx); err != nil {
		return err
	}
	if err := r.subscriptions.InitTable(ctx); err != nil {
		return err
	}
	return nil
}

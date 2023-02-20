package postgrepo

import (
	"context"
	"log"
	"smtp-client/internal/yubin"
	"smtp-client/internal/yubin/subscription"
	"smtp-client/internal/yubin/template"
	"smtp-client/internal/yubin/user"
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

func (r *Repo) Sources() crud.CRUD[string, yubin.NamedSource] {
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

func (r *Repo) Publications() crud.CRUD[string, yubin.Publication] {
	return r.publications
}

func (r *Repo) Reports() crud.CRUD[string, yubin.Report] {
	return r.reports
}

func (r *Repo) PersonalReports() crud.CRUD[crud.PairKey[string, string], yubin.PersonalReport] {
	return r.personalReports
}

func (r *Repo) Init(ctx context.Context) {
	r.initialize(ctx, "users", r.users)
	r.initialize(ctx, "templates", r.templates)
	r.initialize(ctx, "subscriptions", r.subscriptions)
	r.initialize(ctx, "reports", r.reports)
	r.initialize(ctx, "sources", r.sources)
	r.initialize(ctx, "publications", r.publications)
	r.initialize(ctx, "personal reports", r.personalReports)

}

func (r *Repo) initialize(ctx context.Context, label string, i initializer) {
	if err := i.Init(ctx); err != nil {
		log.Printf("failed to initialize '%s': %s", label, err)
	}
}

type initializer interface {
	Init(ctx context.Context) error
}

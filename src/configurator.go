package yubin

import (
	"yubin/common/data/kv"
	"yubin/common/data/record"
	"yubin/src/mail"
	"yubin/src/publication"
	"yubin/src/subscription"
	"yubin/src/template"
	"yubin/src/user"
)

type Configurator struct {
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

func Configure() Configurator {
	return Configurator{}
}

func (cfg Configurator) Transport(t Transport[mail.Package[template.ParametrizedTemplate]]) Configurator {
	cfg.transport = t
	return cfg
}

func (cfg Configurator) Plugins(plugins ...Plugin) Configurator {
	cfg.plugins = append(cfg.plugins, plugins...)
	return cfg
}

func (cfg Configurator) Sources(obj kv.Storage[NamedSource]) Configurator {
	cfg.sources = obj
	return cfg
}
func (cfg Configurator) Publications(obj kv.Storage[publication.Publication]) Configurator {
	cfg.publications = obj
	return cfg
}
func (cfg Configurator) Users(obj kv.Storage[user.User]) Configurator {
	cfg.users = obj
	return cfg
}
func (cfg Configurator) Reports(obj record.Storage) Configurator {
	cfg.reports = obj
	return cfg
}
func (cfg Configurator) Templates(obj template.Repo) Configurator {
	cfg.templates = obj
	return cfg
}
func (cfg Configurator) Subscriptions(obj subscription.Repo) Configurator {
	cfg.subscriptions = obj
	return cfg
}

func (cfg Configurator) Compile() *Yubin {
	return &Yubin{
		transport:     cfg.transport,
		plugins:       cfg.plugins,
		sources:       cfg.sources,
		publications:  cfg.publications,
		users:         cfg.users,
		reports:       cfg.reports,
		templates:     cfg.templates,
		subscriptions: cfg.subscriptions,
	}
}

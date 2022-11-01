package mailer

import (
	"smtp-client/internal/mailer/mail"
	"smtp-client/internal/mailer/template"
	"smtp-client/internal/mailer/user"
)

type Plugin interface {
	Init(*Mailer)
	Intercept(*Mailer,
		Publication,
		user.User,
		*mail.Package[template.ParametrizedTemplate])
}

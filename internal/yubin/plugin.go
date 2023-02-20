package yubin

import (
	"smtp-client/internal/yubin/mail"
	"smtp-client/internal/yubin/template"
	"smtp-client/internal/yubin/user"
)

type Plugin interface {
	Init(*Yubin)
	Intercept(*Yubin,
		Publication,
		user.User,
		*mail.Package[template.ParametrizedTemplate])
}

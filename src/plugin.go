package yubin

import (
	"yubin/src/mail"
	"yubin/src/publication"
	"yubin/src/template"
	"yubin/src/user"
)

type Plugin interface {
	Init(*Yubin) error
	Intercept(
		yubin *Yubin,
		publication publication.Publication,
		user user.User,
		_package *mail.Package[template.ParametrizedTemplate]) error
}

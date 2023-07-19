package yubin

import "yubin/src/mail"

type NamedSource struct {
	Name string
	mail.Source
}

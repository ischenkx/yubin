package templates

import "yubin/src/template"

type Engine interface {
	Convert(model Model) (template.Template, error)
}

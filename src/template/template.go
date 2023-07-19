package template

import (
	"io"
	"yubin/common/data/kv"
)

type Template interface {
	SubTemplate(name string) (Template, bool)
	SubTemplates() map[string]Template
	WriteTo(writer io.Writer, data any) error
	Meta() map[string]any
	Raw() string
	Name() string
}

type Repo kv.Storage[Template]

package template

import (
	"io"
	"smtp-client/pkg/data/crud"
)

type Template interface {
	SubTemplate(name string) (Template, bool)
	SubTemplates() map[string]Template
	WriteTo(writer io.Writer, data any) error
	Meta() map[string]any
	Raw() string
	Name() string
}

type ParametrizedTemplate struct {
	Template  Template
	Parameter any
}

func (t ParametrizedTemplate) Meta() map[string]any {
	return t.Template.Meta()
}

func (t ParametrizedTemplate) SubTemplates() map[string]ParametrizedTemplate {
	m := map[string]ParametrizedTemplate{}
	for name, temp := range t.Template.SubTemplates() {
		m[name] = ParametrizedTemplate{
			Template:  temp,
			Parameter: t.Parameter,
		}
	}
	return m
}

func (t ParametrizedTemplate) SubTemplate(name string) (ParametrizedTemplate, bool) {
	st, ok := t.Template.SubTemplate(name)
	if !ok {
		return ParametrizedTemplate{}, false
	}
	return ParametrizedTemplate{
		Template:  st,
		Parameter: t.Parameter,
	}, true
}

func (t ParametrizedTemplate) WriteTo(writer io.Writer) error {
	return t.Template.WriteTo(writer, t.Parameter)
}

func (t ParametrizedTemplate) Raw() string {
	return t.Template.Raw()
}

func (t ParametrizedTemplate) Name() string {
	return t.Template.Name()
}

type Repo crud.CRUD[string, Template]

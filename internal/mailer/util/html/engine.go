package html

import (
	htmlTemplate "html/template"
	"io"
	"smtp-client/internal/mailer/template"
	"smtp-client/internal/mailer/util/postgrepo"
)

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) convert(model postgrepo.TemplateModel) (*Template, error) {
	t := htmlTemplate.New(model.Name)
	if _, err := t.Parse(model.Data); err != nil {
		return nil, err
	}

	res := &Template{
		Data:         model.Data,
		Sub:          map[string]*Template{},
		HTMLTemplate: t,
		MetaData:     model.Meta,
	}

	for name, sub := range model.SubTemplates {
		t, err := e.convert(sub)
		if err != nil {
			return nil, err
		}
		res.Sub[name] = t
	}

	return res, nil
}

func (e *Engine) Convert(model postgrepo.TemplateModel) (template.Template, error) {
	return e.convert(model)
}

type Template struct {
	Data         string
	MetaData     map[string]any
	Sub          map[string]*Template
	HTMLTemplate *htmlTemplate.Template
}

func (t *Template) SubTemplate(name string) (template.Template, bool) {
	if t.Sub == nil {
		return nil, false
	}
	res, ok := t.Sub[name]
	return res, ok
}

func (t *Template) SubTemplates() map[string]template.Template {
	m := map[string]template.Template{}
	for name, sub := range t.Sub {
		m[name] = sub
	}
	return m
}

func (t *Template) Name() string {
	return t.HTMLTemplate.Name()
}

func (t *Template) Meta() map[string]any {
	return t.MetaData
}

func (t *Template) Raw() string {
	return t.Data
}

func (t *Template) WriteTo(w io.Writer, data any) error {
	return t.HTMLTemplate.Execute(w, data)
}

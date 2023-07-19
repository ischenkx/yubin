package templates

import (
	"errors"
	"io"
	"yubin/lib/api/rest/util"
	"yubin/src/template"
)

type UpdateTemplateDto struct {
	TemplateName string                  `json:"name,omitempty"`
	Data         *string                 `json:"data,omitempty"`
	Meta         *map[string]any         `json:"meta"`
	Sub          *map[string]TemplateDto `json:"sub"`
}

type TemplateDto struct {
	TemplateName string                 `json:"name,omitempty"`
	Data         string                 `json:"data,omitempty"`
	MetaData     map[string]any         `json:"meta"`
	Sub          map[string]TemplateDto `json:"sub"`
}

func (t TemplateDto) SubTemplate(name string) (template.Template, bool) {
	sub, ok := t.Sub[name]
	return sub, ok
}

func (t TemplateDto) SubTemplates() map[string]template.Template {
	m := map[string]template.Template{}
	for name, sub := range t.Sub {
		m[name] = sub
	}
	return m
}

func (t TemplateDto) Meta() map[string]any {
	return t.MetaData
}

func (t TemplateDto) WriteTo(writer io.Writer, data any) error {
	return errors.New("can't write")
}

func (t TemplateDto) Raw() string {
	return t.Data
}

func (t TemplateDto) Name() string {
	return t.TemplateName
}

func template2dto(t template.Template) TemplateDto {
	dto := TemplateDto{
		TemplateName: t.Name(),
		Data:         t.Raw(),
		MetaData:     util.ValidateEmptyMap(t.Meta()),
		Sub:          map[string]TemplateDto{},
	}

	for name, sub := range t.SubTemplates() {
		val := template2dto(sub)
		if val.TemplateName != "" {
			val.TemplateName = name
		}
		dto.Sub[name] = val
	}
	return dto
}

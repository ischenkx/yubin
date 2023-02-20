package viewstat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"smtp-client/internal/yubin"
	"smtp-client/internal/yubin/mail"
	"smtp-client/internal/yubin/template"
	"smtp-client/internal/yubin/user"
	"smtp-client/pkg/data/crud"
)

type ViewStat struct {
	mailer  *yubin.Yubin
	visitor Visitor
}

func New(visitor Visitor) *ViewStat {
	return &ViewStat{
		visitor: visitor,
	}
}

func (v *ViewStat) Run(ctx context.Context) {
	handle := v.visitor.Visits()
	defer handle.Close()
	for {
		select {
		case <-ctx.Done():
			return
		case visit := <-handle.Chan():
			if err := v.postView(visit); err != nil {
				log.Println("failed to post a view:", err)
			}
		}
	}
}

func (v *ViewStat) Init(mailer *yubin.Yubin) {
	v.mailer = mailer
}

func (v *ViewStat) Intercept(mailer *yubin.Yubin,
	publication yubin.Publication,
	user user.User,
	pack *mail.Package[template.ParametrizedTemplate]) {
	if publication.Info.Meta != nil {
		if _, ok := publication.Info.Meta["disable_viewstat"]; ok {
			return
		}
	}

	link, err := v.visitor.GenerateLink(Identifier{
		Publication: publication.ID,
		User:        user.ID,
	})

	if err != nil {
		log.Println("failed to generate a viewstat link:", err)
		return
	}

	suffix := fmt.Sprintf("<img src=\"%s\" alt=\"\">", link)
	pack.Payload.Template = payload{
		suffix:  suffix,
		initial: pack.Payload.Template,
	}
}

func (v *ViewStat) postView(visit Identifier) error {
	if v.mailer == nil {
		return errors.New("yubin not provided")
	}
	report, ok, err := v.mailer.PersonalReports().Get(crud.PairKey[string, string]{
		Item1: visit.Publication,
		Item2: visit.User,
	})
	if err != nil {
		return err
	}
	if !ok {
		log.Println("personal report not found")
		return nil
	}

	if report.Meta == nil {
		report.Meta = map[string]any{}
	}
	report.Meta["viewstat"] = map[string]any{
		"viewed": true,
	}

	if err := v.mailer.PersonalReports().Update(report); err != nil {
		return err
	}
	return nil
}

type payload struct {
	suffix  string
	initial template.Template
}

func (p payload) SubTemplate(name string) (template.Template, bool) {
	return p.initial.SubTemplate(name)
}

func (p payload) SubTemplates() map[string]template.Template {
	return p.initial.SubTemplates()
}

func (p payload) WriteTo(writer io.Writer, data any) error {
	if err := p.initial.WriteTo(writer, data); err != nil {
		return err
	}
	_, err := writer.Write([]byte(p.suffix))
	return err
}

func (p payload) Meta() map[string]any {
	return p.initial.Meta()
}

func (p payload) Raw() string {
	return p.initial.Raw()
}

func (p payload) Name() string {
	return p.initial.Name()
}

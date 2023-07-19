package viewstat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"yubin/common/data/record"
	"yubin/src"
	"yubin/src/mail"
	"yubin/src/publication"
	"yubin/src/template"
	"yubin/src/user"
)

type ViewStat struct {
	yubin   *yubin.Yubin
	visitor Visitor
}

func New(visitor Visitor) *ViewStat {
	return &ViewStat{
		visitor: visitor,
	}
}

func (v *ViewStat) Run(ctx context.Context) {
	channel := v.visitor.Visits(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case visit := <-channel:
			if err := v.postView(ctx, visit); err != nil {
				log.Println("failed to post a view:", err)
			}
		}
	}
}

func (v *ViewStat) Init(yubin *yubin.Yubin) error {
	v.yubin = yubin
	return nil
}

func (v *ViewStat) Intercept(mailer *yubin.Yubin,
	publication publication.Publication,
	user user.User,
	pack *mail.Package[template.ParametrizedTemplate]) error {
	if publication.Properties != nil {
		if _, ok := publication.Properties["disable_viewstat"]; ok {
			return nil
		}
	}

	link, err := v.visitor.GenerateLink(Identifier{
		Publication: publication.ID,
		User:        user.ID,
	})

	if err != nil {
		log.Println("failed to generate a viewstat link:", err)
		return nil
	}

	suffix := fmt.Sprintf("<img src=\"%s\" alt=\"\">", link)
	pack.Payload.Template = payload{
		suffix:  suffix,
		initial: pack.Payload.Template,
	}

	return nil
}

func (v *ViewStat) postView(ctx context.Context, visit Identifier) error {
	if v.yubin == nil {
		return errors.New("yubin not provided")
	}

	err := v.yubin.
		Reports().
		Filter(
			record.E{"publication", visit.Publication},
			record.E{"user", visit.User},
		).
		Update(ctx, record.R{"viewed": true}, record.R{"user": visit.User, "publication": visit.Publication})

	if err != nil {
		return fmt.Errorf("failed to update the report: %s", err)
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

package smtp

import (
	"bytes"
	"fmt"
	"io"
	"net/smtp"
	"smtp-client/internal/mailer/mail"
	"smtp-client/internal/mailer/template"
	"strings"
)

type writerTo interface {
	WriteTo(writer io.Writer) error
}

type Delivery struct{}

func (d Delivery) Deliver(message mail.Package[template.ParametrizedTemplate]) error {
	mes := NewMessage()

	if subject, ok := message.Payload.SubTemplate("subject"); ok {
		s, err := writerTo2String(subject)
		if err != nil {
			return err
		}
		mes.Headers()["Subject"] = s
	}

	if subject, ok := message.Payload.SubTemplate("from"); ok {
		s, err := writerTo2String(subject)
		if err != nil {
			return err
		}
		mes.Headers()["From"] = s
	}

	mes.Headers()["To"] = strings.Join(message.Destination, ",")
	mes.Headers()["Content-Type"] = "text/html"

	if meta := message.Payload.Meta(); meta != nil {
		if rawHeaders, ok := meta["headers"]; ok {
			switch headers := rawHeaders.(type) {
			case map[string]string:
				mes.Headers().Merge(headers)
			case map[string]any:
				m := map[string]string{}
				for key, val := range headers {
					if stringVal, ok := val.(string); ok {
						m[key] = stringVal
					}
				}
				mes.Headers().Merge(m)
			}
		}
	}

	if headers, ok := message.Payload.SubTemplate("headers"); ok {
		for name, val := range headers.SubTemplates() {
			s, err := writerTo2String(val)
			if err != nil {
				return err
			}
			mes.Headers()[name] = s
		}
	}

	if err := message.Payload.WriteTo(mes.Payload()); err != nil {
		return err
	}

	auth := smtp.PlainAuth("",
		message.Source.Address,
		message.Source.Password,
		message.Source.Host)

	payload, err := mes.Build()
	if err != nil {
		return err
	}

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", message.Source.Host, message.Source.Port),
		auth,
		message.Source.Address,
		message.Destination,
		payload,
	)
}

type Headers map[string]string

func (h Headers) WriteTo(w io.Writer) error {
	if len(h) == 0 {
		return nil
	}

	for name, val := range h {
		line := fmt.Sprintf("%s: %s\r\n", name, val)
		if _, err := w.Write([]byte(line)); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}
	return nil
}

func (h Headers) Merge(other Headers) {
	if other == nil {
		return
	}
	for name, val := range other {
		h[name] = val
	}
}

type Message struct {
	headers Headers
	payload *bytes.Buffer
}

func NewMessage() *Message {
	return &Message{
		headers: Headers{},
		payload: bytes.NewBuffer(nil),
	}
}

func (m *Message) Headers() Headers {
	return m.headers
}

func (m *Message) Payload() io.Writer {
	return m.payload
}

func (m *Message) Build() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := m.headers.WriteTo(buf); err != nil {
		return nil, err
	}
	if _, err := m.payload.WriteTo(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writerTo2String(w writerTo) (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := w.WriteTo(buf); err != nil {
		return "", nil
	}
	return buf.String(), nil
}

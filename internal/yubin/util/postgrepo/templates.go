package postgrepo

import (
	"context"
	"fmt"
	"smtp-client/internal/yubin/template"
	"smtp-client/pkg/data/crud"
)

type TemplateEngine interface {
	Convert(model TemplateModel) (template.Template, error)
}

type TemplateModel struct {
	Name         string
	Data         string
	Meta         map[string]any
	SubTemplates map[string]TemplateModel
}

type Templates struct {
	engine TemplateEngine
	conn   *pgx.Conn
}

func (templates *Templates) Init(ctx context.Context) error {
	sql := `
		CREATE TABLE templates (
		    data varchar(65000),
		    name varchar(128) PRIMARY KEY,
		    meta json,
		    subtemplates json
		)
	`

	_, err := templates.conn.Exec(ctx, sql)
	return err
}

func (templates *Templates) Get(name string) (template.Template, bool, error) {
	sql := `
		SELECT data, name, meta, subtemplates FROM templates
		WHERE name=$1
	`

	result, err := templates.conn.Query(context.TODO(), sql, name)
	if err != nil {
		return nil, false, err
	}
	defer result.Close()
	if !result.Next() {
		return nil, false, nil
	}

	var model TemplateModel
	if err := result.Scan(&model.Data, &model.Name, &model.Meta, &model.SubTemplates); err != nil {
		return nil, false, err
	}

	t, err := templates.engine.Convert(model)
	if err != nil {
		return nil, false, err
	}

	return t, true, nil
}

func (templates *Templates) Query(query *crud.Query) ([]template.Template, error) {
	sql := `
		SELECT data, name, meta, subtemplates FROM templates
		ORDER BY name
	`

	if query != nil {
		if query.Offset != nil {
			sql += "\n"
			sql += fmt.Sprintf("OFFSET %d", *query.Offset)
		}
		if query.Limit != nil {
			sql += "\n"
			sql += fmt.Sprintf("LIMIT %d", *query.Limit)
		}
	}

	result, err := templates.conn.Query(context.TODO(), sql)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var output []template.Template
	for result.Next() {
		var model TemplateModel
		if err := result.Scan(&model.Data, &model.Name, &model.Meta, &model.SubTemplates); err != nil {
			return nil, err
		}

		t, err := templates.engine.Convert(model)
		if err != nil {
			return nil, err
		}

		output = append(output, t)
	}

	return output, nil
}

func (templates *Templates) Create(temp template.Template) (template.Template, error) {
	sql := `
		INSERT INTO templates (data, name, meta, subtemplates) VALUES ($1, $2, $3, $4) 
	`
	model := template2model(temp)
	if _, err := templates.conn.Exec(context.TODO(), sql, model.Data, model.Name, model.Meta, model.SubTemplates); err != nil {
		return nil, err
	}
	return templates.engine.Convert(model)
}

func (templates *Templates) Update(temp template.Template) error {
	sql := `
		UPDATE templates
		SET data=$2,
		    meta=$3,
		    subtemplates=$4
		WHERE name = $1
	`

	model := template2model(temp)

	result, err := templates.conn.Query(context.TODO(), sql, model.Name, model.Data, model.Meta, model.SubTemplates)
	if err != nil {
		return err
	}
	defer result.Close()
	return nil
}

func (templates *Templates) Delete(id string) error {
	sql := `
		DELETE FROM templates
		WHERE name = $1
	`
	if _, err := templates.conn.Exec(context.TODO(), sql, id); err != nil {
		return err
	}
	return nil
}

func template2model(t template.Template) TemplateModel {
	model := TemplateModel{
		Name:         t.Name(),
		Data:         t.Raw(),
		Meta:         t.Meta(),
		SubTemplates: map[string]TemplateModel{},
	}

	for name, subTemplate := range t.SubTemplates() {
		data := template2model(subTemplate)
		if data.Name == "" {
			data.Name = name
		}
		model.SubTemplates[name] = data
	}
	return model
}

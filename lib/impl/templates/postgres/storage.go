package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"yubin/common/data/kv"
	"yubin/lib/impl/templates"
	"yubin/src/template"
)

type Storage struct {
	engine templates.Engine
	conn   *pgx.Conn
}

func New(engine templates.Engine, conn *pgx.Conn) *Storage {

}

func (storage *Storage) Init(ctx context.Context) error {
	sql := `
		CREATE TABLE templates (
		    data varchar(65000),
		    name varchar(128) PRIMARY KEY,
		    meta json,
		    subtemplates json
		)
	`

	_, err := storage.conn.Exec(ctx, sql)
	return err
}

func (storage *Storage) Get(ctx context.Context, name string) (template.Template, bool, error) {
	sql := `
		SELECT data, name, meta, subtemplates FROM templates
		WHERE name=$1
	`

	result, err := storage.conn.Query(ctx, sql, name)
	if err != nil {
		return nil, false, err
	}
	defer result.Close()
	if !result.Next() {
		return nil, false, nil
	}

	var model templates.Model
	if err := result.Scan(&model.Data, &model.Name, &model.Meta, &model.SubTemplates); err != nil {
		return nil, false, err
	}

	t, err := storage.engine.Convert(model)
	if err != nil {
		return nil, false, err
	}

	return t, true, nil
}

func (storage *Storage) Query(ctx context.Context, order kv.Order, offset int, limit int) ([]template.Template, error) {
	sql := `
		SELECT data, name, meta, subtemplates FROM templates
		ORDER BY name
	`

	if order == kv.ASC {
		sql += " "
		sql += "ASC"
	}

	if offset >= 0 {
		sql += "\n"
		sql += fmt.Sprintf("OFFSET %d", offset)
	}
	if limit >= 0 {
		sql += "\n"
		sql += fmt.Sprintf("LIMIT %d", limit)
	}

	result, err := storage.conn.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var output []template.Template
	for result.Next() {
		var model templates.Model
		if err := result.Scan(&model.Data, &model.Name, &model.Meta, &model.SubTemplates); err != nil {
			return nil, err
		}

		t, err := storage.engine.Convert(model)
		if err != nil {
			return nil, err
		}

		output = append(output, t)
	}

	return output, nil
}

func (storage *Storage) Create(ctx context.Context, temp template.Template) (template.Template, error) {
	sql := `
		INSERT INTO templates (data, name, meta, subtemplates) VALUES ($1, $2, $3, $4) 
	`
	model := template2model(temp)
	if _, err := storage.conn.Exec(ctx, sql, model.Data, model.Name, model.Meta, model.SubTemplates); err != nil {
		return nil, err
	}
	return storage.engine.Convert(model)
}

func (storage *Storage) Update(ctx context.Context, temp template.Template) error {
	sql := `
		UPDATE templates
		SET data=$2,
		    meta=$3,
		    subtemplates=$4
		WHERE name = $1
	`

	model := template2model(temp)

	result, err := storage.conn.Query(ctx, sql, model.Name, model.Data, model.Meta, model.SubTemplates)
	if err != nil {
		return err
	}
	defer result.Close()
	return nil
}

func (storage *Storage) Delete(ctx context.Context, id string) error {
	sql := `
		DELETE FROM templates
		WHERE name = $1
	`
	if _, err := storage.conn.Exec(ctx, sql, id); err != nil {
		return err
	}
	return nil
}

func template2model(t template.Template) templates.Model {
	model := templates.Model{
		Name:         t.Name(),
		Data:         t.Raw(),
		Meta:         t.Meta(),
		SubTemplates: map[string]templates.Model{},
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

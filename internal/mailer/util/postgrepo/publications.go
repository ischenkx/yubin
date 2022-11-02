package postgrepo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"smtp-client/internal/mailer"
	"smtp-client/pkg/data/crud"
)

type Publications struct {
	conn *pgx.Conn
}

func NewPublications(conn *pgx.Conn) *Publications {
	return &Publications{conn: conn}
}

func (s *Publications) Init(ctx context.Context) error {
	sql := `
		CREATE TABLE publications (
		    id varchar(256) PRIMARY KEY,
		    topics text[],
		    users text[],
		    source varchar(256),
		    template varchar(256),
		    at timestamptz,
		    meta json
		    )
	`
	_, err := s.conn.Exec(ctx, sql)
	return err
}

func (s *Publications) Get(id string) (mailer.Publication, bool, error) {
	sql := `
		SELECT id, topics, users, source, template, at, meta FROM publications
		WHERE id = $1
	`
	result, err := s.conn.Query(context.TODO(), sql, id)
	if err != nil {
		return mailer.Publication{}, false, err
	}
	defer result.Close()
	if !result.Next() {
		return mailer.Publication{}, false, nil
	}
	var publication mailer.Publication
	if err := result.Scan(
		&publication.ID,
		&publication.Info.SendOptions.Topics,
		&publication.Info.SendOptions.Users,
		&publication.Info.SendOptions.SourceID,
		&publication.Info.SendOptions.Template,
		&publication.Info.At,
		&publication.Info.Meta,
	); err != nil {
		return mailer.Publication{}, false, err
	}
	return publication, true, nil
}

func (s *Publications) Query(query *crud.Query) ([]mailer.Publication, error) {
	sql := `
		SELECT id, topics, users, source, template, at, meta
		FROM publications
		ORDER BY id
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

	result, err := s.conn.Query(context.TODO(), sql)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var publications []mailer.Publication
	for result.Next() {
		var publication mailer.Publication
		if err := result.Scan(
			&publication.ID,
			&publication.Info.SendOptions.Topics,
			&publication.Info.SendOptions.Users,
			&publication.Info.SendOptions.SourceID,
			&publication.Info.SendOptions.Template,
			&publication.Info.At,
			&publication.Info.Meta,
		); err != nil {
			return nil, err
		}
		publications = append(publications, publication)
	}
	return publications, nil
}

func (s *Publications) Create(publication mailer.Publication) (mailer.Publication, error) {
	sql := `
		INSERT INTO publications
		    (id, topics, users, source, template, at, meta) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	if _, err := s.conn.Exec(context.TODO(), sql,
		publication.ID,
		publication.Info.SendOptions.Topics,
		publication.Info.SendOptions.Users,
		publication.Info.SendOptions.SourceID,
		publication.Info.SendOptions.Template,
		publication.Info.At,
		publication.Info.Meta,
	); err != nil {
		return mailer.Publication{}, err
	}
	return publication, nil
}

func (s *Publications) Update(publication mailer.Publication) error {
	sql := `
		UPDATE publications
		SET 
		    topics=$2, 
		    users=$3, 
		    source=$4, 
		    template=$5, 
		    at=$7, 
		    meta=$8
		WHERE id=$1
	`
	if _, err := s.conn.Exec(context.TODO(), sql,
		publication.ID,
		publication.Info.SendOptions.Topics,
		publication.Info.SendOptions.Users,
		publication.Info.SendOptions.SourceID,
		publication.Info.SendOptions.Template,
		publication.Info.At,
		publication.Info.Meta,
	); err != nil {
		return err
	}
	return nil
}

func (s *Publications) Delete(id string) error {
	sql := `
		DELETE FROM publications
		WHERE id=$1
	`
	if _, err := s.conn.Exec(context.TODO(), sql, id); err != nil {
		return err
	}
	return nil
}

package postgrepo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"smtp-client/internal/mailer"
	"smtp-client/pkg/data/crud"
)

type Sources struct {
	conn *pgx.Conn
}

func NewSources(conn *pgx.Conn) *Sources {
	return &Sources{conn: conn}
}

func (s *Sources) InitTable() error {
	sql := `
		CREATE TABLE sources (
		    name varchar(1024) PRIMARY KEY,
		    address varchar(1024),
		    password varchar(1024),
		    host varchar(1024),
		    port int
		    )
	`
	_, err := s.conn.Exec(context.TODO(), sql)
	return err
}

func (s *Sources) Get(id string) (mailer.NamedSource, bool, error) {
	sql := `
		SELECT name, address, password, host, port FROM sources
		WHERE name = $1
	`
	result, err := s.conn.Query(context.TODO(), sql, id)
	if err != nil {
		return mailer.NamedSource{}, false, err
	}
	defer result.Close()
	if !result.Next() {
		return mailer.NamedSource{}, false, nil
	}
	var source mailer.NamedSource
	if err := result.Scan(&source.Name,
		&source.Address,
		&source.Password,
		&source.Host,
		&source.Port); err != nil {
		return mailer.NamedSource{}, false, err
	}
	return source, true, nil
}

func (s *Sources) Query(query *crud.Query) ([]mailer.NamedSource, error) {
	sql := `
		SELECT name, address, password, host, port FROM sources
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

	result, err := s.conn.Query(context.TODO(), sql)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var sources []mailer.NamedSource
	for result.Next() {
		var source mailer.NamedSource
		if err := result.Scan(&source.Name, &source.Address, &source.Password, &source.Host, &source.Port); err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}
	return sources, nil
}

func (s *Sources) Create(t mailer.NamedSource) (mailer.NamedSource, error) {
	sql := `
		INSERT INTO sources (name, address, password, host, port) VALUES ($1, $2, $3, $4, $5)
	`
	if _, err := s.conn.Exec(context.TODO(), sql, t.Name, t.Address, t.Password, t.Host, t.Port); err != nil {
		return mailer.NamedSource{}, err
	}
	return t, nil
}

func (s *Sources) Update(t mailer.NamedSource) error {
	sql := `
		UPDATE sources
		SET address=$2, password=$3, host=$4, port=$5
		WHERE name=$1
	`
	if _, err := s.conn.Exec(context.TODO(), sql, t.Name, t.Address, t.Password, t.Host, t.Port); err != nil {
		return err
	}
	return nil
}

func (s *Sources) Delete(id string) error {
	sql := `
		DELETE FROM sources
		WHERE name=$1
	`
	if _, err := s.conn.Exec(context.TODO(), sql, id); err != nil {
		return err
	}
	return nil
}

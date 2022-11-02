package postgrepo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"smtp-client/internal/mailer"
	"smtp-client/pkg/data/crud"
)

type Reports struct {
	conn *pgx.Conn
}

func NewReports(conn *pgx.Conn) *Reports {
	return &Reports{conn: conn}
}

func (s *Reports) Init(ctx context.Context) error {
	sql := `
		CREATE TABLE reports (
		    publication_id varchar(1024) PRIMARY KEY,
		    status varchar(256),
		    failed text[],
		    ok text[]
		    )
	`
	_, err := s.conn.Exec(ctx, sql)
	return err
}

func (s *Reports) Get(id string) (mailer.Report, bool, error) {
	sql := `
		SELECT publication_id, failed, ok, status FROM reports
		WHERE publication_id = $1
	`
	result, err := s.conn.Query(context.TODO(), sql, id)
	if err != nil {
		return mailer.Report{}, false, err
	}
	defer result.Close()
	if !result.Next() {
		return mailer.Report{}, false, nil
	}
	var report mailer.Report
	if err := result.Scan(&report.PublicationID, &report.Failed, &report.OK, &report.Status); err != nil {
		return mailer.Report{}, false, err
	}
	return report, true, nil
}

func (s *Reports) Query(query *crud.Query) ([]mailer.Report, error) {
	sql := `
		SELECT publication_id, failed, ok, status FROM reports
		ORDER BY publication_id
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

	var reports []mailer.Report
	for result.Next() {
		var report mailer.Report
		if err := result.Scan(&report.PublicationID, &report.Failed, &report.OK, &report.Status); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func (s *Reports) Create(t mailer.Report) (mailer.Report, error) {
	sql := `
		INSERT INTO reports (publication_id, failed, ok, status) VALUES ($1, $2, $3, $4)
	`
	if _, err := s.conn.Exec(context.TODO(), sql, t.PublicationID, t.Failed, t.OK, t.Status); err != nil {
		return mailer.Report{}, err
	}
	return t, nil
}

func (s *Reports) Update(t mailer.Report) error {
	sql := `
		UPDATE reports
		SET failed=$2, ok=$3, status=$4
		WHERE publication_id=$1
	`
	if _, err := s.conn.Exec(context.TODO(), sql, t.PublicationID, t.Failed, t.OK, t.Status); err != nil {
		return err
	}
	return nil
}

func (s *Reports) Delete(id string) error {
	sql := `
		DELETE FROM reports
		WHERE publication_id=$1
	`
	if _, err := s.conn.Exec(context.TODO(), sql, id); err != nil {
		return err
	}
	return nil
}

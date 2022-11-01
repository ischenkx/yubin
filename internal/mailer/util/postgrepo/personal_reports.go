package postgrepo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"smtp-client/internal/mailer"
	"smtp-client/pkg/data/crud"
)

type PersonalReports struct {
	conn *pgx.Conn
}

func NewPersonalReports(conn *pgx.Conn) *PersonalReports {
	return &PersonalReports{conn: conn}
}

func (s *PersonalReports) InitTable() error {
	sql := `
		CREATE TABLE personal_reports (
		    publication_id varchar(1024),
		    user_id varchar(1024),
		    meta json,
		    status varchar(1024),
		    primary key (publication_id, user_id)
		    )
	`
	_, err := s.conn.Exec(context.TODO(), sql)
	return err
}

func (s *PersonalReports) Get(key crud.PairKey[string, string]) (mailer.PersonalReport, bool, error) {
	sql := `
		SELECT publication_id, user_id, meta, status FROM personal_reports
		WHERE publication_id = $1 and user_id = $2
	`
	result, err := s.conn.Query(context.TODO(), sql, key.Item1, key.Item2)
	if err != nil {
		return mailer.PersonalReport{}, false, err
	}
	defer result.Close()
	if !result.Next() {
		return mailer.PersonalReport{}, false, nil
	}
	var report mailer.PersonalReport
	if err := result.Scan(&report.PublicationID, &report.UserID, &report.Meta, &report.Status); err != nil {
		return mailer.PersonalReport{}, false, err
	}
	return report, true, nil
}

func (s *PersonalReports) Query(query *crud.Query) ([]mailer.PersonalReport, error) {
	sql := `
		SELECT publication_id, user_id, meta, status FROM personal_reports
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

	var reports []mailer.PersonalReport
	for result.Next() {
		var report mailer.PersonalReport
		if err := result.Scan(&report.PublicationID, &report.PublicationID, &report.UserID, &report.Meta, &report.Status); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func (s *PersonalReports) Create(t mailer.PersonalReport) (mailer.PersonalReport, error) {
	sql := `
		INSERT INTO personal_reports (publication_id, user_id, meta, status) VALUES ($1, $2, $3, $4)
	`
	if _, err := s.conn.Exec(context.TODO(), sql,
		t.PublicationID,
		t.UserID,
		t.Meta,
		t.Status); err != nil {
		return mailer.PersonalReport{}, err
	}
	return t, nil
}

func (s *PersonalReports) Update(t mailer.PersonalReport) error {
	sql := `
		UPDATE personal_reports
		SET status=$3, meta=$4
		WHERE publication_id=$1 and user_id=$2
	`
	if _, err := s.conn.Exec(context.TODO(), sql,
		t.PublicationID,
		t.UserID,
		t.Status,
		t.Meta); err != nil {
		return err
	}
	return nil
}

func (s *PersonalReports) Delete(id crud.PairKey[string, string]) error {
	sql := `
		DELETE FROM personal_reports
		WHERE publication_id = $1 and user_id = $2
	`
	if _, err := s.conn.Exec(context.TODO(), sql, id.Item1, id.Item2); err != nil {
		return err
	}
	return nil
}

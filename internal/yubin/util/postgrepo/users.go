package postgrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"smtp-client/internal/yubin/user"
	"smtp-client/pkg/data/crud"
)

type Users struct {
	conn *pgx.Conn
}

func (users *Users) Init(ctx context.Context) error {
	sql := `
		CREATE TABLE users (
				id varchar(256) PRIMARY KEY,
				email varchar(256),
				name varchar(256),
				surname varchar(256),
				meta json
			);
	`
	_, err := users.conn.Exec(ctx, sql)
	return err
}

func (users *Users) DropTable(ctx context.Context) error {
	sql := `
		DROP TABLE users;
	`
	_, err := users.conn.Exec(ctx, sql)
	return err
}

func (users *Users) Get(id string) (user.User, bool, error) {
	sql := `
		SELECT id, email, name, surname, meta FROM users
		WHERE id=$1
	`

	result, err := users.conn.Query(context.TODO(), sql, id)
	if err != nil {
		return user.User{}, false, err
	}
	defer result.Close()
	if !result.Next() {
		return user.User{}, false, nil
	}

	var u user.User
	if err := result.Scan(&u.ID, &u.Email, &u.Name, &u.Surname, &u.Meta); err != nil {
		return user.User{}, false, err
	}
	return u, true, nil
}

func (users *Users) Query(query *crud.Query) ([]user.User, error) {
	sql := `
		SELECT id, email, name, surname, meta FROM users
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

	result, err := users.conn.Query(context.TODO(), sql)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var output []user.User
	for result.Next() {
		var u user.User
		if err := result.Scan(&u.ID, &u.Email, &u.Name, &u.Surname, &u.Meta); err != nil {
			return nil, err
		}
		output = append(output, u)
	}
	return output, nil
}

func (users *Users) Create(u user.User) (user.User, error) {
	u.ID = uuid.New().String()

	sql := `
		INSERT INTO users (id, email, name, surname, meta)
		VALUES ($1, $2, $3, $4, $5::json)
		RETURNING id, email, name, surname, meta
	`

	meta, err := json.Marshal(u.Meta)
	if err != nil {
		return user.User{}, err
	}

	result, err := users.conn.Query(context.TODO(), sql, u.ID, u.Email, u.Name, u.Surname, string(meta))
	if err != nil {
		return user.User{}, err
	}
	defer result.Close()

	return u, nil
}

func (users *Users) Update(u user.User) error {
	meta, err := json.Marshal(u.Meta)
	if err != nil {
		return err
	}
	sql := `
		UPDATE users
		SET name=$2, surname=$3, email=$4, meta=$5::json
		WHERE id=$1
	`
	_, err = users.conn.Exec(context.TODO(), sql, u.ID, u.Name, u.Surname, u.Email, string(meta))
	return err
}

func (users *Users) Delete(id string) error {
	sql := `
		DELETE FROM users
		WHERE id=$1
	`
	_, err := users.conn.Exec(context.TODO(), sql, id)
	return err
}

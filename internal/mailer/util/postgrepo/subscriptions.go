package postgrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/openpgp/errors"
	"smtp-client/internal/mailer/subscription"
	"smtp-client/pkg/data/crud"
)

type Subscriptions struct {
	conn *pgx.Conn
}

func (s *Subscriptions) InitTable(ctx context.Context) error {
	sql := `
		CREATE TABLE subscriptions (
				subscriber varchar(256),
				topic varchar(256),
				created_at timestamptz,
				meta json
			);
		ALTER TABLE subscriptions
		ADD CONSTRAINT subscriber_topic UNIQUE (subscriber, topic);
	`

	_, err := s.conn.Exec(ctx, sql)

	return err
}

func (s *Subscriptions) DropTable(ctx context.Context) error {
	sql := `
		DROP TABLE subscriptions;
	`
	_, err := s.conn.Exec(ctx, sql)
	return err
}

func (s *Subscriptions) Subscriptions(id string) ([]subscription.Subscription, error) {
	sql := `
			SELECT subscriber, topic, created_at, meta FROM subscriptions
			WHERE subscriber=$1
	`

	result, err := s.conn.Query(context.TODO(), sql, id)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var subscriptions []subscription.Subscription
	for result.Next() {
		var sub subscription.Subscription
		if err := result.Scan(&sub.Subscriber, &sub.Topic, &sub.CreatedAt, &sub.Meta); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}
	return subscriptions, nil
}

func (s *Subscriptions) Subscription(id, topic string) (subscription.Subscription, bool, error) {
	sql := `
			SELECT subscriber, topic, created_at, meta FROM subscriptions
			WHERE subscriber=$1 and topic=$2
	`

	result, err := s.conn.Query(context.TODO(), sql, id)
	if err != nil {
		return subscription.Subscription{}, false, err
	}
	defer result.Close()

	if !result.Next() {
		return subscription.Subscription{}, false, nil
	}

	var sub subscription.Subscription
	if err := result.Scan(&sub.Subscriber, &sub.Topic, &sub.CreatedAt, &sub.Meta); err != nil {
		return subscription.Subscription{}, false, err
	}

	return sub, true, nil
}

func (s *Subscriptions) Subscribers(topic string, query *crud.Query) ([]string, error) {
	sql := `
		SELECT subscriber FROM subscriptions
		WHERE topic=$1
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

	result, err := s.conn.Query(context.TODO(), sql, topic)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var subscribers []string
	for result.Next() {
		var subscriber string
		if err := result.Scan(&subscriber); err != nil {
			return nil, err
		}
		subscribers = append(subscribers, subscriber)
	}
	return subscribers, nil
}

func (s *Subscriptions) Topics(query *crud.Query) ([]string, error) {
	sql := `
		SELECT DISTINCT topic FROM subscriptions
		ORDER BY topic ASC
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

	var topics []string
	for result.Next() {
		var topic string
		if err := result.Scan(&topic); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, nil
}

func (s *Subscriptions) DeleteTopic(topic string) error {
	sql := `
		DELETE FROM subscriptions
		WHERE topic=$1
	`
	_, err := s.conn.Exec(context.TODO(), sql, topic)
	return err
}

func (s *Subscriptions) Subscribe(id, topic string) (subscription.Subscription, error) {
	sql := `
		INSERT INTO subscriptions (subscriber, topic, created_at, meta)
			VALUES (
				$1,
				$2,
				now(),
				'{}'::json
		   )
		RETURNING subscriber, topic, created_at, meta
	`

	res, err := s.conn.Query(context.TODO(), sql, id, topic)
	if err != nil {
		return subscription.Subscription{}, err
	}
	defer res.Close()

	if !res.Next() {
		return subscription.Subscription{}, errors.SignatureError("not found")
	}

	var sub subscription.Subscription
	if err := res.Scan(&sub.Subscriber, &sub.Topic, &sub.CreatedAt, &sub.Meta); err != nil {
		return subscription.Subscription{}, err
	}

	return sub, nil
}

func (s *Subscriptions) Unsubscribe(id, topic string) error {
	sql := `
		DELETE FROM subscriptions
		WHERE subscriber=$1 and topic=$2
	`

	_, err := s.conn.Exec(context.TODO(), sql, id, topic)
	return err
}

func (s *Subscriptions) UnsubscribeAll(id string) error {
	sql := `
		DELETE FROM subscriptions
		WHERE subscriber=$1
	`
	_, err := s.conn.Exec(context.TODO(), sql, id)
	return err
}

func (s *Subscriptions) Update(subscription subscription.Subscription) error {
	meta, err := json.Marshal(subscription.Meta)
	if err != nil {
		return err
	}
	sql := `
		UPDATE subscriptions
		SET meta=$3::json
		WHERE subscriber=$1 and topic=$2
	`
	_, err = s.conn.Exec(context.TODO(), sql, subscription.Subscriber, subscription.Topic, string(meta))
	return err
}

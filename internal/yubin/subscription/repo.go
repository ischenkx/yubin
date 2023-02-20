package subscription

import "smtp-client/pkg/data/crud"

type Repo interface {
	Subscription(id, topic string) (Subscription, bool, error)
	Subscriptions(id string) ([]Subscription, error)
	Subscribers(topic string, query *crud.Query) ([]string, error)
	Topics(query *crud.Query) ([]string, error)
	DeleteTopic(topic string) error
	Subscribe(id, topic string) (Subscription, error)
	Unsubscribe(id, topic string) error
	UnsubscribeAll(id string) error
	Update(Subscription) error
}

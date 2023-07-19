package subscription

import "context"

type Repo interface {
	Subscription(ctx context.Context, id, topic string) (Subscription, bool, error)
	Subscriptions(ctx context.Context, id string) ([]Subscription, error)
	Subscribers(ctx context.Context, topic string) ([]string, error)
	Topics(ctx context.Context) ([]string, error)
	DeleteTopic(ctx context.Context, topic string) error
	Subscribe(ctx context.Context, id, topic string) (Subscription, error)
	Unsubscribe(ctx context.Context, id, topic string) error
	UnsubscribeAll(ctx context.Context, id string) error
	Update(ctx context.Context, sub Subscription) error
}

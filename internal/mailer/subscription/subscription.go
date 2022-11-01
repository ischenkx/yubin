package subscription

import "time"

type Subscription struct {
	Subscriber string
	Topic      string
	CreatedAt  time.Time
	Meta       map[string]any
}

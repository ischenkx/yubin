package redistat

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"smtp-client/internal/mailer/util/plugins/viewstat"
	"smtp-client/pkg/channel"
	"smtp-client/pkg/pubsub"
	"strings"
)

type Visitor struct {
	urls    []string
	channel string
	client  redis.UniversalClient
	pubsub  *pubsub.PubSub[viewstat.Identifier]
}

func New(urls []string, channel string, client redis.UniversalClient) *Visitor {
	return &Visitor{
		urls:    urls,
		channel: channel,
		client:  client,
		pubsub:  pubsub.New[viewstat.Identifier](),
	}
}

func (v *Visitor) Run(ctx context.Context) {
	ps := v.client.Subscribe(ctx, v.channel)
	messages := ps.Channel()
	defer ps.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case mes := <-messages:
			id, ok := v.parsePayload(mes.Payload)
			if !ok {
				log.Println("failed to parse payload:", mes.Payload)
				continue
			}
			v.pubsub.Publish(id)
		}
	}
}

func (v *Visitor) GenerateLink(identifier viewstat.Identifier) (string, error) {
	if len(v.urls) == 0 {
		return "", errors.New("no urls provided")
	}
	url := v.urls[rand.Intn(len(v.urls))]
	return fmt.Sprintf("%s/%s:%s", url, identifier.Publication, identifier.User), nil
}

func (v *Visitor) Visits() channel.Handle[viewstat.Identifier] {
	return v.pubsub.Subscribe()
}

func (v *Visitor) parsePayload(payload string) (viewstat.Identifier, bool) {
	parts := strings.Split(payload, ":")
	if len(parts) != 2 {
		return viewstat.Identifier{}, false
	}
	return viewstat.Identifier{
		Publication: parts[0],
		User:        parts[1],
	}, true
}

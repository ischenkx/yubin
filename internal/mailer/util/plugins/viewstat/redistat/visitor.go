package redistat

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"smtp-client/internal/mailer/util/plugins/viewstat"
	"strings"
	"sync"
)

type Visitor struct {
	urls    []string
	channel string
	client  redis.UniversalClient
	handles []chan viewstat.Identifier
	mu      sync.RWMutex
}

func New(urls []string, channel string, client redis.UniversalClient) *Visitor {
	return &Visitor{
		urls:    urls,
		channel: channel,
		client:  client,
	}
}

func (v *Visitor) Run(ctx context.Context) {
	pubsub := v.client.Subscribe(ctx, v.channel)
	messages := pubsub.Channel()
	defer pubsub.Close()

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
			v.publish(id)
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

func (v *Visitor) Visits() viewstat.Handle {
	return v.newHandle()
}

func (v *Visitor) newHandle() *Handle {
	v.mu.Lock()
	defer v.mu.Unlock()
	c := make(chan viewstat.Identifier, 2048)
	v.handles = append(v.handles, c)
	return &Handle{
		channel: c,
		visitor: v,
	}
}

func (v *Visitor) removeHandle(c chan viewstat.Identifier) {
	v.mu.Lock()
	defer v.mu.Unlock()
	for i, candidate := range v.handles {
		if candidate == c {
			v.handles[i] = v.handles[len(v.handles)-1]
			v.handles = v.handles[:len(v.handles)-1]
			break
		}
	}
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

func (v *Visitor) publish(id viewstat.Identifier) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	for _, c := range v.handles {
		select {
		case c <- id:
		default:
		}
	}
}

type Handle struct {
	channel chan viewstat.Identifier
	visitor *Visitor
}

func (h Handle) Chan() <-chan viewstat.Identifier {
	return h.channel
}

func (h Handle) Close() {
	close(h.channel)
	h.visitor.removeHandle(h.channel)
}

package piper

import (
	"context"
	"fmt"
	"log"
	"yubin/common/data/queue"
	yubin "yubin/src"
)

type Piper struct {
	queue queue.Queue[string]
	yubin *yubin.Yubin
}

func (piper *Piper) Schedule(ctx context.Context, id string) error {
	return piper.queue.Push(context.WithValue(ctx, "yubin", piper.yubin), id)
}

func (piper *Piper) Process(ctx context.Context) error {
	channel, err := piper.queue.Read(ctx)
	if err != nil {
		return fmt.Errorf("failed to read from the queue: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case id := <-channel:
			if _, err := piper.yubin.Deliver(ctx, id); err != nil {
				log.Println("failed to deliver the publication:", err)
			}
		}
	}
}

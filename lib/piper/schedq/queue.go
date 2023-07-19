package schedq

import (
	"context"
	"fmt"
	"log"
	"time"
	"yubin/common/data/queue"
	"yubin/common/scheduler"
	yubin "yubin/src"
)

type Queue struct {
	real      queue.Queue[string]
	scheduler scheduler.Scheduler[string]
}

func (queue *Queue) Read(ctx context.Context) (<-chan string, error) {
	return queue.real.Read(ctx)
}

func (queue *Queue) Push(ctx context.Context, id string) error {
	if val := ctx.Value("yubin"); val != nil {
		yub := val.(*yubin.Yubin)

		pub, err := yub.Publications().Get(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get the publication(id='%s'): %s", id, err)
		}

		if pub.Properties == nil {
			pub.Properties = map[string]any{}
		}

		at, ok := pub.Properties["at"]
		if ok {
			unixTime, ok := at.(int64)
			if !ok {
				return fmt.Errorf("failed to cast '%s' to unix time", at)
			}

			if err := queue.scheduler.Schedule(ctx, id, time.Unix(unixTime, 0)); err != nil {
				return fmt.Errorf("failed to schedule: %s", err)
			}
		}
	}

	if err := queue.real.Push(ctx, id); err != nil {
		return fmt.Errorf("failed to push to the real queue: %s", err)
	}

	return nil
}

func (queue *Queue) Process(ctx context.Context) error {
	ready := queue.scheduler.Ready(ctx)
	for id := range ready {
		if err := queue.real.Push(ctx, id); err != nil {
			log.Println("failed to push to the real queue:", err)
		}
	}

	return nil
}

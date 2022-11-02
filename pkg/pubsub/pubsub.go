package pubsub

import (
	"sync"
)

type PubSub[T any] struct {
	seq     int64
	handles map[int64]*Handle[T]
	mu      sync.RWMutex
}

func New[T any]() *PubSub[T] {
	return &PubSub[T]{
		seq:     0,
		handles: map[int64]*Handle[T]{},
	}
}

func (p *PubSub[T]) Publish(item T) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, handle := range p.handles {
		select {
		case handle.channel <- item:
		default:
		}
	}
}

func (p *PubSub[T]) Subscribe() *Handle[T] {
	p.mu.Lock()
	defer p.mu.Unlock()
	for {
		if _, ok := p.handles[p.seq]; !ok {
			break
		}
		p.seq++
	}
	handle := newHandle(p, p.seq, 2048)
	p.handles[p.seq] = handle
	return handle
}

func (p *PubSub[T]) remove(id int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.handles, id)
}



package events

import "sync"

type Bus struct {
	mu          sync.RWMutex
	subscribers []chan Event
	bufferSize  int
}

func NewBus(bufferSize int) *Bus {
	return &Bus{
		subscribers: make([]chan Event, 0),
		bufferSize:  bufferSize,
	}
}

// Subscribe add a new consumer (TUI, GUI, ...)
func (b *Bus) Subscribe() <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Event, b.bufferSize)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

// Publish stream events without blocking the main tread
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subscribers {
		select {
		case ch <- e:
		default:
			// Drop event if consumer is too slow
		}
	}
}

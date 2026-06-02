package events

import (
	"encoding/json"

	"github.com/bellapacx/kids-utopia/pkg/sqs"
)

type Handler func(Event)

type Bus struct {
	handlers map[EventType][]Handler
	queue    *sqs.Client
}

func NewBus(queue *sqs.Client) *Bus {
	return &Bus{
		handlers: make(map[EventType][]Handler),
		queue:    queue,
	}
}
func (b *Bus) Subscribe(eventType EventType, h Handler) {
	b.handlers[eventType] = append(b.handlers[eventType], h)
}
func (b *Bus) Publish(e Event) {

	// 1. realtime (streak)
	if hs, ok := b.handlers[e.Type]; ok {
		for _, h := range hs {
			go h(e)
		}
	}

	// 2. durable (SQS)
	if b.queue != nil {

		data, err := json.Marshal(e)
		if err != nil {
			return
		}

		_ = b.queue.Send(string(data))
	}
}
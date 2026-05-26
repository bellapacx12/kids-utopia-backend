package events

type Handler func(Event)

type Bus struct {
	handlers map[EventType][]Handler
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[EventType][]Handler),
	}
}
func (b *Bus) Subscribe(eventType EventType, h Handler) {
	b.handlers[eventType] = append(b.handlers[eventType], h)
}
func (b *Bus) Publish(e Event) {
	if hs, ok := b.handlers[e.Type]; ok {
		for _, h := range hs {
			go h(e) // async, non-blocking
		}
	}
}
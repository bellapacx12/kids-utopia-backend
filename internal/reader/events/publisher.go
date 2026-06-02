package events

type Publisher interface {
	Publish(event Event) error
}
package events

type EventType string

const (
	ProgressUpdated EventType = "progress.updated"
	SessionStarted  EventType = "session.started"
	SessionEnded    EventType = "session.ended"
)

type Event struct {
	Type    EventType
	ChildID string
	BookID  string
	Page    int
}
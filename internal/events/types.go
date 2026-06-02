package events

// =========================
// EVENT TYPES (single source of truth)
// =========================

type EventType string

const (
	BookUploaded    EventType = "book.uploaded"
	ProgressUpdated EventType = "progress.updated"
	SessionStarted  EventType = "session.started"
	SessionEnded    EventType = "session.ended"
)
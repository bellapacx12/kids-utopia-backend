package events

type BookUploadedEvent struct {
	BookID   string `json:"book_id"`
	ObjectKey string `json:"object_key"`
	Status   string `json:"status"`
}
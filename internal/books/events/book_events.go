package events

type BookUploadedEvent struct {
	 Type      string `json:"type"`
	BookID   string `json:"book_id"`
	ObjectKey string `json:"object_key"`
	Status   string `json:"status"`
}
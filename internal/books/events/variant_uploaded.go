package events

type BookVariantUploaded struct {
	Type      string `json:"type"`
	BookID    string `json:"book_id"`
	VariantID string `json:"variant_id"`
	FileURL   string `json:"file_url"`
	Language  string `json:"language"`
	ObjectKey string `json:"object_key"`
}
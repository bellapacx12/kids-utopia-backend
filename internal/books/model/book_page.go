package model

type BookPage struct {
	PageNumber int    `json:"page_number"`
	Content    string `json:"content"`
	ImageURL   string `json:"image_url"`
	AudioURL   *string `json:"audio_url,omitempty"`
}

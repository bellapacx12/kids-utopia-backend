package model

import (
	"github.com/bellapacx/kids-utopia/internal/books/dto"
	"github.com/bellapacx/kids-utopia/internal/books/model"
)

type ReaderAccess struct {
	Allowed bool `json:"allowed"`
	Preview bool `json:"preview"`
	Locked  bool `json:"locked"`
	MaxPage int  `json:"max_page"`
}

type ReaderProgress struct {
	CurrentPage    int  `json:"current_page"`
	Completed      bool `json:"completed"`
	ProgressPercent int `json:"progress_percent"`
}

type ReaderFeatures struct {
	Audio     bool `json:"audio"`
	Bookmarks bool `json:"bookmarks"`
}

type ReadingState struct {
	SessionID string           `json:"session_id"`
	Book      interface{}      `json:"book"`
	Reader    ReaderProgress   `json:"reader"`
	Access    ReaderAccess     `json:"access"`
	Features  ReaderFeatures   `json:"features"`
}
type BookResponse struct {
    Info     *model.Book        `json:"info"`
    Variants []dto.ReaderVariant `json:"variants"`
}
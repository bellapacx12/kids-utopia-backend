package repository

import (
	"context"

	progressmodel "github.com/bellapacx/kids-utopia/internal/progress/model"
)

type Repository interface {
	GetContinueReading(
		ctx context.Context,
		childID string,
		limit int,
	) ([]progressmodel.BookProgress, error)
}
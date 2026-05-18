package storage

import (
	"context"
	"mime/multipart"
)

type Storage interface {
	UploadFile(
		ctx context.Context,
		file multipart.File,
		fileName string,
	) (string, error)

	GetPresignedURL(
		ctx context.Context,
		objectName string,
	) (string, error)
GetPublicURL(objectName string) string
}
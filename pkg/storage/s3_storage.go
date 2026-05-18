package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Storage struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	publicURL  string
}

func NewS3Storage(endpoint, accessKey, secretKey, bucket string, publicURL  string) (*S3Storage, error) {

	// =========================
	// R2 / S3 COMPATIBLE CLIENT
	// =========================
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true, // IMPORTANT: R2 = HTTPS
		Region: "auto",
	})
	if err != nil {
		return nil, err
	}

	return &S3Storage{
		client:     client,
		bucketName: bucket,
		endpoint:   endpoint,
		publicURL: publicURL,
	}, nil
}
func (s *S3Storage) UploadFile(
	ctx context.Context,
	file multipart.File,
	objectName string,
) (string, error) {

	buf := new(bytes.Buffer)

	_, err := io.Copy(buf, file)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buf.Bytes())

	_, err = s.client.PutObject(
		ctx,
		s.bucketName,
		objectName,
		bytes.NewReader(buf.Bytes()),
		int64(buf.Len()),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		return "", err
	}

	// return object key (NOT full URL)
	return objectName, nil
}
func (s *S3Storage) GetPresignedURL(
	ctx context.Context,
	objectName string,
) (string, error) {

	reqParams := make(url.Values)

	presignedURL, err := s.client.PresignedGetObject(
		ctx,
		s.bucketName,
		objectName,
		time.Hour*24,
		reqParams,
	)

	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
func (s *S3Storage) GetPublicURL(key string) string {
	return fmt.Sprintf("%s/%s", s.publicURL, key)
}
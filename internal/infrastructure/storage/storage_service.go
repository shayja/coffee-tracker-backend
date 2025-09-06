package storage

import (
	"context"
	"io"
)

// StorageService defines the interface for storage operations
type StorageService interface {
	//ListBuckets(ctx context.Context) ([]map[string]any, error)
	//CreateBucket(ctx context.Context, name string, isPublic bool) error
	UploadFile(ctx context.Context, bucket, filename string, file io.Reader) (string, error)
	GenerateSignedURL(ctx context.Context, bucket, filename string, expiresInSeconds int) (string, error)
}
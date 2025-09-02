package storage

import (
	"context"
	"io"
)

type StorageService interface {
    UploadFile(ctx context.Context, bucket, filename string, file io.Reader) (string, error)
}

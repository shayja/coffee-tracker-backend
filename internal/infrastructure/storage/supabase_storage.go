package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"

	storage "github.com/supabase-community/storage-go"
)

type SupabaseStorageService struct {
	client  *storage.Client
	baseURL string
}

func NewSupabaseStorageService(client *storage.Client) *SupabaseStorageService {
	return &SupabaseStorageService{
		client:  client,
		baseURL: "https://<project-ref>.supabase.co/storage/v1/object/public",
	}
}

// UploadFile uploads a file to a Supabase Storage bucket
func (s *SupabaseStorageService) UploadFile(ctx context.Context, bucket, filename string, file io.Reader) (string, error) {
	// Upload the file to the specified bucket
	result, err := s.client.UploadFile(bucket, filename, file)
	log.Println(result)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return the public URL of the uploaded file
	return fmt.Sprintf("%s/%s/%s", s.baseURL, bucket, url.PathEscape(filename)), nil
}
package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"

	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseStorageService struct {
	storageClient  *storage_go.Client
	baseURL string
}

func NewSupabaseStorageService(storageClient *storage_go.Client, baseURL string) *SupabaseStorageService {
    return &SupabaseStorageService{
        storageClient:  storageClient,
        baseURL: baseURL,
    }
}

// UploadFile uploads a file to a Supabase Storage bucket
func (s *SupabaseStorageService) UploadFile(ctx context.Context, bucketId, filename string, file io.Reader) (string, error) {


	err22 := s.ListBuckets(ctx)
	if err22 != nil {
	log.Printf("ListBuckets: %v", err22)
	}


    _, err := s.GetBucket(ctx, bucketId)
    if err != nil {
		log.Fatalf("GetBucket: %v", err)
        /*
        _, err := s.CreateBucket(ctx, bucketId)
        if err != nil {
            return "", fmt.Errorf("failed to create bucket %s: %v", bucketId, err)
        }
        */
    }

    upsert := true
	opts := &storage_go.FileOptions{
		Upsert: &upsert,
	}

    result, err := s.storageClient.UploadOrUpdateFile(bucketId, filename, file, true, *opts)
    if err != nil {
        log.Printf("failed to upload file: %v", err)
        return "", fmt.Errorf("failed to upload file: %w", err)
    }

    log.Printf("upload result: %v", result)

    return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.baseURL, bucketId, url.PathEscape(filename)), nil
}

func (s *SupabaseStorageService) GetBucket(ctx context.Context, bucketId string) (bool, error) {
    result, err := s.storageClient.GetBucket(bucketId)
    if err != nil {
        return false, fmt.Errorf("failed to get bucket %s: %v", bucketId, err)
    }
    log.Printf("bucket details: %v", result)
    return true, nil
}


/*
func(s *SupabaseStorageService) CreateBucket(ctx context.Context, bucketId string) (bool, error) {

	result, err := s.storageClient.CreateBucket(bucketId, storage_go.BucketOptions{
		Public: true,
	})

  	log.Print(err, result)
	if err != nil {
		return false, fmt.Errorf("failed to create bucket %s: %v", bucketId, err)
	}

	return true, nil
}
*/

func (s *SupabaseStorageService) ListBuckets(ctx context.Context) error {
    buckets, err := s.storageClient.ListBuckets()
    if err != nil {
       log.Printf("list buckets failed: %v", err)
       return fmt.Errorf("list buckets failed: %v", err)
    }

	for _, bucket := range buckets {
		log.Printf("Bucket Name: %s, ID: %s\n", bucket.Name, bucket.Id)
	}
    return nil
}

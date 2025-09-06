package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
)

// SupabaseStorageService handles interactions with Supabase storage
type SupabaseStorageService struct {
	storageURL string
	apiKey     string
	client     *http.Client
}

// NewSupabaseStorageService creates a new storage service
func NewSupabaseStorageService(storageURL, apiKey string) StorageService {
	return &SupabaseStorageService{
		storageURL: storageURL,
		apiKey:     apiKey,
		client:     &http.Client{},
	}
}

// ---------------- Buckets ----------------
/*
// ListBuckets retrieves all buckets in the Supabase project
func (s *SupabaseStorageService) ListBuckets(ctx context.Context) ([]map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.storageURL+"/bucket", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	s.addAuthHeaders(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("ListBuckets response: %s", string(body))

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("list buckets failed: %s", string(body))
	}

	var buckets []map[string]any
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&buckets); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return buckets, nil
}

// CreateBucket creates a new storage bucket
func (s *SupabaseStorageService) CreateBucket(ctx context.Context, name string, isPublic bool) error {
	payload := map[string]any{
		"name":   name,
		"public": isPublic,
	}
	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", s.storageURL+"/bucket", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	s.addAuthHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create bucket failed: %s", string(body))
	}

	return nil
}
*/
// ---------------- Objects ----------------

// UploadFile uploads a file to a specified bucket and returns a signed URL
func (s *SupabaseStorageService) UploadFile(ctx context.Context, bucket, filename string, file io.Reader) (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}
	writer.Close()

	req, err := http.NewRequestWithContext(ctx,
		"POST",
		fmt.Sprintf("%s/object/%s/%s", s.storageURL, bucket, url.PathEscape(filename)),
		&buf,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	s.addAuthHeaders(req)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("UploadFile response for %s: %s", filename, string(body))
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("upload failed for %s: %s", filename, string(body))
	}

	// Generate a signed URL for the uploaded file (expires in 1 hour)
	signedURL, err := s.GenerateSignedURL(ctx, bucket, filename, 3600)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL for %s: %v", filename, err)
	}
	return signedURL, nil
}

// ---------------- Signed URL ----------------

// GenerateSignedURL creates a temporary signed URL for accessing a private file
func (s *SupabaseStorageService) GenerateSignedURL(ctx context.Context, bucket, filename string, expiresInSeconds int) (string, error) {
	payload := map[string]any{
		"expiresIn": expiresInSeconds,
	}
	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx,
		"POST",
		fmt.Sprintf("%s/object/sign/%s/%s", s.storageURL, bucket, url.PathEscape(filename)),
		bytes.NewBuffer(data),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request for %s: %v", filename, err)
	}
	s.addAuthHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request for %s: %v", filename, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("GenerateSignedURL response for %s: %s", filename, string(body))
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("signed URL failed for %s: %s", filename, string(body))
	}

	var result struct {
		SignedURL string `json:"signedURL"`
	}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response for %s: %v", filename, err)
	}

	// Ensure the signed URL is absolute by combining with storageURL correctly
	if result.SignedURL[0] == '/' {
		return s.storageURL + result.SignedURL, nil
	}
	return fmt.Sprintf("%s/%s", s.storageURL, result.SignedURL), nil
}

// ---------------- Helpers ----------------

// addAuthHeaders adds authentication headers to the request
func (s *SupabaseStorageService) addAuthHeaders(req *http.Request) {
	req.Header.Set("apikey", s.apiKey)
	req.Header.Set("Authorization", "Bearer " + s.apiKey)
}
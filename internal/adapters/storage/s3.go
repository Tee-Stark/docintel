package storage

import "context"

// S3 implements domain.FileStorage using MinIO/S3.
type S3 struct{}

func NewS3() *S3 { return &S3{} }

func (s *S3) Upload(ctx context.Context, key string, data []byte, mimeType string) error {
	panic("not implemented")
}

func (s *S3) Download(ctx context.Context, key string) ([]byte, error) {
	panic("not implemented")
}

func (s *S3) Delete(ctx context.Context, key string) error {
	panic("not implemented")
}

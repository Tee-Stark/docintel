package config

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient() (*minio.Client, string, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("MINIO_BUCKET")
	useSSLStr := os.Getenv("MINIO_USE_SSL")

	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	if bucketName == "" {
		bucketName = "docintel-documents"
	}

	useSSL, _ := strconv.ParseBool(useSSLStr)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create minio client: %w", err)
	}

	ctx := context.Background()

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check minio bucket: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, "", fmt.Errorf("failed to create minio bucket: %w", err)
		}
	}

	return client, bucketName, nil
}

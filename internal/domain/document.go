package domain

import "time"

type DocumentStatus string

const (
	DocumentStatusPending    DocumentStatus = "pending"
	DocumentStatusProcessing DocumentStatus = "processing"
	DocumentStatusReady      DocumentStatus = "ready"
	DocumentStatusFailed     DocumentStatus = "failed"
)

type Document struct {
	ID               string
	UserID           string
	Title            string
	OriginalFilename string
	MimeType         string
	SizeBytes        int64
	StorageKey       string
	Status           DocumentStatus
	ErrorMessage     *string
	PageCount        *int
	ChunkCount       int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ProcessedAt      *time.Time
}

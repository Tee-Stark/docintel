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
	ID         string
	UserID     string
	Name       string
	MimeType   string
	SizeBytes  int64
	StorageKey string
	Status     DocumentStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

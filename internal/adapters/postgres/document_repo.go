package postgres

import (
	"context"

	"docintel/internal/domain"
)

const (
	UploadDir = "./uploads"
)

func (r *Repository) CreateDocument(ctx context.Context, doc *domain.Document) error {
	query := `
		INSERT INTO documents (
			id, user_id, title, original_filename, mime_type,
			size_bytes, storage_key, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		doc.ID,
		doc.UserID,
		doc.Title,
		doc.OriginalFilename,
		doc.MimeType,
		doc.SizeBytes,
		doc.StorageKey,
	)

	return err
}

// func (r *Repository) FindByID(ctx context.Context, id string) (*domain.Document, error) {
// 	panic("not implemented")
// }

// func (r *Repository) ListByUserID(ctx context.Context, userID string) ([]*domain.Document, error) {
// 	panic("not implemented")
// }

// func (r *Repository) UpdateStatus(ctx context.Context, id string, status domain.DocumentStatus) error {
// 	panic("not implemented")
// }

// func (r *Repository) Delete(ctx context.Context, id string) error {
// 	panic("not implemented")
// }

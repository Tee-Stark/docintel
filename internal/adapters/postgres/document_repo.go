package postgres

import (
	"context"

	"docintel/internal/domain"
)

type DocumentRepo struct{}

func NewDocumentRepo() *DocumentRepo { return &DocumentRepo{} }

func (r *DocumentRepo) Create(ctx context.Context, doc *domain.Document) error {
	panic("not implemented")
}

func (r *DocumentRepo) FindByID(ctx context.Context, id string) (*domain.Document, error) {
	panic("not implemented")
}

func (r *DocumentRepo) ListByUserID(ctx context.Context, userID string) ([]*domain.Document, error) {
	panic("not implemented")
}

func (r *DocumentRepo) UpdateStatus(ctx context.Context, id string, status domain.DocumentStatus) error {
	panic("not implemented")
}

func (r *DocumentRepo) Delete(ctx context.Context, id string) error {
	panic("not implemented")
}

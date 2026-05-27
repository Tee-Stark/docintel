package postgres

import (
	"context"

	"docintel/internal/domain"
)

type ChunkRepo struct{}

func NewChunkRepo() *ChunkRepo { return &ChunkRepo{} }

func (r *ChunkRepo) CreateBatch(ctx context.Context, chunks []*domain.Chunk) error {
	panic("not implemented")
}

func (r *ChunkRepo) SearchByEmbedding(ctx context.Context, embedding []float32, topK int) ([]*domain.Chunk, error) {
	panic("not implemented")
}

func (r *ChunkRepo) DeleteByDocumentID(ctx context.Context, documentID string) error {
	panic("not implemented")
}

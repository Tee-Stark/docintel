package domain

import "time"

type Chunk struct {
	ID         string
	DocumentID string
	Content    string
	Embedding  []float32
	ChunkIndex int
	CreatedAt  time.Time
}

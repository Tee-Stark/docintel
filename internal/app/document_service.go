package app

import (
	"context"
	"docintel/internal/domain"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

const (
	UploadDir = "./uploads"
)

type DocumentService struct {
	docRepo domain.DocumentRepository
}

func NewDocumentService(docRepo domain.DocumentRepository) *DocumentService {
	return &DocumentService{docRepo: docRepo}
}

func (s *DocumentService) UploadDocument(ctx context.Context, file multipart.File, doc *domain.Document) error {
	// Implement document upload logic, e.g., save file path to database
	err := os.MkdirAll(UploadDir, os.ModePerm)
	if err != nil {
		return err
	}

	destinationPath := filepath.Join(UploadDir, doc.StorageKey)

	dst, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		return err
	}

	title := strings.TrimSuffix(doc.OriginalFilename, filepath.Ext(doc.OriginalFilename))

	doc.SizeBytes = size
	doc.Title = title

	return s.docRepo.CreateDocument(ctx, doc)
}

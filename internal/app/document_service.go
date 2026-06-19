package app

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"

	"docintel/internal/domain"
)

type DocumentService struct {
	docRepo    domain.DocumentRepository
	minio      *minio.Client
	bucketName string
}

func NewDocumentService(docRepo domain.DocumentRepository, minioClient *minio.Client, bucketName string) *DocumentService {
	return &DocumentService{
		docRepo:    docRepo,
		minio:      minioClient,
		bucketName: bucketName,
	}
}

func (s *DocumentService) UploadDocument(ctx context.Context, file multipart.File, doc *domain.Document) error {
	title := strings.TrimSuffix(doc.OriginalFilename, filepath.Ext(doc.OriginalFilename))
	objectInfo, err := s.minio.PutObject(
		ctx,
		s.bucketName,
		doc.StorageKey,
		file,
		-1,
		minio.PutObjectOptions{
			ContentType: doc.MimeType,
		},
	)
	if err != nil {
		return err
	}
	doc.SizeBytes = objectInfo.Size
	doc.Title = title

	return s.docRepo.CreateDocument(ctx, doc)
}

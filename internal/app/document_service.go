package app

import (
	"context"
	"docintel/internal/domain"

	"mime/multipart"

	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
)

const (
	UploadDir = "./uploads"
)

type DocumentService struct {
	docRepo    domain.DocumentRepository
	minio      *minio.Client
	bucketName string
}

func NewDocumentService(docRepo domain.DocumentRepository, minioClient *minio.Client,
	bucketName string) *DocumentService {
	return &DocumentService{docRepo: docRepo,
		minio:      minioClient,
		bucketName: bucketName}
}

func (s *DocumentService) UploadDocument(ctx context.Context, file multipart.File, doc *domain.Document) error {
	// Implement document upload logic, e.g., save file path to database
	// err := os.MkdirAll(UploadDir, os.ModePerm)
	// if err != nil {
	// 	return err
	// }

	// destinationPath := filepath.Join(UploadDir, doc.StorageKey)

	// dst, err := os.Create(destinationPath)
	// if err != nil {
	// 	return err
	// }
	// defer dst.Close()

	// size, err := io.Copy(dst, file)
	// if err != nil {
	// 	return err
	// }

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

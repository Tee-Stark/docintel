package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DocumentHandler struct {
	DB        *pgxpool.Pool
	UploadDir string
}

func NewDocumentHandler(db *pgxpool.Pool, uploadDir string) *DocumentHandler {
	return &DocumentHandler{
		DB:        db,
		UploadDir: uploadDir,
	}
}

func UploadDocument(w http.ResponseWriter, r *http.Request) {
	var h DocumentHandler

	ctx := r.Context()

	err := r.ParseMultipartForm(20 << 20) // 20MB max memory
	if err != nil {
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > 20*1024*1024 {
		http.Error(w, "File too large. Max size is 20MB", http.StatusBadRequest)
		return
	}

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
		http.Error(w, "Only PDF files are allowed for now", http.StatusBadRequest)
		return
	}

	documentID := uuid.New()
	storageKey := fmt.Sprintf("%s.pdf", documentID.String())

	err = os.MkdirAll(h.UploadDir, os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to prepare upload directory", http.StatusInternalServerError)
		return
	}

	destinationPath := filepath.Join(h.UploadDir, storageKey)

	dst, err := os.Create(destinationPath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	// Temporary hardcoded user ID for now.
	// Later, get this from JWT middleware.
	userID := "00000000-0000-0000-0000-000000000001"

	title := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	mimeType := header.Header.Get("Content-Type")

	query := `
		INSERT INTO documents (
			id, user_id, title, original_filename, mime_type,
			size_bytes, storage_key, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending', $8, $8)
	`

	_, err = h.DB.Exec(
		context.Background(),
		query,
		documentID,
		userID,
		title,
		header.Filename,
		mimeType,
		size,
		storageKey,
		time.Now(),
	)

	if err != nil {
		http.Error(w, "Failed to create document record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, `{
		"document_id": "%s",
		"status": "pending",
		"filename": "%s",
		"size_bytes": %d
	}`, documentID.String(), header.Filename, size)

	_ = ctx
}

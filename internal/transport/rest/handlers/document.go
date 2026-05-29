package handlers

import (
	"docintel/internal/domain"
	"docintel/internal/transport/rest/response"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type DocumentHandler struct {
	DocService domain.DocumentService
}

func NewDocumentHandler(docService domain.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		DocService: docService,
	}
}

func (h *DocumentHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := r.ParseMultipartForm(30 << 20)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid multipart form"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, fmt.Errorf("file is required"))
		return
	}
	defer file.Close()

	if header.Size > 30*1024*1024 {
		response.WriteError(w, http.StatusBadRequest, fmt.Errorf("file too large. Max size is 30MB"))
		return
	}

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
		response.WriteError(w, http.StatusBadRequest, fmt.Errorf("only PDF files are allowed for now"))
		return
	}

	documentID, err := uuid.NewV7()
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	userID := ctx.Value("userID").(string)
	mimeType := header.Header.Get("Content-Type")

	doc := &domain.Document{
		ID:               documentID,
		UserID:           userID,
		OriginalFilename: header.Filename,
		MimeType:         mimeType,
		StorageKey:       fmt.Sprintf("%s.pdf", documentID.String()),
	}

	err = h.DocService.UploadDocument(ctx, file, doc)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resp := response.Response{
		Message: "Document uploaded successfully",
		Data: map[string]interface{}{
			"document_id": documentID.String(),
			"status":      "pending",
			"filename":    header.Filename,
			"size_bytes":  header.Size,
		},
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}

package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mdarkanurl/ResumeBuilder/apps/api/internal/service"
)

type ResumeHandler struct {
	service       *service.ResumeService
	maxUploadSize int64
}

func NewResumeHandler(svc *service.ResumeService, maxUploadSize int64) *ResumeHandler {
	return &ResumeHandler{service: svc, maxUploadSize: maxUploadSize}
}

func (h *ResumeHandler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.maxUploadSize)
	if err := r.ParseMultipartForm(h.maxUploadSize); err != nil {
		var maxBytesErr *http.MaxBytesError
		switch {
		case errors.As(err, &maxBytesErr):
			http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
		case errors.Is(err, http.ErrNotMultipart):
			http.Error(w, "Content-Type must be multipart/form-data", http.StatusBadRequest)
		case errors.Is(err, http.ErrMissingBoundary):
			http.Error(w, "Missing multipart boundary", http.StatusBadRequest)
		default:
			http.Error(w, "Invalid multipart request", http.StatusBadRequest)
		}
		return
	}
	defer r.MultipartForm.RemoveAll()

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	result, err := h.service.Upload(r.Context(), header.Filename, file)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotPDF),
			errors.Is(err, service.ErrEmptyFile),
			errors.Is(err, service.ErrInvalidFilename):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Printf("upload %q failed: %v", header.Filename, err)
			http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode response: %v", err)
	}
}

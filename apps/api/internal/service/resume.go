// Package service contains the business rules of the API.
package service

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	ErrEmptyFile = errors.New("uploaded file is empty")
	ErrNotPDF = errors.New("only PDF files are allowed")
	ErrInvalidFilename = errors.New("invalid file name")
)

type ResumeRepository interface {
	Save(ctx context.Context, filename string, src io.Reader) (int64, error)
}

type UploadResult struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type ResumeService struct {
	repo ResumeRepository
}

func NewResumeService(repo ResumeRepository) *ResumeService {
	return &ResumeService{repo: repo}
}

func (s *ResumeService) Upload(ctx context.Context, filename string, src io.Reader) (UploadResult, error) {
	safeName, err := sanitizeFilename(filename)
	if err != nil {
		return UploadResult{}, err
	}

	reader := bufio.NewReader(src)
	header, err := reader.Peek(512)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, bufio.ErrBufferFull) {
		return UploadResult{}, err
	}
	if len(header) == 0 {
		return UploadResult{}, ErrEmptyFile
	}
	if http.DetectContentType(header) != "application/pdf" {
		return UploadResult{}, ErrNotPDF
	}

	size, err := s.repo.Save(ctx, safeName, reader)
	if err != nil {
		return UploadResult{}, err
	}

	return UploadResult{Filename: safeName, Size: size}, nil
}

func sanitizeFilename(filename string) (string, error) {
	name := filepath.Base(strings.TrimSpace(filename))
	if name == "" || name == "." || name == string(filepath.Separator) {
		return "", ErrInvalidFilename
	}
	return name, nil
}

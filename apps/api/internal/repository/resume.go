package repository

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileResumeRepository struct {
	dir string
}

func NewFileResumeRepository(dir string) *FileResumeRepository {
	return &FileResumeRepository{dir: dir}
}

func (r *FileResumeRepository) Save(ctx context.Context, filename string, src io.Reader) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if err := os.MkdirAll(r.dir, 0o755); err != nil {
		return 0, fmt.Errorf("create upload directory: %w", err)
	}

	dst, err := os.Create(filepath.Join(r.dir, filename))
	if err != nil {
		return 0, fmt.Errorf("create file: %w", err)
	}

	written, copyErr := io.Copy(dst, src)
	closeErr := dst.Close()

	if copyErr != nil {
		return written, fmt.Errorf("write file: %w", copyErr)
	}
	if closeErr != nil {
		return written, fmt.Errorf("close file: %w", closeErr)
	}
	return written, nil
}

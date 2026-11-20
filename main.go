package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	mux := http.NewServeMux()
	const maxUploadSize = 10 << 20

	mux.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
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

		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		buff := make([]byte, 512)
		if _, err := file.Read(buff); err != nil {
			http.Error(w, "Error reading file header", http.StatusInternalServerError)
			return
		}
		if http.DetectContentType(buff) != "application/pdf" {
			http.Error(w, "Only PDF files are allowed", http.StatusBadRequest)
			return
		}

		if _, err := file.Seek(0, io.SeekStart); err != nil {
			http.Error(w, "Error resetting file pointer", http.StatusInternalServerError)
			return
		}

		os.MkdirAll("./uploads", os.ModePerm)

		safeFilename := filepath.Base(handler.Filename)
		dst, err := os.Create(filepath.Join("./uploads", safeFilename))
		if err != nil {
			http.Error(w, "Unable to create file on server", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Successfully uploaded %s\n", safeFilename)
	})

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}

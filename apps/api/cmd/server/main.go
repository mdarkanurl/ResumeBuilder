package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mdarkanurl/ResumeBuilder/apps/api/internal/config"
	"github.com/mdarkanurl/ResumeBuilder/apps/api/internal/handler"
	"github.com/mdarkanurl/ResumeBuilder/apps/api/internal/middleware"
	"github.com/mdarkanurl/ResumeBuilder/apps/api/internal/repository"
	"github.com/mdarkanurl/ResumeBuilder/apps/api/internal/service"
)

func main() {
	cfg := config.Load()

	repo := repository.NewFileResumeRepository(cfg.UploadDir)
	resumeService := service.NewResumeService(repo)
	resumeHandler := handler.NewResumeHandler(resumeService, cfg.MaxUploadSize)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("POST /upload", resumeHandler.Upload)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           middleware.Chain(mux, middleware.Recover, middleware.Logger, middleware.CORS(cfg.AllowedOrigin)),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("API listening on http://localhost:%s (uploads: %s)", cfg.Port, cfg.UploadDir)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

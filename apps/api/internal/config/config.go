package config

import (
	"os"
	"strconv"
	"strings"
)

const defaultMaxUploadSize = 5 << 20

type Config struct {
	Port          string
	UploadDir     string
	MaxUploadSize int64
	AllowedOrigin string
	KafkaBrokers  []string
	KafkaTopic    string
}

func Load() Config {
	return Config{
		Port:          getEnv("PORT", "8080"),
		UploadDir:     getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSize: getEnvInt64("MAX_UPLOAD_SIZE", defaultMaxUploadSize),
		AllowedOrigin: getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:5173"),
		KafkaBrokers:  splitList(getEnv("KAFKA_BROKERS", "localhost:9092")),
		KafkaTopic:    getEnv("KAFKA_TOPIC", "resume-uploaded"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func splitList(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return items
}

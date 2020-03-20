package utils

import (
	"os"
	"strings"
)

// Getenv …
func Getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value != "" {
		return value
	}
	return fallback
}

package utils

import (
	"os"
	"strings"
)

func GetEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value != "" {
		return value
	}
	return fallback
}

package utils

import (
	"fmt"
	"log/slog"
	"os"
)

// SanityCheck checks that all required environment variables are set.
// If any of the required variables is not defined, it sets a default value and prints a warning message.
// If error happened during setting env variable, then logs error and exits application.
func SanityCheck(l *slog.Logger) {
	defaultEnvVars := map[string]string{
		"API_HOST":  "127.0.0.1",
		"API_PORT":  "8000",
		"DB_USER":   "postgres",
		"DB_PASSWD": "postgres",
		"DB_HOST":   "127.0.0.1",
		"DB_PORT":   "5432",
		"DB_NAME":   "fileverse",
	}

	for key, defaultValue := range defaultEnvVars {
		if os.Getenv(key) == "" {
			if err := os.Setenv(key, defaultValue); err != nil {
				l.Error(fmt.Sprintf(
					"failed to set environment variable %s to default value %s. Exiting application.",
					key,
					defaultValue,
				))
				os.Exit(1)
			}

			l.Warn(fmt.Sprintf("environment variable %s not defined. Setting to default: %s", key, defaultValue))
		}
	}
}

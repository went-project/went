package handlers

import (
	"os"
	"path/filepath"

	tmplResponses "went/internal/templates/responses"
)

// CreateGlobalErrorResponse creates the shared error response helper file.
func CreateGlobalErrorResponse(basePath string) error {
	targetDir := filepath.Join(basePath, "internal", "responses")
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return err
	}

	targetFile := filepath.Join(targetDir, "error_response.go")
	return os.WriteFile(targetFile, []byte(tmplResponses.ErrorTemplate()), 0644)
}

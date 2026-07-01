package handlers

import (
	"os"
	"path/filepath"

	tmplHTTP "went/internal/templates/http"
)

// CreateHelpers writes internal/helpers/helpers.go once per project.
func CreateHelpers(basePath string) error {
	targetDir := filepath.Join(basePath, "internal", "helpers")
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return err
	}

	targetFile := filepath.Join(targetDir, "helpers.go")
	if _, err := os.Stat(targetFile); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	return os.WriteFile(targetFile, []byte(tmplHTTP.ControllerHelperTemplate()), 0644)
}

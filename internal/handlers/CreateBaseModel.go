package handlers

import (
	"os"
	"path/filepath"

	tmplDB "went/internal/templates/database"
)

// CreateBaseModel writes database/models/base.go once per project.
func CreateBaseModel(basePath string) error {
	targetDir := filepath.Join(basePath, "database", "models")
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return err
	}

	targetFile := filepath.Join(targetDir, "base.go")
	if _, err := os.Stat(targetFile); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	return os.WriteFile(targetFile, []byte(tmplDB.BaseModelTemplate()), 0644)
}

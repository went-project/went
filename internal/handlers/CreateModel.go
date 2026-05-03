package handlers

import (
	"os"
	"path/filepath"
	"strings"

	tmplDB "went/internal/templates/database"
)

func CreateModel(name string, full bool) error {
	var (
		templateContent string
		err             error
	)

	if full {
		templateContent, err = tmplDB.ModelTemplate(name)
	} else {
		templateContent, err = tmplDB.ModelSkeletonTemplate(name)
	}
	if err != nil {
		return err
	}

	err = os.MkdirAll("database/models", os.ModePerm)
	if err != nil {
		return err
	}

	fileName := strings.ToLower(name) + ".go"
	filePath := filepath.Join("database", "models", fileName)
	return os.WriteFile(filePath, []byte(templateContent), 0644)
}

package handlers

import (
	"os"
	"path/filepath"
	"strings"

	tmplHTTP "went/internal/templates/http"
)

func CreateController(name string, full bool) error {
	var (
		templateContent string
		err             error
	)

	if full {
		templateContent, err = tmplHTTP.ControllerTemplate(name)
	} else {
		templateContent, err = tmplHTTP.ControllerSkeletonTemplate(name)
	}
	if err != nil {
		return err
	}

	err = os.MkdirAll("http/controllers", os.ModePerm)
	if err != nil {
		return err
	}

	fileName := strings.ToLower(name) + "_controller.go"
	filePath := filepath.Join("http", "controllers", fileName)
	return os.WriteFile(filePath, []byte(templateContent), 0644)
}

package handlers

import (
	"os"
	"path/filepath"
	"strings"

	tmplHTTP "went/internal/templates/http"
)

func CreateResource(name string, full bool) error {
	var (
		templateContent string
		err             error
	)

	if full {
		templateContent, err = tmplHTTP.ResourceTemplate(name)
	} else {
		templateContent, err = tmplHTTP.ResourceSkeletonTemplate(name)
	}
	if err != nil {
		return err
	}

	if err := os.MkdirAll("http/resources", os.ModePerm); err != nil {
		return err
	}

	if full {
		paginationPath := filepath.Join("http", "resources", "pagination.go")
		if _, err := os.Stat(paginationPath); os.IsNotExist(err) {
			if err := os.WriteFile(paginationPath, []byte(tmplHTTP.PaginationTemplate()), 0644); err != nil {
				return err
			}
		}

		baseResourcePath := filepath.Join("http", "resources", "base_resource.go")
		if _, err := os.Stat(baseResourcePath); os.IsNotExist(err) {
			baseContent, err := tmplHTTP.BaseResourceTemplateFromConfig()
			if err != nil {
				return err
			}
			if err := os.WriteFile(baseResourcePath, []byte(baseContent), 0644); err != nil {
				return err
			}
		}
	}

	fileName := strings.ToLower(name) + "_resource.go"
	filePath := filepath.Join("http", "resources", fileName)
	return os.WriteFile(filePath, []byte(templateContent), 0644)
}

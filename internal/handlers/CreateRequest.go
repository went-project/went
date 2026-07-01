package handlers

import (
	"os"
	"path/filepath"
	"strings"

	tmplHTTP "went/internal/templates/http"
)

func CreateRequest(name string, full bool) error {
	var (
		templateContent string
		err             error
	)

	if full {
		templateContent, err = tmplHTTP.RequestTemplate(name)
	} else {
		templateContent, err = tmplHTTP.RequestSkeletonTemplate(name)
	}
	if err != nil {
		return err
	}

	err = os.MkdirAll("http/requests", os.ModePerm)
	if err != nil {
		return err
	}

	if full {
		validatorPath := filepath.Join("http", "requests", "validator.go")
		if _, err := os.Stat(validatorPath); os.IsNotExist(err) {
			if err := os.WriteFile(validatorPath, []byte(tmplHTTP.ValidatorTemplate()), 0644); err != nil {
				return err
			}
		}
	}

	fileName := strings.ToLower(name) + "_request.go"
	filePath := filepath.Join("http", "requests", fileName)
	return os.WriteFile(filePath, []byte(templateContent), 0644)
}

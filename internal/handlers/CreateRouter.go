package handlers

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	tmplHTTP "went/internal/templates/http"
)

func CreateRouter(name string, full bool) error {

	var (
		templateContent string
		err             error
	)

	if full {
		templateContent, err = tmplHTTP.RouterTemplate(name)
	} else {
		templateContent, err = tmplHTTP.RouterSkeletonTemplate(name)
	}
	if err != nil {
		return err
	}

	tpl, err := template.New("router").Parse(templateContent)
	if err != nil {
		return err
	}

	err = os.MkdirAll("routes", os.ModePerm)
	if err != nil {
		return err
	}

	fileName := strings.ToLower(name) + "_router.go"
	f, err := os.Create(filepath.Join("routes", fileName))
	if err != nil {
		return err
	}
	defer f.Close()

	return tpl.Execute(f, nil)

}

package http

import (
	"went/internal/utils"
)

func RequestTemplate(name string) (string, error) {
	config, err := utils.ReadConfigFile()
	if err != nil {
		return "", err
	}
	appName := config["name"].(string)
	return RequestTemplateWithApp(name, appName)
}

func RequestTemplateWithApp(name string, appName string) (string, error) {
	_ = appName

	return `package requests

// ` + name + `Payload represents required fields for create operations.
// @Description ` + name + ` creation payload
type ` + name + `Payload struct {
	Email string ` + "`json:\"email\" db:\"email\" validate:\"required,email,max=255\" example:\"john.doe@example.com\"`" + `
	Name  string ` + "`json:\"name\"  db:\"name\"  validate:\"required,max=255\"      example:\"John Doe\"`" + `
}

// ` + name + `UpdatePayload represents optional fields for partial updates.
// @Description ` + name + ` update payload
type ` + name + `UpdatePayload struct {
	Email *string ` + "`json:\"email,omitempty\" db:\"email\" validate:\"omitempty,email,max=255\" example:\"john.doe@example.com\"`" + `
	Name  *string ` + "`json:\"name,omitempty\"  db:\"name\"  validate:\"omitempty,max=255\"      example:\"John Doe\"`" + `
}

func (p *` + name + `Payload) Validate() error {
	return validate.Struct(p)
}

func (p *` + name + `UpdatePayload) Validate() error {
	return validate.Struct(p)
}
`, nil
}

func RequestSkeletonTemplate(name string) (string, error) {
	return `package requests

// ` + name + `Payload is a skeleton request DTO generated without --all.
// Fill fields required by create operations.
type ` + name + `Payload struct{}

// ` + name + `UpdatePayload is a skeleton update request DTO generated without --all.
// Fill optional fields for partial updates.
type ` + name + `UpdatePayload struct{}
`, nil
}

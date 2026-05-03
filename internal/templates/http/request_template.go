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
type ` + name + `Payload struct {
	Email string ` + "`json:\"email\" example:\"john.doe@example.com\"`" + `
	Name  string ` + "`json:\"name\" example:\"John Doe\"`" + `
}

// ` + name + `UpdatePayload represents optional fields for partial updates.
type ` + name + `UpdatePayload struct {
	Email *string ` + "`json:\"email,omitempty\" example:\"john.doe@example.com\"`" + `
	Name  *string ` + "`json:\"name,omitempty\" example:\"John Doe\"`" + `
}
`, nil
}

func RequestSkeletonTemplate(name string) (string, error) {
	return `package requests

// ` + name + `Payload is a skeleton request DTO generated without --all.
// Fill fields required by create operations.
type ` + name + `Payload struct {}

// ` + name + `UpdatePayload is a skeleton update request DTO generated without --all.
// Fill optional fields for partial updates.
type ` + name + `UpdatePayload struct {}
`, nil
}

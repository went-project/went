package http

import (
	"went/internal/utils"
)

func ResourceTemplate(name string) (string, error) {
	config, err := utils.ReadConfigFile()
	if err != nil {
		return "", err
	}
	appName := config["name"].(string)
	return ResourceTemplateWithApp(name, appName)
}

func ResourceTemplateWithApp(name string, appName string) (string, error) {
	return `package resources

import (
	"` + appName + `/database/models"
)

// ` + name + `Resource is the single-item representation returned by the API.
type ` + name + `Resource struct {
	BaseResource
	Name  string ` + "`json:\"name\"`" + `
	Email string ` + "`json:\"email\"`" + `
}

// ` + name + `Collection is the paginated list representation returned by the API.
type ` + name + `Collection struct {
	Data []` + name + `Resource ` + "`json:\"data\"`" + `
	Meta PaginationMeta      ` + "`json:\"meta\"`" + `
}

// New` + name + `Resource transforms a models.` + name + ` into a ` + name + `Resource.
func New` + name + `Resource(m models.` + name + `) ` + name + `Resource {
	return ` + name + `Resource{
		BaseResource: NewBaseResource(m.ID, m.CreatedAt, m.UpdatedAt),
		Name:         m.Name,
		Email:        m.Email,
	}
}

// New` + name + `Collection builds a paginated ` + name + `Collection.
func New` + name + `Collection(items []models.` + name + `, total, page, perPage int64) ` + name + `Collection {
	data := make([]` + name + `Resource, len(items))
	for i, item := range items {
		data[i] = New` + name + `Resource(item)
	}
	return ` + name + `Collection{
		Data: data,
		Meta: BuildMeta(total, page, perPage, int64(len(items))),
	}
}
`, nil
}

func ResourceSkeletonTemplate(name string) (string, error) {
	return `package resources

// ` + name + `Resource is a skeleton resource generated without --all.
type ` + name + `Resource struct {
	BaseResource
}

// ` + name + `Collection is a skeleton collection generated without --all.
type ` + name + `Collection struct {
	Data []` + name + `Resource ` + "`json:\"data\"`" + `
	Meta PaginationMeta      ` + "`json:\"meta\"`" + `
}
`, nil
}

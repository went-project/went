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
	"time"
	"` + appName + `/database/models"
	"gorm.io/gorm"
)

// ` + name + `Resource is the single-item representation returned by the API.
type ` + name + `Resource struct {
	ID        uint      ` + "`json:\"id\"`" + `
	Name      string    ` + "`json:\"name\"`" + `
	Email     string    ` + "`json:\"email\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

// ` + name + `Collection is the paginated list representation returned by the API.
type ` + name + `Collection struct {
	Data []` + name + `Resource ` + "`json:\"data\"`" + `
	Meta PaginationMeta     ` + "`json:\"meta\"`" + `
}

// ` + name + `Query handles all DB read operations for ` + name + ` and transforms results.
type ` + name + `Query struct {
	db *gorm.DB
}

// New` + name + `Query creates a new ` + name + `Query bound to the given DB instance.
func New` + name + `Query(db *gorm.DB) *` + name + `Query {
	return &` + name + `Query{db: db}
}

// Paginate runs a COUNT + paginated SELECT and returns a ` + name + `Collection.
func (q *` + name + `Query) Paginate(page, perPage int64) (` + name + `Collection, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}

	var total int64
	if err := q.db.Model(&models.` + name + `{}).Count(&total).Error; err != nil {
		return ` + name + `Collection{}, err
	}

	var items []models.` + name + `
	if err := q.db.Offset(int((page-1)*perPage)).Limit(int(perPage)).Find(&items).Error; err != nil {
		return ` + name + `Collection{}, err
	}

	return New` + name + `Collection(items, total, page, perPage), nil
}

// Find fetches a single ` + name + ` by primary key and returns a ` + name + `Resource.
func (q *` + name + `Query) Find(id uint) (` + name + `Resource, error) {
	var m models.` + name + `
	if err := q.db.First(&m, id).Error; err != nil {
		return ` + name + `Resource{}, err
	}
	return New` + name + `Resource(m), nil
}

// New` + name + `Resource transforms a models.` + name + ` into a ` + name + `Resource.
func New` + name + `Resource(m models.` + name + `) ` + name + `Resource {
	return ` + name + `Resource{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
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

import "gorm.io/gorm"

// ` + name + `Resource is a skeleton resource generated without --all.
type ` + name + `Resource struct {}

// ` + name + `Collection is a skeleton collection generated without --all.
type ` + name + `Collection struct {
	Data []` + name + `Resource ` + "`json:\"data\"`" + `
	Meta PaginationMeta     ` + "`json:\"meta\"`" + `
}

// ` + name + `Query is a skeleton query builder generated without --all.
type ` + name + `Query struct {
	db *gorm.DB
}

func New` + name + `Query(db *gorm.DB) *` + name + `Query {
	return &` + name + `Query{db: db}
}
`, nil
}

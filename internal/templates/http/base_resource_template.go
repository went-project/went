package http

import "went/internal/utils"

func BaseResourceTemplate(appName string) string {
	return `package resources

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseResource holds the common fields shared by all API resource types.
type BaseResource struct {
	ID        uuid.UUID ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

// NewBaseResource constructs a BaseResource from the model's base fields.
func NewBaseResource(id uuid.UUID, createdAt, updatedAt time.Time) BaseResource {
	return BaseResource{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// PaginatedResponse is a generic paginated envelope for any resource type.
type PaginatedResponse[T any] struct {
	Data []T            ` + "`json:\"data\"`" + `
	Meta PaginationMeta ` + "`json:\"meta\"`" + `
}

// GenericQuery provides reusable DB read operations for any model type.
type GenericQuery[T any] struct {
	db *gorm.DB
}

// NewGenericQuery creates a new GenericQuery bound to the given DB instance.
func NewGenericQuery[T any](db *gorm.DB) *GenericQuery[T] {
	return &GenericQuery[T]{db: db}
}

// Paginate runs a COUNT + paginated SELECT and returns a PaginatedResponse.
func (q *GenericQuery[T]) Paginate(page, perPage int64) (PaginatedResponse[T], error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}

	var total int64
	var model T
	if err := q.db.Model(&model).Count(&total).Error; err != nil {
		return PaginatedResponse[T]{}, err
	}

	var items []T
	if err := q.db.Offset(int((page-1)*perPage)).Limit(int(perPage)).Find(&items).Error; err != nil {
		return PaginatedResponse[T]{}, err
	}

	return NewPaginatedResponse(items, total, page, perPage), nil
}

// FindByID fetches a single record by UUID primary key.
func (q *GenericQuery[T]) FindByID(id uuid.UUID) (T, error) {
	var m T
	err := q.db.First(&m, "id = ?", id).Error
	return m, err
}

// NewPaginatedResponse builds a PaginatedResponse from a slice and pagination params.
func NewPaginatedResponse[T any](data []T, total, page, perPage int64) PaginatedResponse[T] {
	return PaginatedResponse[T]{
		Data: data,
		Meta: BuildMeta(total, page, perPage, int64(len(data))),
	}
}
`
}

func BaseResourceTemplateFromConfig() (string, error) {
	config, err := utils.ReadConfigFile()
	if err != nil {
		return "", err
	}
	appName := config["name"].(string)
	return BaseResourceTemplate(appName), nil
}

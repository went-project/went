package http

func ControllerHelperTemplate() string {
	return `package helpers

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ParsePagination extracts page and per_page query params from the request.
func ParsePagination(r *http.Request) (page, perPage int64) {
	page, _ = strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	perPage, _ = strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}
	return
}

// Paginate runs a COUNT + paginated SELECT for any model type T.
func Paginate[T any](db *gorm.DB, page, perPage int64) (items []T, total int64, err error) {
	var model T
	if err = db.Model(&model).Count(&total).Error; err != nil {
		return
	}
	err = db.Offset(int((page-1)*perPage)).Limit(int(perPage)).Find(&items).Error
	return
}

// ParseID parses the {id} URL parameter as a UUID.
func ParseID(r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

// BuildCreateMap builds a column→value map from a payload's db struct tags.
func BuildCreateMap(payload interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	v := reflect.ValueOf(payload)
	t := reflect.TypeOf(payload)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		column := t.Field(i).Tag.Get("db")
		if column != "" {
			result[column] = v.Field(i).Interface()
		}
	}

	return result
}

// CreateFromPayload maps a payload into a model T via db tags and persists it.
func CreateFromPayload[T any](db *gorm.DB, payload interface{}) (T, error) {
	data := BuildCreateMap(payload)
	var model T
	mv := reflect.ValueOf(&model).Elem()
	mt := mv.Type()

	for i := 0; i < mt.NumField(); i++ {
		field := mt.Field(i)
		if field.Anonymous {
			fillEmbedded(mv.Field(i), data)
			continue
		}
		if v, ok := data[toSnakeCase(field.Name)]; ok {
			mv.Field(i).Set(reflect.ValueOf(v))
		}
	}

	err := db.Create(&model).Error
	return model, err
}

func fillEmbedded(v reflect.Value, data map[string]interface{}) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.Anonymous {
			fillEmbedded(v.Field(i), data)
			continue
		}
		if val, ok := data[toSnakeCase(field.Name)]; ok {
			v.Field(i).Set(reflect.ValueOf(val))
		}
	}
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(r + 32)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// BuildUpdateMap builds a column→value map from non-nil pointer fields using db tags.
func BuildUpdateMap(payload interface{}) map[string]interface{} {
	updates := map[string]interface{}{}
	v := reflect.ValueOf(payload)
	t := reflect.TypeOf(payload)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			column := t.Field(i).Tag.Get("db")
			if column != "" {
				updates[column] = field.Elem().Interface()
			}
		}
	}

	return updates
}
`
}

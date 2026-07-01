package http

import (
	"strings"
	"went/internal/utils"
)

func ControllerTemplate(name string) (string, error) {
	config, err := utils.ReadConfigFile()
	if err != nil {
		return "", err
	}
	appName := config["name"].(string)

	return ControllerTemplateWithApp(name, appName)
}

func ControllerTemplateWithApp(name string, appName string) (string, error) {
	resourceName := strings.ToLower(string(name[0]) + string(name[1:]))

	return `package controllers

import (
	"encoding/json"
	"net/http"
	"` + appName + `/database/models"
	"` + appName + `/http/requests"
	"` + appName + `/http/resources"
	"` + appName + `/internal/helpers"
	"` + appName + `/internal/responses"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ` + name + ` struct {
	DB *gorm.DB
}

func (c *` + name + `) findByID(id uuid.UUID, item *models.` + name + `) error {
	return c.DB.First(item, "id = ?", id).Error
}

// GetAll` + name + ` godoc
// @Summary List all ` + name + ` records
// @Description Get a paginated list of ` + name + ` records. Use ?page=1&per_page=15 to control pagination.
// @Tags ` + name + `
// @Produce json
// @Param page     query int false "Page number (default 1)"
// @Param per_page query int false "Items per page (default 15)"
// @Success 200 {object} resources.` + name + `Collection
// @Failure 500 {object} responses.ErrorBody
// @Router /` + resourceName + ` [get]
func (c *` + name + `) GetAll` + name + `(w http.ResponseWriter, r *http.Request) {
	page, perPage := helpers.ParsePagination(r)

	items, total, err := helpers.Paginate[models.` + name + `](c.DB, page, perPage)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, resources.New` + name + `Collection(items, total, page, perPage))
}

// Get` + name + `ByID godoc
// @Summary Get a ` + name + ` by ID
// @Description Get a single ` + name + ` record by ID
// @Tags ` + name + `
// @Produce json
// @Param id path string true "` + name + ` ID (UUID)"
// @Success 200 {object} resources.` + name + `Resource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /` + resourceName + `/{id} [get]
func (c *` + name + `) Get` + name + `ByID(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}

	var item models.` + name + `
	if err := c.DB.First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, resources.New` + name + `Resource(item))
}

// Create` + name + ` godoc
// @Summary Create a ` + name + `
// @Description Create a new ` + name + ` record
// @Tags ` + name + `
// @Accept json
// @Produce json
// @Param payload body requests.` + name + `Payload true "` + name + ` payload"
// @Success 201 {object} resources.` + name + `Resource
// @Failure 400 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /` + resourceName + ` [post]
func (c *` + name + `) Create` + name + `(w http.ResponseWriter, r *http.Request) {
	var payload requests.` + name + `Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := helpers.CreateFromPayload[models.` + name + `](c.DB, &payload)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resources.New` + name + `Resource(item))
}

// Update` + name + ` godoc
// @Summary Update a ` + name + `
// @Description Update an existing ` + name + ` record by ID
// @Tags ` + name + `
// @Accept json
// @Produce json
// @Param id path string true "` + name + ` ID (UUID)"
// @Param payload body requests.` + name + `UpdatePayload true "` + name + ` update payload"
// @Success 200 {object} resources.` + name + `Resource
// @Failure 400 {object} responses.ErrorBody
// @Failure 404 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /` + resourceName + `/{id} [put]
func (c *` + name + `) Update` + name + `(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}

	var item models.` + name + `
	err = c.findByID(id, &item)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.JSONError(w, r, http.StatusNotFound, "not found")
			return
		}
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	var payload requests.` + name + `UpdatePayload
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	updates := helpers.BuildUpdateMap(&payload)
	if len(updates) == 0 {
		responses.JSONError(w, r, http.StatusBadRequest, "no fields to update")
		return
	}

	err = c.DB.Model(&item).Updates(updates).Error
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	err = c.findByID(id, &item)
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, resources.New` + name + `Resource(item))
}

// Delete` + name + ` godoc
// @Summary Delete a ` + name + `
// @Description Delete an existing ` + name + ` record by ID
// @Tags ` + name + `
// @Param id path string true "` + name + ` ID (UUID)"
// @Success 204
// @Failure 400 {object} responses.ErrorBody
// @Failure 500 {object} responses.ErrorBody
// @Router /` + resourceName + `/{id} [delete]
func (c *` + name + `) Delete` + name + `(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ParseID(r)
	if err != nil {
		responses.JSONError(w, r, http.StatusBadRequest, "invalid id")
		return
	}

	err = c.DB.Delete(&models.` + name + `{}, "id = ?", id).Error
	if err != nil {
		responses.JSONError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusNoContent)
}
`, nil
}

func ControllerSkeletonTemplate(name string) (string, error) {
	config, err := utils.ReadConfigFile()
	if err != nil {
		return "", err
	}
	appName := config["name"].(string)
	return ControllerSkeletonTemplateWithApp(name, appName)
}

func ControllerSkeletonTemplateWithApp(name string, appName string) (string, error) {
	return `package controllers

import (
	"net/http"
	"` + appName + `/internal/responses"

	"gorm.io/gorm"
)

// ` + name + ` is a skeleton controller generated without --all.
// Fill handlers based on your resource and request DTOs.
type ` + name + ` struct {
	DB *gorm.DB
}

func (c *` + name + `) GetAll` + name + `(w http.ResponseWriter, r *http.Request) {
	responses.JSONError(w, r, http.StatusNotImplemented, "not implemented")
}

func (c *` + name + `) Get` + name + `ByID(w http.ResponseWriter, r *http.Request) {
	responses.JSONError(w, r, http.StatusNotImplemented, "not implemented")
}

func (c *` + name + `) Create` + name + `(w http.ResponseWriter, r *http.Request) {
	responses.JSONError(w, r, http.StatusNotImplemented, "not implemented")
}

func (c *` + name + `) Update` + name + `(w http.ResponseWriter, r *http.Request) {
	responses.JSONError(w, r, http.StatusNotImplemented, "not implemented")
}

func (c *` + name + `) Delete` + name + `(w http.ResponseWriter, r *http.Request) {
	responses.JSONError(w, r, http.StatusNotImplemented, "not implemented")
}
`, nil
}

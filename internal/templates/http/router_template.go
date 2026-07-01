package http

import (
	"strings"
	"went/internal/utils"
)

func RouterTemplate(name string) (string, error) {
	config, err := utils.ReadConfigFile()
	if err != nil {
		return "", err
	}
	app_name := config["name"].(string)
	return RouterTemplateWithApp(name, app_name)
}

func RouterTemplateWithApp(name string, appName string) (string, error) {
	controller_name := name
	resource_name := strings.ToLower(string(name[0]) + string(name[1:]))

	return `package routes

import (
	"` + appName + `/http/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func ` + controller_name + `Routes(r *chi.Mux, db *gorm.DB) {

	controller := &controllers.` + controller_name + `{DB: db}

	r.Get("/` + resource_name + `", controller.GetAll` + controller_name + `)
	r.Get("/` + resource_name + `/{id}", controller.Get` + controller_name + `ByID)
	r.Post("/` + resource_name + `", controller.Create` + controller_name + `)
	r.Put("/` + resource_name + `/{id}", controller.Update` + controller_name + `)
	r.Delete("/` + resource_name + `/{id}", controller.Delete` + controller_name + `)
}`, nil
}

func RouterSkeletonTemplate(name string) (string, error) {
	config, err := utils.ReadConfigFile()
	if err != nil {
		return "", err
	}
	appName := config["name"].(string)
	return RouterSkeletonTemplateWithApp(name, appName)
}

func RouterSkeletonTemplateWithApp(name string, appName string) (string, error) {
	controllerName := name

	return `package routes

import (
	"` + appName + `/http/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func ` + controllerName + `Routes(r *chi.Mux, db *gorm.DB) {
	_ = db
	_ = controllers.` + controllerName + `{}
	// TODO: define routes for ` + controllerName + ` controller
	// Example:
	// controller := &controllers.` + controllerName + `{DB: db}
	// r.Get("/...", controller.GetAll` + controllerName + `)
}
`, nil
}

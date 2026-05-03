package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	tmplApp "went/internal/templates/app"
	tmplConfig "went/internal/templates/config"
	tmplDB "went/internal/templates/database"
	tmplHTTP "went/internal/templates/http"
	tmplProject "went/internal/templates/project"
	"went/internal/utils"
)

type Project struct {
	Name        string   `json:"name"`
	ProjectType string   `json:"projectType"`
	RouterType  string   `json:"routerType"`
	Packages    []string `json:"packages"`
}

func runCommand(projectName string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = projectName
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %v failed: %w\n%s", name, args, err, string(output))
	}
	return nil
}

func CreateProject(req Project) error {

	projectName := utils.Slugify(req.Name)
	req.Name = projectName

	folders := []string{
		projectName,
		projectName + "/internal",
		projectName + "/internal/config",
		projectName + "/internal/providers",
		projectName + "/internal/responses",
		projectName + "/database/models",
		projectName + "/database/seeders",
		projectName + "/database/migrations",
		projectName + "/http",
		projectName + "/http/controllers",
		projectName + "/http/middlewares",
		projectName + "/routes",
		projectName + "/http/requests",
		projectName + "/http/resources",
		projectName + "/database",
	}

	for _, folder := range folders {
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// create a wentconfig.json file with the project configuration
	err := utils.CreateConfigFile(projectName, req)
	if err != nil {
		return err
	}

	templateContent, err := tmplApp.MainTemplate(projectName)
	if err != nil {
		return err
	}

	mainFilePath := filepath.Join(projectName, "main.go")
	err = utils.CreateFileWithContent(mainFilePath, templateContent)
	if err != nil {
		return err
	}

	templateContent, err = tmplConfig.ConfigTemplate()
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(projectName, "internal", "config", "config.go")
	err = utils.CreateFileWithContent(configFilePath, templateContent)
	if err != nil {
		return err
	}

	err = CreateGlobalErrorResponse(projectName)
	if err != nil {
		return err
	}

	templateContent, err = tmplDB.DatabaseProviderTemplate(projectName)
	if err != nil {
		return err
	}

	databaseProviderFilePath := filepath.Join(projectName, "internal", "providers", "database_provider.go")
	err = utils.CreateFileWithContent(databaseProviderFilePath, templateContent)
	if err != nil {
		return err
	}

	templateContent, err = tmplHTTP.AuthTemplate()
	if err != nil {
		return err
	}

	authFilePath := filepath.Join(projectName, "http", "middlewares", "auth.go")
	err = utils.CreateFileWithContent(authFilePath, templateContent)
	if err != nil {
		return err
	}

	templateContent, err = tmplHTTP.CorsTemplate()
	if err != nil {
		return err
	}

	corsFilePath := filepath.Join(projectName, "http", "middlewares", "cors.go")
	err = utils.CreateFileWithContent(corsFilePath, templateContent)
	if err != nil {
		return err
	}

	templateContent, err = tmplHTTP.MainRouterTemplate()
	if err != nil {
		return err
	}

	mainRouterFilePath := filepath.Join(projectName, "routes", "main_router.go")
	err = utils.CreateFileWithContent(mainRouterFilePath, templateContent)
	if err != nil {
		return err
	}

	templateContent, err = tmplHTTP.RouterTemplateWithApp("User", projectName)
	if err != nil {
		return err
	}

	userRouterFilePath := filepath.Join(projectName, "routes", "user_router.go")
	err = utils.CreateFileWithContent(userRouterFilePath, templateContent)
	if err != nil {
		return err
	}

	templateContent, err = tmplDB.ModelTemplate("User")
	if err != nil {
		return err
	}

	userModelFilePath := filepath.Join(projectName, "database", "models", "user.go")
	err = utils.CreateFileWithContent(userModelFilePath, templateContent)
	if err != nil {
		return err
	}

	templateContent, err = tmplHTTP.RequestTemplateWithApp("User", projectName)
	if err != nil {
		return err
	}

	err = utils.CreateFileWithContent(filepath.Join(projectName, "http", "requests", "user_request.go"), templateContent)
	if err != nil {
		return err
	}

	err = CreateControllerHelper(projectName)
	if err != nil {
		return err
	}

	templateContent, err = tmplHTTP.ControllerTemplateWithApp("User", projectName)
	if err != nil {
		return err
	}

	userControllerFilePath := filepath.Join(projectName, "http", "controllers", "user_controller.go")
	err = utils.CreateFileWithContent(userControllerFilePath, templateContent)
	if err != nil {
		return err
	}

	// resources
	err = utils.CreateFileWithContent(
		filepath.Join(projectName, "http", "resources", "pagination.go"),
		tmplHTTP.PaginationTemplate(),
	)
	if err != nil {
		return err
	}

	templateContent, err = tmplHTTP.ResourceTemplateWithApp("User", projectName)
	if err != nil {
		return err
	}

	err = utils.CreateFileWithContent(filepath.Join(projectName, "http", "resources", "user_resource.go"), templateContent)
	if err != nil {
		return err
	}

	upMigrationContent, downMigrationContent, err := tmplDB.MigrationTemplates("User")
	if err != nil {
		return err
	}

	upMigrationFilePath := filepath.Join(projectName, "database", "migrations", "000001_create_users_table.up.sql")
	err = utils.CreateFileWithContent(upMigrationFilePath, upMigrationContent)
	if err != nil {
		return err
	}

	downMigrationFilePath := filepath.Join(projectName, "database", "migrations", "000001_create_users_table.down.sql")
	err = utils.CreateFileWithContent(downMigrationFilePath, downMigrationContent)
	if err != nil {
		return err
	}

	err = utils.CreateFileWithContent(filepath.Join(projectName, ".env.example"), tmplProject.EnvExampleTemplate())
	if err != nil {
		return err
	}

	err = utils.CreateFileWithContent(filepath.Join(projectName, ".gitignore"), tmplProject.GitignoreTemplate())
	if err != nil {
		return err
	}

	err = utils.CreateFileWithContent(filepath.Join(projectName, ".devwatchignore"), tmplProject.DevwatchignoreTemplate())
	if err != nil {
		return err
	}

	err = utils.CreateFileWithContent(filepath.Join(projectName, "AGENTS.md"), tmplProject.AgentsTemplate(projectName))
	if err != nil {
		return err
	}

	err = utils.CreateFileWithContent(filepath.Join(projectName, "README.md"), tmplProject.ReadmeTemplate(projectName))
	if err != nil {
		return err
	}

	err = utils.CreateFileWithContent(filepath.Join(projectName, "CLAUDE.md"), tmplProject.ClaudeTemplate())
	if err != nil {
		return err
	}

	err = runCommand(projectName, "go", "mod", "init", projectName)
	if err != nil {
		return err
	}

	err = runCommand(projectName, "go", "mod", "tidy")
	if err != nil {
		return err
	}

	err = runCommand(projectName, "go", "get", "github.com/swaggo/swag/cmd/swag@latest")
	if err != nil {
		return err
	}

	err = runCommand(projectName, "swag", "init")
	if err != nil {
		return err
	}

	return nil
}

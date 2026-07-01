package handlers

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestConfig(t *testing.T, dir string) {
	t.Helper()
	configPath := filepath.Join(dir, "wentconfig.json")
	data := []byte(`{"name":"testapp"}`)
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("failed to write wentconfig.json: %v", err)
	}
}

func TestGenerateModelSingleArtifact(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	writeTestConfig(t, dir)

	if err := Generate(GenerationOptions{Name: "User", Type: GeneratorTypeModel}); err != nil {
		t.Fatal(err)
	}

	modelPath := filepath.Join(dir, "database", "models", "user.go")
	if _, err := os.Stat(modelPath); err != nil {
		t.Fatalf("expected model file %s to exist: %v", modelPath, err)
	}

	if _, err := os.Stat(filepath.Join(dir, "database", "migrations")); err == nil {
		t.Fatal("expected migrations directory to not exist for model-only generation")
	}
}

func TestGenerateModelSingleArtifactWithMigration(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	writeTestConfig(t, dir)

	if err := Generate(GenerationOptions{Name: "User", Type: GeneratorTypeModel, IncludeMigration: true}); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dir, "database", "models", "user.go")); err != nil {
		t.Fatalf("expected model file to exist: %v", err)
	}

	matches, err := filepath.Glob(filepath.Join(dir, "database", "migrations", "*.up.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected one up migration, got %d", len(matches))
	}

	matches, err = filepath.Glob(filepath.Join(dir, "database", "migrations", "*.down.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected one down migration, got %d", len(matches))
	}
}

func TestGenerateControllerSingleArtifact(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	writeTestConfig(t, dir)

	if err := Generate(GenerationOptions{Name: "User", Type: GeneratorTypeController}); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dir, "http", "controllers", "user_controller.go")); err != nil {
		t.Fatalf("expected controller file to exist: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "http", "requests")); err == nil {
		t.Fatal("expected no requests directory for controller-only generation")
	}

	if _, err := os.Stat(filepath.Join(dir, "http", "controllers", "helpers.go")); err == nil {
		t.Fatal("expected no controller helper for controller-only generation")
	}
}

func TestGenerateResourceSingleArtifact(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	writeTestConfig(t, dir)

	if err := Generate(GenerationOptions{Name: "User", Type: GeneratorTypeResource}); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dir, "http", "resources", "user_resource.go")); err != nil {
		t.Fatalf("expected resource file to exist: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "http", "resources", "pagination.go")); err == nil {
		t.Fatal("expected no pagination file for resource-only generation")
	}
}

func TestGenerateFullRelatedWithController(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	writeTestConfig(t, dir)

	if err := Generate(GenerationOptions{Name: "User", Type: GeneratorTypeController, All: true}); err != nil {
		t.Fatal(err)
	}

	paths := []string{
		filepath.Join(dir, "database", "models", "user.go"),
		filepath.Join(dir, "database", "migrations"),
		filepath.Join(dir, "http", "controllers", "helpers.go"),
		filepath.Join(dir, "http", "requests", "user_request.go"),
		filepath.Join(dir, "http", "controllers", "user_controller.go"),
		filepath.Join(dir, "routes", "user_router.go"),
		filepath.Join(dir, "http", "resources", "user_resource.go"),
		filepath.Join(dir, "http", "resources", "pagination.go"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("expected file or directory %s to exist: %v", p, err)
		}
	}
}

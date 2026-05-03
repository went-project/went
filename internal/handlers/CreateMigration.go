package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	tmplDB "went/internal/templates/database"
	"went/internal/utils"
)

func CreateMigration(name string, full bool) error {
	var (
		upContent   string
		downContent string
		err         error
	)

	if full {
		upContent, downContent, err = tmplDB.MigrationTemplates(name)
	} else {
		upContent, downContent, err = tmplDB.MigrationSkeletonTemplates(name)
	}
	if err != nil {
		return err
	}

	err = os.MkdirAll("database/migrations", os.ModePerm)
	if err != nil {
		return err
	}

	tableName := utils.Pluralize(name)

	prefix := time.Now().Format("20060102150405")
	baseName := fmt.Sprintf("%s_create_%s_table", prefix, tableName)

	upPath := filepath.Join("database", "migrations", baseName+".up.sql")
	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return err
	}

	downPath := filepath.Join("database", "migrations", baseName+".down.sql")
	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return err
	}

	return nil
}

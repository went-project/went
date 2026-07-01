package commands

import (
	"os"

	"went/internal/handlers"
	"went/internal/output"

	"github.com/spf13/cobra"
)

var MigrateFresh = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "Drop all tables and re-run every migration",
	Long:  "Runs all down migration files in reverse order, then re-runs all up migration files",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintInfo("Running fresh migration...")
		if err := handlers.FreshMigrations(); err != nil {
			output.PrintError("Fresh migration failed: %v", err)
			os.Exit(1)
		}
		output.PrintSuccess("Fresh migration completed successfully.")
	},
}

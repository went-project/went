package commands

import (
	"os"

	"went/internal/handlers"
	"went/internal/output"

	"github.com/spf13/cobra"
)

var Migrate = &cobra.Command{
	Use:   "migrate",
	Short: "Run pending migrations",
	Long:  "Runs all pending up migration files in database/migrations in order",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintInfo("Running migrations...")
		if err := handlers.RunMigrations(); err != nil {
			output.PrintError("Migration failed: %v", err)
			os.Exit(1)
		}
		output.PrintSuccess("Migrations completed successfully.")
	},
}

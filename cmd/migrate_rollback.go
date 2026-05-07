package commands

import (
	"os"

	"went/internal/handlers"
	"went/internal/output"

	"github.com/spf13/cobra"
)

var MigrateRollback = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last N migrations",
	Long:  "Runs the down migration files for the last N applied migrations (default: 1)",
	Run: func(cmd *cobra.Command, args []string) {
		step, _ := cmd.Flags().GetInt("step")

		output.PrintInfo("Rolling back %d migration(s)...", step)
		if err := handlers.RollbackMigrations(step); err != nil {
			output.PrintError("Rollback failed: %v", err)
			os.Exit(1)
		}
		output.PrintSuccess("Rollback completed successfully.")
	},
}

func init() {
	MigrateRollback.Flags().IntP("step", "s", 1, "Number of migrations to rollback")
}

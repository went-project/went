package commands

import (
	"fmt"
	"os"

	"went/internal/handlers"

	"github.com/spf13/cobra"
)

var MigrateRollback = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last N migrations",
	Long:  "Runs the down migration files for the last N applied migrations (default: 1)",
	Run: func(cmd *cobra.Command, args []string) {
		step, _ := cmd.Flags().GetInt("step")

		fmt.Printf("Rolling back %d migration(s)...\n", step)
		if err := handlers.RollbackMigrations(step); err != nil {
			fmt.Printf("Rollback failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Done.")
	},
}

func init() {
	MigrateRollback.Flags().IntP("step", "s", 1, "Number of migrations to rollback")
}

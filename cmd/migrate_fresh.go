package commands

import (
	"fmt"
	"os"

	"went/internal/handlers"

	"github.com/spf13/cobra"
)

var MigrateFresh = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "Drop all tables and re-run every migration",
	Long:  "Runs all down migration files in reverse order, then re-runs all up migration files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running fresh migration...")
		if err := handlers.FreshMigrations(); err != nil {
			fmt.Printf("Fresh migration failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Done.")
	},
}

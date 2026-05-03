package commands

import (
	"fmt"
	"os"

	"went/internal/handlers"

	"github.com/spf13/cobra"
)

var Migrate = &cobra.Command{
	Use:   "migrate",
	Short: "Run pending migrations",
	Long:  "Runs all pending up migration files in database/migrations in order",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running migrations...")
		if err := handlers.RunMigrations(); err != nil {
			fmt.Printf("Migration failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Done.")
	},
}

package commands

import (
	"os"

	"went/internal/handlers"
	"went/internal/output"

	"github.com/spf13/cobra"
)

var Run = &cobra.Command{
	Use:   "run",
	Short: "Run the current project with hot reload",
	Long:  "Run the current project with hot reload, watching local files and restarting when source changes.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := handlers.RunWithWatcher("."); err != nil {
			output.PrintError("Error: %v", err)
			os.Exit(1)
		}
	},
}

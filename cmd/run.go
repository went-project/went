package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"went/internal/handlers"
)

var Run = &cobra.Command{
	Use:   "run",
	Short: "Run the current project with hot reload",
	Long:  "Run the current project with hot reload, watching local files and restarting when source changes.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := handlers.RunWithWatcher("."); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

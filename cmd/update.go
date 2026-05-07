package commands

import (
	"went/internal/handlers"

	"github.com/spf13/cobra"
)

var Update = &cobra.Command{
	Use:   "update",
	Short: "Check for updates and install the latest went version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := handlers.UpdateWent(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
}

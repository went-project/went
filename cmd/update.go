package commands

import (
	"went/internal/handlers"
	"went/internal/utils"

	"github.com/spf13/cobra"
)

var (
	updateBetaFlag  bool
	updateForceFlag bool
)

var Update = &cobra.Command{
	Use:   "update",
	Short: "Check for updates and install the latest went version",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := handlers.UpdateOptions{
			Channel: utils.ParseChannelFlag(updateBetaFlag),
			Force:   updateForceFlag,
		}
		if err := handlers.UpdateWent(opts); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	Update.Flags().BoolVarP(&updateBetaFlag, "beta", "b", false, "Install latest beta release instead of stable")
	Update.Flags().BoolVarP(&updateForceFlag, "force", "f", false, "Force reinstall even if went is already up to date")
}

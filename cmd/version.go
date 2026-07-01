package commands

import (
	"os"

	"went/internal/handlers"
	"went/internal/output"
	"went/internal/utils"

	"github.com/spf13/cobra"
)

var versionBetaFlag bool

func printVersionInfo() {
	channel := utils.ParseChannelFlag(versionBetaFlag)
	output.PrintHeader("Version Info")
	output.PrintLine("Binary", utils.GetCurrentVersion())
	output.PrintLine("Channel", string(channel))

	if _, err := os.Stat("wentconfig.json"); err == nil {
		config, err := utils.ReadConfigFile()
		if err != nil {
			output.PrintWarning("Project config loaded, but version could not be read: %v", err)
		} else if versionValue, ok := config["version"]; ok {
			if versionStr, ok := versionValue.(string); ok && versionStr != "" {
				output.PrintLine("Project", versionStr)
			} else {
				output.PrintLine("Project", "unavailable")
			}
		} else {
			output.PrintLine("Project", "unavailable")
		}
	}

	result, err := handlers.CheckLatestVersion(handlers.UpdateOptions{Channel: channel})
	if err != nil {
		output.PrintWarning("Unable to check latest release: %v", err)
		return
	}

	output.PrintLine("Latest", result.Latest)
	if result.UpdateAvailable {
		if channel == utils.ChannelBeta {
			output.PrintSuccess("Beta update available — run 'went update --beta' to install the latest beta release.")
		} else {
			output.PrintSuccess("Update available — run 'went update' to install the latest version.")
		}
	} else {
		output.PrintInfo("No updates available.")
	}
}

var Version = &cobra.Command{
	Use:   "version",
	Short: "Display Went version information",
	Run: func(cmd *cobra.Command, args []string) {
		printVersionInfo()
	},
}

func init() {
	Version.Flags().BoolVarP(&versionBetaFlag, "beta", "b", false, "Check for the latest beta release instead of stable")
}

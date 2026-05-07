package commands

import (
	"fmt"
	"os"

	"went/internal/handlers"
	"went/internal/utils"

	"github.com/spf13/cobra"
)

var Version = &cobra.Command{
	Use:   "version",
	Short: "Sürümü göster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Went binary version: %s\n", utils.GetCurrentVersion())

		if _, err := os.Stat("wentconfig.json"); err == nil {
			config, err := utils.ReadConfigFile()
			if err != nil {
				fmt.Printf("Project config yüklendi, ancak version okunamadi: %v\n", err)
			} else if versionValue, ok := config["version"]; ok {
				if versionStr, ok := versionValue.(string); ok && versionStr != "" {
					fmt.Printf("Project version: %s\n", versionStr)
				} else {
					fmt.Println("Project version: unavailable")
				}
			} else {
				fmt.Println("Project version: unavailable")
			}
		}

		result, err := handlers.CheckLatestVersion()
		if err != nil {
			fmt.Printf("Update kontrolü yapılamadı: %v\n", err)
			return
		}

		fmt.Printf("Latest release: %s\n", result.Latest)
		if result.UpdateAvailable {
			fmt.Println("Update available: yes (run 'went update')")
		} else {
			fmt.Println("Update available: no")
		}
	},
}

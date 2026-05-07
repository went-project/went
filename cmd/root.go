package commands

import (
	"errors"
	"fmt"
	"os"

	"went/internal/utils"

	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use:     "went",
	Short:   "Went uygulaması",
	Long:    "Cobra ile yazılmış basit bir Went uygulaması",
	Version: utils.GetCurrentVersion(),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "create" || cmd.Name() == "help" || cmd.Name() == "update" || cmd.Name() == "version" {
			return nil
		}

		if cmd.Parent() == nil {
			return nil
		}

		if _, err := os.Stat("wentconfig.json"); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("wentconfig.json bulunamadı: bu komutu proje klasoru icinde calistirin")
			}
			return fmt.Errorf("wentconfig.json kontrol edilirken hata: %w", err)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Went uygulamasına hoş geldiniz!")
		cmd.Help()
	},
}

func init() {
	Root.AddCommand(Version)
	Root.AddCommand(Update)
	Root.AddCommand(Create)
	Root.AddCommand(GenerateRouter)
	Root.AddCommand(GenerateModel)
	Root.AddCommand(GenerateController)
	Root.AddCommand(GenerateMigration)
	Root.AddCommand(GenerateResource)
	Root.AddCommand(Run)
	Root.AddCommand(Migrate)
	Root.AddCommand(MigrateFresh)
	Root.AddCommand(MigrateRollback)
}

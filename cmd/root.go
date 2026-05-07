package commands

import (
	"errors"
	"fmt"
	"os"

	"went/internal/output"
	"went/internal/utils"

	"github.com/spf13/cobra"
)

var showVersionFlag bool

var Root = &cobra.Command{
	Use:     "went",
	Short:   "Went Framework CLI",
	Long:    "WENT.\nThe lightweight framework for Go developers.\nA clean framework for modern Go services.\nBuild fast. Stay simple. Minimal by design. Powerful by default.\n\nFor further documentation, visit https://wentframework.com/docs.html.",
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
				return fmt.Errorf("wentconfig.json not found: run this command in a project directory")
			}
			return fmt.Errorf("failed to validate wentconfig.json: %w", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersionFlag {
			printVersionInfo()
			return nil
		}

		output.PrintHeader("Went CLI")
		output.PrintInfo("Welcome to Went. Run `went --help` to explore available commands.")
		cmd.Help()
		return nil
	},
}

func init() {
	Root.PersistentFlags().BoolVarP(&showVersionFlag, "version", "v", false, "Show version information")
	Root.AddCommand(Version)
	Root.AddCommand(Update)
	Root.AddCommand(Create)
	Root.AddCommand(GenerateRouter)
	Root.AddCommand(GenerateModel)
	Root.AddCommand(GenerateController)
	Root.AddCommand(GenerateMigration)
	Root.AddCommand(GenerateResource)
	Root.AddCommand(GenerateRequest)
	Root.AddCommand(Run)
	Root.AddCommand(Migrate)
	Root.AddCommand(MigrateFresh)
	Root.AddCommand(MigrateRollback)
}

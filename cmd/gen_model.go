package commands

import (
	"went/internal/handlers"
	"went/internal/output"

	"github.com/spf13/cobra"
)

var GenerateModel = &cobra.Command{
	Use:   "gen:model [name]",
	Short: "Generate model",
	Long: `Generates a model file from template.

Flags:
  -m, --migration   Also generate up/down migration files for this model (skeleton if --all is not set)
  -a, --all         Generate all related templates in full form`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		withMigration, _ := cmd.Flags().GetBool("migration")
		withAll, _ := cmd.Flags().GetBool("all")

		if withAll {
			output.PrintInfo("Generating full related templates for %s", name)
			if err := handlers.CreateAllRelated(name); err != nil {
				output.PrintError("Failed to generate related templates: %v", err)
				return
			}
			output.PrintSuccess("All related templates for '%s' generated successfully!", name)
			return
		}

		if withMigration {
			output.PrintInfo("Note: --migration is implicit in skeleton related generation mode.")
		}

		output.PrintInfo("Generating skeleton related templates for %s", name)
		if err := handlers.CreateAllRelatedSkeleton(name); err != nil {
			output.PrintError("Failed to generate skeleton related templates: %v", err)
			return
		}
		output.PrintSuccess("Skeleton related templates for '%s' generated successfully!", name)
	},
}

func init() {
	GenerateModel.Flags().BoolP("migration", "m", false, "Also generate migration files for this model")
	GenerateModel.Flags().BoolP("all", "a", false, "Generate all related templates in full form")
}

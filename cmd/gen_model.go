package commands

import (
	"fmt"

	"went/internal/handlers"

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
			fmt.Printf("Generating full related templates: %s\n", name)
			if err := handlers.CreateAllRelated(name); err != nil {
				fmt.Printf("Error generating related templates: %v\n", err)
				return
			}
			fmt.Printf("All related templates for '%s' generated successfully!\n", name)
			return
		}

		if withMigration {
			fmt.Printf("Note: --migration is implicit in skeleton related generation mode.\n")
		}

		fmt.Printf("Generating skeleton related templates: %s\n", name)
		if err := handlers.CreateAllRelatedSkeleton(name); err != nil {
			fmt.Printf("Error generating skeleton related templates: %v\n", err)
			return
		}
		fmt.Printf("Skeleton related templates for '%s' generated successfully!\n", name)
	},
}

func init() {
	GenerateModel.Flags().BoolP("migration", "m", false, "Also generate migration files for this model")
	GenerateModel.Flags().BoolP("all", "a", false, "Generate all related templates in full form")
}

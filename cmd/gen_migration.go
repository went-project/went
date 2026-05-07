package commands

import (
	"went/internal/handlers"
	"went/internal/output"

	"github.com/spf13/cobra"
)

var GenerateMigration = &cobra.Command{
	Use:   "gen:migration [name]",
	Short: "Generate migration",
	Long:  "Generates skeleton up/down SQL migration files by default. Use --all to generate all related templates in full form.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
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

		output.PrintInfo("Generating skeleton related templates for %s", name)
		err := handlers.CreateAllRelatedSkeleton(name)
		if err != nil {
			output.PrintError("Failed to generate skeleton related templates: %v", err)
		} else {
			output.PrintSuccess("Skeleton related templates for '%s' generated successfully!", name)
		}
	},
}

func init() {
	GenerateMigration.Flags().BoolP("all", "a", false, "Generate all related templates in full form")
}

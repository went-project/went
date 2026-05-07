package commands

import (
	"went/internal/handlers"
	"went/internal/output"

	"github.com/spf13/cobra"
)

var GenerateResource = &cobra.Command{
	Use:   "gen:resource [name]",
	Short: "Generate resource",
	Long:  `Generates a skeleton resource file by default. Use --all to generate all related templates in full form. Also writes http/resources/pagination.go if it does not exist.`,
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
		if err := handlers.CreateAllRelatedSkeleton(name); err != nil {
			output.PrintError("Failed to generate skeleton related templates: %v", err)
			return
		}
		output.PrintSuccess("Skeleton related templates for '%s' generated successfully!", name)
	},
}

func init() {
	GenerateResource.Flags().BoolP("all", "a", false, "Generate all related templates in full form")
}

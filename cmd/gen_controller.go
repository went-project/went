package commands

import (
	"went/internal/handlers"
	"went/internal/output"

	"github.com/spf13/cobra"
)

var GenerateController = &cobra.Command{
	Use:   "gen:controller [name]",
	Short: "Generate controller",
	Long:  "Generates a skeleton controller by default. Use --all to generate all related templates in full form.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		withAll, _ := cmd.Flags().GetBool("all")

		opts := handlers.GenerationOptions{
			Name: name,
			Type: handlers.GeneratorTypeController,
			All:  withAll,
		}

		if withAll {
			output.PrintInfo("Generating full related templates for %s", name)
		} else {
			output.PrintInfo("Generating controller skeleton for %s", name)
		}

		if err := handlers.Generate(opts); err != nil {
			output.PrintError("Failed to generate templates: %v", err)
			return
		}

		output.PrintSuccess("Templates for '%s' generated successfully!", name)
	},
}

func init() {
	GenerateController.Flags().BoolP("all", "a", false, "Generate all related templates in full form")
}

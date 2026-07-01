package commands

import (
	"went/internal/handlers"
	"went/internal/output"

	"github.com/spf13/cobra"
)

var GenerateRouter = &cobra.Command{
	Use:   "gen:router [name]",
	Short: "Generate router",
	Long:  "Generates a skeleton router by default. Use --all to generate all related templates in full form.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		withAll, _ := cmd.Flags().GetBool("all")

		opts := handlers.GenerationOptions{
			Name: name,
			Type: handlers.GeneratorTypeRouter,
			All:  withAll,
		}

		if withAll {
			output.PrintInfo("Generating full related templates for %s", name)
		} else {
			output.PrintInfo("Generating router skeleton for %s", name)
		}

		if err := handlers.Generate(opts); err != nil {
			output.PrintError("Failed to generate templates: %v", err)
			return
		}

		output.PrintSuccess("Templates for '%s' generated successfully!", name)
	},
}

func init() {
	GenerateRouter.Flags().BoolP("all", "a", false, "Generate all related templates in full form")
}

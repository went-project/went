package commands

import (
	"fmt"

	"went/internal/handlers"

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
			fmt.Printf("Generating full related templates: %s\n", name)
			if err := handlers.CreateAllRelated(name); err != nil {
				fmt.Printf("Error generating related templates: %v\n", err)
				return
			}
			fmt.Printf("All related templates for '%s' generated successfully!\n", name)
			return
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
	GenerateResource.Flags().BoolP("all", "a", false, "Generate all related templates in full form")
}

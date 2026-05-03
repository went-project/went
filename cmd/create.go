package commands

import (
	"fmt"

	"went/internal/handlers"

	"github.com/spf13/cobra"
)

// Project variables
var ProjectType string // web,api,cli
var RouterType string  // gin, echo, fiber
var Packages []string  // database, auth, logger

var Create = &cobra.Command{
	Use:   "create [name]",
	Short: "Create project",
	Long:  "Creates a new project with the given name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		project := handlers.Project{
			Name:        name,
			ProjectType: ProjectType,
			RouterType:  RouterType,
			Packages:    Packages,
		}

		err := handlers.CreateProject(project)
		if err != nil {
			fmt.Printf("Error creating project: %v\n", err)
		} else {
			fmt.Printf("Project '%s' created successfully!\n", name)
		}
	},
}

func init() {
	Create.Flags().StringVarP(&ProjectType, "type", "t", "", "Type of the project (web, api, cli)")
	Create.Flags().StringVarP(&RouterType, "router", "r", "", "Type of the router (gin, echo, fiber)")
	Create.Flags().StringSliceVarP(&Packages, "packages", "p", []string{}, "Packages to include (database, auth, logger)")
}

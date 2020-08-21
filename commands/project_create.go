package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var projectCreateCmd = &cobra.Command{
	Use:   "create [path] [flags]",
	Short: `Create Gitlab project`,
	Long:  ``,
	RunE:  runCreateProject,
}

func runCreateProject(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		_ = cmd.Help()
		return nil
	}

	var (
		path      string
		visiblity gitlab.VisibilityValue
	)
	if len(args) == 1 {
		path = args[0]
	}

	name, _ := cmd.Flags().GetString("name")

	if path == "" && name == "" {
		fmt.Println("ERROR: Path or Name required to create project.")
		_ = cmd.Usage()
		return nil
	}

	description, _ := cmd.Flags().GetString("description")

	if internal, _ := cmd.Flags().GetBool("internal"); internal {
		visiblity = gitlab.InternalVisibility
	} else if private, _ := cmd.Flags().GetBool("private"); private {
		visiblity = gitlab.PrivateVisibility
	} else if public, _ := cmd.Flags().GetBool("public"); public {
		visiblity = gitlab.PublicVisibility
	}

	opts := &gitlab.CreateProjectOptions{
		Name:        gitlab.String(name),
		Path:        gitlab.String(path),
		Description: gitlab.String(description),
	}

	if visiblity != "" {
		opts.Visibility = &visiblity
	}

	project, err := createProject(opts)

	if err == nil {
		fmt.Println("Project created: ", project.WebURL)
	} else {
		fmt.Println("Error creating project: ", err)
	}
	return err
}

func init() {
	projectCreateCmd.Flags().StringP("name", "n", "", "name of the new project")
	projectCreateCmd.Flags().StringP("description", "d", "", "Description of the new project")
	projectCreateCmd.Flags().Bool("internal", false, "Make project internal: visible to any authenticated user (default)")
	projectCreateCmd.Flags().BoolP("private", "p", false, "Make project private: visible only to project members")
	projectCreateCmd.Flags().BoolP("public", "P", false, "Make project public: visible without any authentication")
	projectCmd.AddCommand(projectCreateCmd)
}

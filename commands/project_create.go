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

	defaultBranch, _ := cmd.Flags().GetString("defaultBranch")
	tags, _ := cmd.Flags().GetStringArray("tag")
	readme, _ := cmd.Flags().GetBool("readme")

	opts := &gitlab.CreateProjectOptions{
		Name:                 gitlab.String(name),
		Path:                 gitlab.String(path),
		Description:          gitlab.String(description),
		DefaultBranch:        gitlab.String(defaultBranch),
		TagList:              &tags,
		InitializeWithReadme: gitlab.Bool(readme),
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
	projectCreateCmd.Flags().StringP("name", "n", "", "Name of the new project")
	projectCreateCmd.Flags().StringP("description", "d", "", "Description of the new project")
	projectCreateCmd.Flags().String("defaultBranch", "", "Default branch of the project. If not provided, `master` by default.")
	projectCreateCmd.Flags().StringArrayP("tag", "t", []string{}, "The list of tags for the project.")
	projectCreateCmd.Flags().Bool("internal", false, "Make project internal: visible to any authenticated user (default)")
	projectCreateCmd.Flags().BoolP("private", "p", false, "Make project private: visible only to project members")
	projectCreateCmd.Flags().BoolP("public", "P", false, "Make project public: visible without any authentication")
	projectCreateCmd.Flags().Bool("readme", false, "Initialize project with README.md")
	projectCmd.AddCommand(projectCreateCmd)
}

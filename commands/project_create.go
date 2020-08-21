package commands

import (
	"glab/internal/git"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func createProject(opts *gitlab.CreateProjectOptions) (*gitlab.Project, error) {
	gitlabClient, _ := git.InitGitlabClient()
	project, _, err := gitlabClient.Projects.CreateProject(opts)
	if err != nil {
		return nil, err
	}
	return project, nil
}

var projectCreateCmd = &cobra.Command{
	Use:   "create [flags]",
	Short: `Create Gitlab project`,
	Long:  ``,
	Run:   runCreateProject,
}

func runCreateProject(cmd *cobra.Command, args []string) {
	if len(args) > 1 {
		_ = cmd.Help()
		return
	}

	// var path string
	// if len(args) > 0 {
	// 	path
	// }

	name, _ := cmd.Flags().GetString("name")
	opts := &gitlab.CreateProjectOptions{
		Name: gitlab.String(name),
	}
	createProject(opts)
}

func init() {
	projectCreateCmd.Flags().StringP("description", "d", "", "Description of the new project")
	projectCreateCmd.Flags().Bool("internal", false, "Make project internal: visible to any authenticated user (default)")
	projectCreateCmd.Flags().StringP("name", "n", "", "name of the new project")
	projectCreateCmd.Flags().BoolP("private", "p", false, "Make project private: visible only to project members")
	projectCreateCmd.Flags().BoolP("public", "P", false, "Make project public: visible without any authentication")
	projectCreateCmd.Flags().BoolP("create-source-branch", "", false, "Create source branch if it does not exist")
	projectCmd.AddCommand(projectCreateCmd)
}

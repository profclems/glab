package commands

import (
	"glab/internal/git"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func getProject(projectID interface{}) (*gitlab.Project, error) {
	gitlabClient, _ := git.InitGitlabClient()
	opts := &gitlab.GetProjectOptions{
		Statistics:           gitlab.Bool(true),
		License:              gitlab.Bool(true),
		WithCustomAttributes: gitlab.Bool(true),
	}
	project, _, err := gitlabClient.Projects.GetProject(projectID, opts)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func createProject(opts *gitlab.CreateProjectOptions) (*gitlab.Project, error) {
	gitlabClient, _ := git.InitGitlabClient()
	project, _, err := gitlabClient.Projects.CreateProject(opts)
	if err != nil {
		return nil, err
	}
	return project, nil
}

var projectCmd = &cobra.Command{
	Use:   "project <command> [flags]",
	Short: `Work with GitLab projects`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 2 {
			_ = cmd.Help()
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(projectCmd)
}

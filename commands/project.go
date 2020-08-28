package commands

import (
	"glab/internal/git"

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

func listProjectTree(projectID interface{}, opts *gitlab.ListTreeOptions) ([]*gitlab.TreeNode, error)  {
	gitlabClient, _ := git.InitGitlabClient()
	projectTree, _, err := gitlabClient.Repositories.ListTree(projectID, opts)
	if err != nil {
		return nil, err
	}
	return projectTree, nil
}

// projectCmd is the same as the repoCmd since repo has project as an alias
var projectCmd = repoCmd

package commands

import (
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
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

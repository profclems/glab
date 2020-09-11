package commands

import (
	"github.com/profclems/glab/internal/git"

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
func deleteProject(projectID interface{}) (*gitlab.Response, error) {
	gitlabClient, _ := git.InitGitlabClient()
	project, err := gitlabClient.Projects.DeleteProject(projectID)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func createProject(opts *gitlab.CreateProjectOptions) (*gitlab.Project, error) {
	gitlabClient, _ := git.InitGitlabClient(false)
	project, _, err := gitlabClient.Projects.CreateProject(opts)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func getGroup(groupID interface{}) (*gitlab.Group, error)  {
	gitlabClient, _ := git.InitGitlabClient(false)
	group, _, err := gitlabClient.Groups.GetGroup(groupID)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// projectCmd is the same as the repoCmd since repo has project as an alias
var projectCmd = repoCmd

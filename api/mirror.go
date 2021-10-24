package api

import (
	"github.com/xanzy/go-gitlab"
)

var CreatePushMirror = func(
	client *gitlab.Client,
	projectID interface{},
	url string,
	enabled bool,
	onlyProtectedBranches bool,
	keepDivergentRefs bool,
) (*gitlab.ProjectMirror, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	opt := &gitlab.AddProjectMirrorOptions{
		URL:                   &url,
		Enabled:               &enabled,
		OnlyProtectedBranches: &onlyProtectedBranches,
		KeepDivergentRefs:     &keepDivergentRefs,
	}
	pm, _, err := client.ProjectMirrors.AddProjectMirror(projectID, opt)
	return pm, err
}

var CreatePullMirror = func(
	client *gitlab.Client,
	projectID interface{},
	url string,
	enabled bool,
	onlyProtectedBranches bool,
) error {
	if client == nil {
		client = apiClient.Lab()
	}
	opt := &gitlab.EditProjectOptions{
		ImportURL:                   &url,
		Mirror:                      &enabled,
		OnlyMirrorProtectedBranches: &onlyProtectedBranches,
	}
	_, _, err := client.Projects.EditProject(projectID, opt)
	return err
}

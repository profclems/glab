package api

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type AddPullMirrorOptions struct {
	URL                   *string `url:"import_url,omitempty" json:"import_url,omitempty"`
	Enabled               *bool   `url:"mirror,omitempty" json:"mirror,omitempty"`
	OnlyProtectedBranches *bool   `url:"only_mirror_protected_branches,omitempty" json:"only_mirror_protected_branches,omitempty"`
}

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
	projectID int,
	url string,
	enabled bool,
	onlyProtectedBranches bool,
) error {
	opt := &AddPullMirrorOptions{
		URL:                   &url,
		Enabled:               &enabled,
		OnlyProtectedBranches: &onlyProtectedBranches,
	}
	u := fmt.Sprintf("projects/%d/remote_mirrors/", projectID)
	// TODO
	fmt.Println(u, opt)
	return nil
}

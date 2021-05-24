package api

import "github.com/xanzy/go-gitlab"

var CreateBranch = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateBranchOptions) (*gitlab.Branch, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	branch, _, err := client.Branches.CreateBranch(projectID, opts)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

var ListBranches = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListBranchesOptions) ([]*gitlab.Branch, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	branches, _, err := client.Branches.ListBranches(projectID, opts)
	if err != nil {
		return nil, err
	}

	return branches, nil
}

package api

import "github.com/xanzy/go-gitlab"

var ListMilestones = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListMilestonesOptions) ([]*gitlab.Milestone, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	if opts.PerPage == 0 {
		opts.PerPage = DefaultListLimit
	}

	milestone, _, err := client.Milestones.ListMilestones(projectID, opts)
	if err != nil {
		return nil, err
	}
	return milestone, nil
}

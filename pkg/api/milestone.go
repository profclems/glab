package api

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

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

var MilestoneByTitle = func(client *gitlab.Client, projectID interface{}, name string) (*gitlab.Milestone, error) {
	opts := &gitlab.ListMilestonesOptions{Title: gitlab.String(name)}

	if client == nil {
		client = apiClient.Lab()
	}

	if opts.PerPage == 0 {
		opts.PerPage = DefaultListLimit
	}

	milestones, _, err := client.Milestones.ListMilestones(projectID, opts)
	if err != nil {
		return nil, err
	}

	if len(milestones) != 1 {
		return nil, fmt.Errorf("failed to find milestone by title: %s", name)
	}

	return milestones[0], nil
}

package api

import "github.com/xanzy/go-gitlab"

var CreateLabel = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateLabelOptions) (*gitlab.Label, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	label, _, err := client.Labels.CreateLabel(projectID, opts)
	if err != nil {
		return nil, err
	}
	return label, nil
}

var ListLabels = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListLabelsOptions) ([]*gitlab.Label, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	if opts.PerPage == 0 {
		opts.PerPage = DefaultListLimit
	}

	label, _, err := client.Labels.ListLabels(projectID, opts)
	if err != nil {
		return nil, err
	}
	return label, nil
}

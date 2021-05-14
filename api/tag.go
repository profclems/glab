package api

import "github.com/xanzy/go-gitlab"

var ListTags = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListTagsOptions) ([]*gitlab.Tag, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	tags, _, err := client.Tags.ListTags(projectID, opts)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

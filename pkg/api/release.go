package api

import "github.com/xanzy/go-gitlab"

var GetRelease = func(client *gitlab.Client, projectID interface{}, tag string) (*gitlab.Release, error) {
	if client == nil {
		client = apiClient
	}

	release, _, err := client.Releases.GetRelease(projectID, tag)
	if err != nil {
		return nil, err
	}

	return release, nil
}
var ListReleases = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListReleasesOptions) ([]*gitlab.Release, error) {
	if client == nil {
		client = apiClient
	}

	releases, _, err := apiClient.Releases.ListReleases(projectID, opts)
	if err != nil {
		return nil, err
	}

	return releases, nil
}

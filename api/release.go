package api

import "github.com/xanzy/go-gitlab"

var CreateRelease = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateReleaseOptions) (*gitlab.Release, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	release, _, err := client.Releases.CreateRelease(projectID, opts)
	if err != nil {
		return nil, err
	}

	return release, nil
}

var GetRelease = func(client *gitlab.Client, projectID interface{}, tag string) (*gitlab.Release, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	release, _, err := client.Releases.GetRelease(projectID, tag)
	if err != nil {
		return nil, err
	}

	return release, nil
}

var ListReleases = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListReleasesOptions) ([]*gitlab.Release, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	if opts.PerPage == 0 {
		opts.PerPage = DefaultListLimit
	}

	releases, _, err := client.Releases.ListReleases(projectID, opts)
	if err != nil {
		return nil, err
	}

	return releases, nil
}

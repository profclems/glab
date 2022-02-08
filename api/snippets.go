package api

import "github.com/xanzy/go-gitlab"

// CreateSnippet for the user inside the users snippets
var CreateSnippet = func(
	client *gitlab.Client,
	projectID interface{},
	opts *gitlab.CreateSnippetOptions,
) (*gitlab.Snippet, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	snipet, _, err := client.Snippets.CreateSnippet(opts)
	if err != nil {
		return nil, err
	}
	return snipet, err
}

// CreateProjectSnippet inside the project
var CreateProjectSnippet = func(
	client *gitlab.Client,
	projectID interface{},
	opts *gitlab.CreateProjectSnippetOptions,
) (*gitlab.Snippet, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	snipet, _, err := client.ProjectSnippets.CreateSnippet(projectID, opts)
	if err != nil {
		return nil, err
	}
	return snipet, err
}

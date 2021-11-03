package api

import (
	"github.com/hashicorp/go-retryablehttp"

	"github.com/xanzy/go-gitlab"
)

var CreateProjectVariable = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateProjectVariableOptions) (*gitlab.ProjectVariable, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	vars, _, err := client.ProjectVariables.CreateVariable(projectID, opts)
	if err != nil {
		return nil, err
	}

	return vars, nil
}

var ListProjectVariables = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListProjectVariablesOptions) ([]*gitlab.ProjectVariable, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	vars, _, err := client.ProjectVariables.ListVariables(projectID, opts)
	if err != nil {
		return nil, err
	}

	return vars, nil
}

var DeleteProjectVariable = func(client *gitlab.Client, projectID interface{}, key string, scope string) error {
	if client == nil {
		client = apiClient.Lab()
	}

	var filter = func(request *retryablehttp.Request) error {
		q := request.URL.Query()
		q.Add("filter[environment_scope]", scope)

		request.URL.RawQuery = q.Encode()

		return nil
	}

	_, err := client.ProjectVariables.RemoveVariable(projectID, key, filter)

	if err != nil {
		return err
	}

	return nil
}

var ListGroupVariables = func(client *gitlab.Client, groupID interface{}, opts *gitlab.ListGroupVariablesOptions) ([]*gitlab.GroupVariable, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	vars, _, err := client.GroupVariables.ListVariables(groupID, opts)
	if err != nil {
		return nil, err
	}

	return vars, nil
}

var CreateGroupVariable = func(client *gitlab.Client, groupID interface{}, opts *gitlab.CreateGroupVariableOptions) (*gitlab.GroupVariable, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	vars, _, err := client.GroupVariables.CreateVariable(groupID, opts)
	if err != nil {
		return nil, err
	}

	return vars, nil
}

var DeleteGroupVariable = func(client *gitlab.Client, groupID interface{}, key string) error {
	if client == nil {
		client = apiClient.Lab()
	}

	_, err := client.GroupVariables.RemoveVariable(groupID, key)

	if err != nil {
		return err
	}

	return nil
}

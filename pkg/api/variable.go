package api

import "github.com/xanzy/go-gitlab"

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

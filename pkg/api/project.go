package api

import "github.com/xanzy/go-gitlab"

var GetProject = func(client *gitlab.Client, projectID interface{}) (*gitlab.Project, error) {
	if client == nil {
		client = apiClient
	}
	opts := &gitlab.GetProjectOptions{
		Statistics:           gitlab.Bool(true),
		License:              gitlab.Bool(true),
		WithCustomAttributes: gitlab.Bool(true),
	}
	project, _, err := client.Projects.GetProject(projectID, opts)
	if err != nil {
		return nil, err
	}
	return project, nil
}
var GetRepository = func(client *gitlab.Client, projectID interface{}) (*gitlab.Project, error) {
	if client == nil {
		client = apiClient
	}
	opts := &gitlab.GetProjectOptions{
		Statistics:           gitlab.Bool(true),
		License:              gitlab.Bool(true),
		WithCustomAttributes: gitlab.Bool(true),
	}
	project, _, err := client.Projects.GetProject(projectID, opts)
	if err != nil {
		return nil, err
	}
	return project, nil
}

var DeleteProject = func(client *gitlab.Client, projectID interface{}) (*gitlab.Response, error) {
	if client == nil {
		client = apiClient
	}
	project, err := client.Projects.DeleteProject(projectID)
	if err != nil {
		return nil, err
	}
	return project, nil
}

var CreateProject = func(client *gitlab.Client, opts *gitlab.CreateProjectOptions) (*gitlab.Project, error) {
	if client == nil {
		client = apiClient
	}
	project, _, err := client.Projects.CreateProject(opts)
	if err != nil {
		return nil, err
	}
	return project, nil
}

var GetGroup = func(client *gitlab.Client, groupID interface{}) (*gitlab.Group, error) {
	if client == nil {
		client = apiClient
	}
	group, _, err := client.Groups.GetGroup(groupID)
	if err != nil {
		return nil, err
	}
	return group, nil
}

var ListGroupProjects = func(client *gitlab.Client, groupID interface{}, opts *gitlab.ListGroupProjectsOptions) ([]*gitlab.Project, error) {
	if client == nil {
		client = apiClient
	}
	project, _, err := client.Groups.ListGroupProjects(groupID, opts)
	if err != nil {
		return nil, err
	}
	return project, nil
}
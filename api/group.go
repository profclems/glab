package api

import "github.com/xanzy/go-gitlab"

var GetGroup = func(client *gitlab.Client, groupID interface{}) (*gitlab.Group, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	group, _, err := client.Groups.GetGroup(groupID)
	if err != nil {
		return nil, err
	}
	return group, nil
}

var ListGroups = func(client *gitlab.Client, opts *gitlab.ListGroupsOptions) ([]*gitlab.Group, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	groups, _, err := client.Groups.ListGroups(opts)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

var ListSubgroups = func(client *gitlab.Client, groupId interface{}, opts *gitlab.ListSubgroupsOptions) ([]*gitlab.Group, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	groups, _, err := client.Groups.ListSubgroups(groupId, opts)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

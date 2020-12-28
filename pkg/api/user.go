package api

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

var CurrentUser = func(client *gitlab.Client) (*gitlab.User, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	u, _, err := client.Users.CurrentUser()
	if err != nil {
		return nil, err
	}
	return u, nil
}

var UserByName = func(client *gitlab.Client, name string) (*gitlab.User, error) {
	opts := &gitlab.ListUsersOptions{Username: gitlab.String(name)}

	if client == nil {
		client = apiClient.Lab()
	}

	if opts.PerPage == 0 {
		opts.PerPage = DefaultListLimit
	}

	// Handle special case of '@me' which maps to the currently authenticated user
	if name == "@me" {
		return CurrentUser(client)
	}

	users, _, err := client.Users.ListUsers(opts)
	if err != nil {
		return nil, err
	}

	if len(users) != 1 {
		return nil, fmt.Errorf("failed to find user by name : %s", name)
	}

	return users[0], nil
}

var UsersByNames = func(client *gitlab.Client, names []string) ([]*gitlab.User, error) {
	users := make([]*gitlab.User, 0)
	for _, name := range names {
		user, err := UserByName(client, name)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

package api

import (
	"errors"
	"fmt"

	"github.com/xanzy/go-gitlab"
)

var CurrentUser = func(client *gitlab.Client) (*gitlab.User, error) {
	if client == nil {
		client = apiClient
	}
	u, _, err := client.Users.CurrentUser()
	if err != nil {
		return nil, err
	}
	return u, nil
}

var UserByName = func(client *gitlab.Client, name string) (*gitlab.User, error) {
	opts := &gitlab.ListUsersOptions{Username: gitlab.String(name)}
	users, _, err := apiClient.Users.ListUsers(opts)
	if err != nil {
		return nil, err
	}

	if len(users) != 1 {
		return nil, errors.New(fmt.Sprintf("failed to find user by name : %s", name))
	}

	return users[0], nil
}

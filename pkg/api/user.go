package api

import (
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

var GetUser = func(client *gitlab.Client, uid int) (*gitlab.User, error) {
	if client == nil {
		client = apiClient
	}
	u, _, err := client.Users.GetUser(uid)
	if err != nil {
		return nil, err
	}
	return u, nil
}

var GetUsername = func(client *gitlab.Client, uid int) (string, error) {
	if client == nil {
		client = apiClient
	}
	u, err := GetUser(client, uid)
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

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

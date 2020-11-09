package api

import (
	"github.com/xanzy/go-gitlab"
)

var CurrentUserEvents = func(client *gitlab.Client) ([]*gitlab.ContributionEvent, error) {
	if client == nil {
		client = apiClient
	}

	events, _, err := client.Events.ListCurrentUserContributionEvents(&gitlab.ListContributionEventsOptions{})
	if err != nil {
		return nil, err
	}
	return events, nil
}

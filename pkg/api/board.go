package api

import "github.com/xanzy/go-gitlab"

var CreateIssueBoard = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateIssueBoardOptions) (*gitlab.IssueBoard, error) {
	if client == nil {
		client = apiClient
	}
	board, _, err := client.Boards.CreateIssueBoard(projectID, opts)
	if err != nil {
		return nil, err
	}

	return board, nil
}

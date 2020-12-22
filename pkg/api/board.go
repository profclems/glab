package api

import "github.com/xanzy/go-gitlab"

var CreateIssueBoard = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateIssueBoardOptions) (*gitlab.IssueBoard, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	board, _, err := client.Boards.CreateIssueBoard(projectID, opts)
	if err != nil {
		return nil, err
	}

	return board, nil
}

var ListIssueBoards = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListIssueBoardsOptions) ([]*gitlab.IssueBoard, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	boards, _, err := client.Boards.ListIssueBoards(projectID, opts)
	if err != nil {
		return nil, err
	}

	return boards, nil
}

var GetIssueBoardLists = func(client *gitlab.Client, projectID interface{}, boardID int, opts *gitlab.GetIssueBoardListsOptions) ([]*gitlab.BoardList, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	boardLists, _, err := client.Boards.GetIssueBoardLists(projectID, boardID, opts)
	if err != nil {
		return nil, err
	}

	return boardLists, nil
}

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

var ListGroupIssueBoards = func(client *gitlab.Client, groupID interface{}, opts *gitlab.ListGroupIssueBoardsOptions) ([]*gitlab.GroupIssueBoard, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	boards, _, err := client.GroupIssueBoards.ListGroupIssueBoards(groupID, opts)
	if err != nil {
		return nil, err
	}

	return boards, nil
}

var ListProjectIssueBoards = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListIssueBoardsOptions) ([]*gitlab.IssueBoard, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	boards, _, err := client.Boards.ListIssueBoards(projectID, opts)
	if err != nil {
		return nil, err
	}

	return boards, nil
}

var GetPojectIssueBoardLists = func(client *gitlab.Client, projectID interface{}, boardID int, opts *gitlab.GetIssueBoardListsOptions) ([]*gitlab.BoardList, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	boardLists, _, err := client.Boards.GetIssueBoardLists(projectID, boardID, opts)
	if err != nil {
		return nil, err
	}

	return boardLists, nil
}

var GetGroupIssueBoardLists = func(client *gitlab.Client, groupID interface{}, boardID int, opts *gitlab.ListGroupIssueBoardListsOptions) ([]*gitlab.BoardList, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	boardLists, _, err := client.GroupIssueBoards.ListGroupIssueBoardLists(groupID, boardID, opts)
	if err != nil {
		return nil, err
	}

	return boardLists, nil
}

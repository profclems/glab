// This is a silly wrapper for go-gitlab but helps maintain consistency
package api

import (
	"github.com/xanzy/go-gitlab"
)

var ListIssueNotes = func(client *gitlab.Client, projectID interface{}, issueID int, opts *gitlab.ListIssueNotesOptions) ([]*gitlab.Note, error)  {
	if client == nil {
		client = apiClient
	}
	notes, _, err := client.Notes.ListIssueNotes(projectID, issueID, opts)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

var UpdateIssue = func(client *gitlab.Client, projectID interface{}, issueID int, opts *gitlab.UpdateIssueOptions) (*gitlab.Issue, error)  {
	if client == nil {
		client = apiClient
	}
	issue, _, err := apiClient.Issues.UpdateIssue(projectID, issueID, opts)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

var GetIssue = func(client *gitlab.Client, projectID interface{}, issueID int) (*gitlab.Issue, error)  {
	if client == nil {
		client = apiClient
	}
	issue, _, err := apiClient.Issues.GetIssue(projectID, issueID)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

var ListIssues = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListProjectIssuesOptions) ([]*gitlab.Issue, error)  {
	issues, _, err := apiClient.Issues.ListProjectIssues(projectID, opts)
	if err != nil {
		return nil, err
	}

	return issues, nil
}

var CreateIssue = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateIssueOptions) (*gitlab.Issue, error)  {
	if client == nil {
		client = apiClient
	}
	issue, _, err := apiClient.Issues.CreateIssue(projectID, opts)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

var DeleteIssue = func(client *gitlab.Client, projectID interface{}, issueID int) error {
	if client == nil {
		client = apiClient
	}

	_, err := apiClient.Issues.DeleteIssue(projectID, issueID)
	if err != nil {
		return err
	}

	return nil
}

var CreateIssueNote = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.CreateIssueNoteOptions) (*gitlab.Note, error) {
	if client == nil {
		client = apiClient
	}

	note, _, err := client.Notes.CreateIssueNote(projectID, mrID, opts)
	if err != nil {
		return note, err
	}

	return note, nil
}

var SubscribeToIssue = func(client *gitlab.Client, projectID interface{}, issueID int, opts gitlab.RequestOptionFunc) (*gitlab.Issue, error) {
	if client == nil {
		client = apiClient
	}

	issue, _, err := client.Issues.SubscribeToIssue(projectID, issueID, opts)
	if err != nil {
		return issue, err
	}

	return issue, nil
}

var UnsubscribeFromIssue = func(client *gitlab.Client, projectID interface{}, issueID int, opts gitlab.RequestOptionFunc) (*gitlab.Issue, error) {
	if client == nil {
		client = apiClient
	}

	issue, _, err := client.Issues.UnsubscribeFromIssue(projectID, issueID, opts)
	if err != nil {
		return issue, err
	}

	return issue, nil
}

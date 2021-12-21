// This is a silly wrapper for go-gitlab but helps maintain consistency
package api

import (
	"errors"

	"github.com/xanzy/go-gitlab"
)

var ListIssueNotes = func(client *gitlab.Client, projectID interface{}, issueID int, opts *gitlab.ListIssueNotesOptions) ([]*gitlab.Note, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultListLimit
	}
	notes, _, err := client.Notes.ListIssueNotes(projectID, issueID, opts)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

var UpdateIssue = func(client *gitlab.Client, projectID interface{}, issueID int, opts *gitlab.UpdateIssueOptions) (*gitlab.Issue, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	issue, _, err := client.Issues.UpdateIssue(projectID, issueID, opts)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

var GetIssue = func(client *gitlab.Client, projectID interface{}, issueID int) (*gitlab.Issue, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	issue, _, err := client.Issues.GetIssue(projectID, issueID)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

var ProjectListIssueOptionsToGroup = func(l *gitlab.ListProjectIssuesOptions) *gitlab.ListGroupIssuesOptions {
	return &gitlab.ListGroupIssuesOptions{
		ListOptions:        l.ListOptions,
		State:              l.State,
		Labels:             l.Labels,
		NotLabels:          l.NotLabels,
		WithLabelDetails:   l.WithLabelDetails,
		IIDs:               l.IIDs,
		Milestone:          l.Milestone,
		Scope:              l.Scope,
		AuthorID:           l.AuthorID,
		NotAuthorID:        l.NotAuthorID,
		AssigneeID:         l.AssigneeID,
		NotAssigneeID:      l.NotAssigneeID,
		AssigneeUsername:   l.AssigneeUsername,
		MyReactionEmoji:    l.MyReactionEmoji,
		NotMyReactionEmoji: l.NotMyReactionEmoji,
		OrderBy:            l.OrderBy,
		Sort:               l.Sort,
		Search:             l.Search,
		In:                 l.In,
		CreatedAfter:       l.CreatedAfter,
		CreatedBefore:      l.CreatedBefore,
		UpdatedAfter:       l.UpdatedAfter,
		UpdatedBefore:      l.UpdatedBefore,
		IssueType:          l.IssueType,
	}
}

var ListGroupIssues = func(client *gitlab.Client, groupID interface{}, opts *gitlab.ListGroupIssuesOptions) ([]*gitlab.Issue, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultListLimit
	}
	issues, _, err := client.Issues.ListGroupIssues(groupID, opts)
	if err != nil {
		return nil, err
	}

	return issues, nil
}

var ListIssues = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListProjectIssuesOptions) ([]*gitlab.Issue, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	if opts.PerPage == 0 {
		opts.PerPage = DefaultListLimit
	}
	issues, _, err := client.Issues.ListProjectIssues(projectID, opts)
	if err != nil {
		return nil, err
	}

	return issues, nil
}

var CreateIssue = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateIssueOptions) (*gitlab.Issue, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	issue, _, err := client.Issues.CreateIssue(projectID, opts)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

var DeleteIssue = func(client *gitlab.Client, projectID interface{}, issueID int) error {
	if client == nil {
		client = apiClient.Lab()
	}

	_, err := client.Issues.DeleteIssue(projectID, issueID)
	if err != nil {
		return err
	}

	return nil
}

var CreateIssueNote = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.CreateIssueNoteOptions) (*gitlab.Note, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	note, _, err := client.Notes.CreateIssueNote(projectID, mrID, opts)
	if err != nil {
		return note, err
	}

	return note, nil
}

var SubscribeToIssue = func(client *gitlab.Client, projectID interface{}, issueID int, opts gitlab.RequestOptionFunc) (*gitlab.Issue, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	issue, resp, err := client.Issues.SubscribeToIssue(projectID, issueID, opts)
	if err != nil {
		if resp != nil {
			// If the user is already subscribed to the issue, the status code 304 is returned.
			if resp.StatusCode == 304 {
				return nil, errors.New("you are already subscribed to this issue")
			}
		}
		return issue, err
	}

	return issue, nil
}

var UnsubscribeFromIssue = func(client *gitlab.Client, projectID interface{}, issueID int, opts gitlab.RequestOptionFunc) (*gitlab.Issue, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	issue, resp, err := client.Issues.UnsubscribeFromIssue(projectID, issueID, opts)
	if err != nil {
		if resp != nil {
			// If the user is not subscribed to the issue, the status code 304 is returned.
			if resp.StatusCode == 304 {
				return nil, errors.New("you are not subscribed to this issue")
			}
		}
		return issue, err
	}

	return issue, nil
}

var LinkIssues = func(client *gitlab.Client, projectID interface{}, issueIDD int, opts *gitlab.CreateIssueLinkOptions) (*gitlab.Issue, *gitlab.Issue, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	issueLink, _, err := client.IssueLinks.CreateIssueLink(projectID, issueIDD, opts)
	if err != nil {
		return nil, nil, err
	}

	return issueLink.SourceIssue, issueLink.TargetIssue, nil
}

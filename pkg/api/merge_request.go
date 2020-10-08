package api

import (
	"github.com/xanzy/go-gitlab"
)

var ApproveMR = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.ApproveMergeRequestOptions) (*gitlab.MergeRequestApprovals, error) {
	if client == nil {
		client = apiClient
	}
	mr, _, err := client.MergeRequestApprovals.ApproveMergeRequest(projectID, mrID, opts)
	if err != nil {
		return nil, err
	}

	return mr, nil
}

var GetMRApprovalState = func(client *gitlab.Client, projectID interface{}, mrID int, opts ...gitlab.RequestOptionFunc) (*gitlab.MergeRequestApprovalState, error) {
	if client == nil {
		client = apiClient
	}
	mrApprovals, _, err := client.MergeRequestApprovals.GetApprovalState(projectID, mrID, opts...)
	if err != nil {
		return nil, err
	}

	return mrApprovals, nil
}

var GetMR = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.GetMergeRequestsOptions) (*gitlab.MergeRequest, error) {
	if client == nil {
		client = apiClient
	}
	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrID, opts)
	if err != nil {
		return nil, err
	}

	return mr, nil
}

var ListMRs = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
	if client == nil {
		client = apiClient
	}
	mrs, _, err := client.MergeRequests.ListProjectMergeRequests(projectID, opts)
	if err != nil {
		return nil, err
	}

	return mrs, nil
}

var ListMRsWithAssignees = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListProjectMergeRequestsOptions, assigneeIds []int) ([]*gitlab.MergeRequest, error) {
	if client == nil {
		client = apiClient
	}
	mrs := make([]*gitlab.MergeRequest, 0)
	for _, id := range assigneeIds {
		opts.AssigneeID = gitlab.Int(id)
		assingeMrs, _, err := client.MergeRequests.ListProjectMergeRequests(projectID, opts)
		if err != nil {
			return nil, err
		}
		mrs = append(mrs, assingeMrs...)
	}
	return mrs, nil
}

var UpdateMR = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.UpdateMergeRequestOptions) (*gitlab.MergeRequest, error) {
	if client == nil {
		client = apiClient
	}
	mr, _, err := client.MergeRequests.UpdateMergeRequest(projectID, mrID, opts)
	if err != nil {
		return nil, err
	}

	return mr, nil
}

var DeleteMR = func(client *gitlab.Client, projectID interface{}, mrID int) error {
	if client == nil {
		client = apiClient
	}
	_, err := client.MergeRequests.DeleteMergeRequest(projectID, mrID)
	if err != nil {
		return err
	}

	return nil
}

var MergeMR = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.AcceptMergeRequestOptions) (*gitlab.MergeRequest, error) {
	if client == nil {
		client = apiClient
	}
	mrs, _, err := client.MergeRequests.AcceptMergeRequest(projectID, mrID, opts)
	if err != nil {
		return nil, err
	}

	return mrs, nil
}

var CreateMR = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateMergeRequestOptions) (*gitlab.MergeRequest, error) {
	if client == nil {
		client = apiClient
	}
	mr, _, err := client.MergeRequests.CreateMergeRequest(projectID, opts)
	if err != nil {
		return nil, err
	}

	return mr, nil
}

var GetMRLinkedIssues = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.GetIssuesClosedOnMergeOptions) ([]*gitlab.Issue, error) {
	if client == nil {
		client = apiClient
	}
	mrIssues, _, err := client.MergeRequests.GetIssuesClosedOnMerge(projectID, mrID, opts)
	if err != nil {
		return nil, err
	}

	return mrIssues, nil
}

var CreateMRNote = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.CreateMergeRequestNoteOptions) (*gitlab.Note, error) {
	if client == nil {
		client = apiClient
	}

	note, _, err := client.Notes.CreateMergeRequestNote(projectID, mrID, opts)
	if err != nil {
		return note, err
	}

	return note, nil
}

var ListMRNotes = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.ListMergeRequestNotesOptions) ([]*gitlab.Note, error) {
	if client == nil {
		client = apiClient
	}

	notes, _, err := client.Notes.ListMergeRequestNotes(projectID, mrID, opts)
	if err != nil {
		return notes, err
	}

	return notes, nil
}

var RebaseMR = func(client *gitlab.Client, projectID interface{}, mrID int) error {
	if client == nil {
		client = apiClient
	}

	_, err := client.MergeRequests.RebaseMergeRequest(projectID, mrID)
	if err != nil {
		return err
	}

	return nil
}

var UnapproveMR = func(client *gitlab.Client, projectID interface{}, mrID int) error {
	if client == nil {
		client = apiClient
	}

	_, err := client.MergeRequestApprovals.UnapproveMergeRequest(projectID, mrID)
	if err != nil {
		return err
	}

	return nil
}

var SubscribeToMR = func(client *gitlab.Client, projectID interface{}, mrID int, opts gitlab.RequestOptionFunc) (*gitlab.MergeRequest, error) {
	if client == nil {
		client = apiClient
	}

	mr, _, err := client.MergeRequests.SubscribeToMergeRequest(projectID, mrID, opts)
	if err != nil {
		return nil, err
	}

	return mr, nil
}

var UnsubscribeFromMR = func(client *gitlab.Client, projectID interface{}, mrID int, opts gitlab.RequestOptionFunc) (*gitlab.MergeRequest, error) {
	if client == nil {
		client = apiClient
	}

	mr, _, err := client.MergeRequests.UnsubscribeFromMergeRequest(projectID, mrID, opts)
	if err != nil {
		return nil, err
	}

	return mr, nil
}

var MRTodo = func(client *gitlab.Client, projectID interface{}, mrID int, opts gitlab.RequestOptionFunc) (*gitlab.Todo, error) {
	if client == nil {
		client = apiClient
	}

	mr, _, err := client.MergeRequests.CreateTodo(projectID, mrID, opts)
	if err != nil {
		return nil, err
	}

	return mr, nil
}

package api

import "github.com/xanzy/go-gitlab"

var GetCommitStatuses = func(client *gitlab.Client, pid interface{}, sha string) ([]*gitlab.CommitStatus, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	opt := &gitlab.GetCommitStatusesOptions{
		All: gitlab.Bool(true),
	}

	statuses, _, err := client.Commits.GetCommitStatuses(pid, sha, opt, nil)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

package mrutils

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/xanzy/go-gitlab"
)

func Test_DisplayMR(t *testing.T) {
	testCases := []struct {
		name   string
		mr     *gitlab.MergeRequest
		output string
	}{
		{
			name: "opened",
			mr: &gitlab.MergeRequest{
				IID:          1,
				State:        "opened",
				Title:        "This is open",
				SourceBranch: "trunk",
				WebURL:       "https://gitlab.com/profclems/glab/-/merge_requests/1",
			},
			output: "!1 This is open (trunk)\n https://gitlab.com/profclems/glab/-/merge_requests/1\n",
		},
		{
			name: "merged",
			mr: &gitlab.MergeRequest{
				IID:          2,
				State:        "merged",
				Title:        "This is merged",
				SourceBranch: "trunk",
				WebURL:       "https://gitlab.com/profclems/glab/-/merge_requests/2",
			},
			output: "!2 This is merged (trunk)\n https://gitlab.com/profclems/glab/-/merge_requests/2\n",
		},
		{
			name: "closed",
			mr: &gitlab.MergeRequest{
				IID:          3,
				State:        "closed",
				Title:        "This is closed",
				SourceBranch: "trunk",
				WebURL:       "https://gitlab.com/profclems/glab/-/merge_requests/3",
			},
			output: "!3 This is closed (trunk)\n https://gitlab.com/profclems/glab/-/merge_requests/3\n",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := DisplayMR(tC.mr)
			assert.Equal(t, tC.output, got)
		})
	}
}

func Test_MRCheckErrors(t *testing.T) {
	testCases := []struct {
		name    string
		mr      *gitlab.MergeRequest
		errOpts MRCheckErrOptions
		output  string
	}{
		{
			name: "draft",
			mr: &gitlab.MergeRequest{
				IID:            1,
				WorkInProgress: true,
			},
			errOpts: MRCheckErrOptions{
				WorkInProgress: true,
			},
			output: "this merge request is still a work in progress. Run `glab mr update 1 --ready` to mark it as ready for review",
		},
		{
			name: "pipeline",
			mr: &gitlab.MergeRequest{
				IID:                       1,
				MergeWhenPipelineSucceeds: true,
				Pipeline: &gitlab.PipelineInfo{
					Status: "failure",
				},
			},
			errOpts: MRCheckErrOptions{
				PipelineStatus: true,
			},
			output: "pipeline for this merge request has failed. Pipeline is required to succeed before merging",
		},
		{
			name: "merged",
			mr: &gitlab.MergeRequest{
				IID:   1,
				State: "merged",
			},
			errOpts: MRCheckErrOptions{
				Merged: true,
			},
			output: "this merge request has already been merged",
		},
		{
			name: "closed",
			mr: &gitlab.MergeRequest{
				IID:   1,
				State: "closed",
			},
			errOpts: MRCheckErrOptions{
				Closed: true,
			},
			output: "this merge request has been closed",
		},
		{
			name: "opened",
			mr: &gitlab.MergeRequest{
				IID:   1,
				State: "opened",
			},
			errOpts: MRCheckErrOptions{
				Opened: true,
			},
			output: "this merge request is already open",
		},
		{
			name: "subscribed",
			mr: &gitlab.MergeRequest{
				IID:        1,
				Subscribed: true,
			},
			errOpts: MRCheckErrOptions{
				Subscribed: true,
			},
			output: "you are already subscribed to this merge request",
		},
		{
			name: "unsubscribed",
			mr: &gitlab.MergeRequest{
				IID:        1,
				Subscribed: false,
			},
			errOpts: MRCheckErrOptions{
				Unsubscribed: true,
			},
			output: "you are not subscribed to this merge request",
		},
		{
			name: "merge-privilege",
			mr: &gitlab.MergeRequest{
				IID: 1,
				User: struct {
					CanMerge bool "json:\"can_merge\""
				}{CanMerge: false},
			},
			errOpts: MRCheckErrOptions{
				MergePrivilege: true,
			},
			output: "you do not have enough priviledges to merge this merge request",
		},
		{
			name: "conflicts",
			mr: &gitlab.MergeRequest{
				IID:          1,
				HasConflicts: true,
			},
			errOpts: MRCheckErrOptions{
				Conflict: true,
			},
			output: "there are merge conflicts. Resolve conflicts and try again or merge locally",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			err := MRCheckErrors(tC.mr, tC.errOpts)
			assert.EqualError(t, err, tC.output)
		})
	}

	t.Run("nil", func(t *testing.T) {
		err := MRCheckErrors(&gitlab.MergeRequest{}, MRCheckErrOptions{})
		assert.Nil(t, err)
	})
}

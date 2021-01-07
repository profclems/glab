package mrutils

import (
	"errors"
	"testing"

	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/api"
	"github.com/profclems/glab/pkg/prompt"
	"github.com/stretchr/testify/assert"
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

func Test_GetOpenMRForBranchFails(t *testing.T) {
	baseRepo := glrepo.NewWithHost("foo", "bar", "gitlab.com")

	t.Run("API-call-failed", func(t *testing.T) {
		api.ListMRs = func(_ *gitlab.Client, _ interface{}, _ *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
			return nil, errors.New("API call failed")
		}

		got, err := GetOpenMRForBranch(&gitlab.Client{}, baseRepo, "foo")
		assert.Nil(t, got)
		assert.EqualError(t, err, `failed to get open merge request for "foo": API call failed`)
	})

	t.Run("no-return", func(t *testing.T) {
		api.ListMRs = func(_ *gitlab.Client, _ interface{}, _ *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
			return []*gitlab.MergeRequest{}, nil
		}

		got, err := GetOpenMRForBranch(&gitlab.Client{}, baseRepo, "foo")
		assert.Nil(t, got)
		assert.EqualError(t, err, `no open merge request available for "foo"`)
	})

	t.Run("owner-no-match", func(t *testing.T) {
		api.ListMRs = func(_ *gitlab.Client, _ interface{}, _ *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
			return []*gitlab.MergeRequest{
				{
					IID: 1,
					Author: &gitlab.BasicUser{
						Username: "profclems",
					},
				},
				{
					IID: 2,
					Author: &gitlab.BasicUser{
						Username: "maxice8",
					},
				},
			}, nil
		}

		got, err := GetOpenMRForBranch(&gitlab.Client{}, baseRepo, "zemzale:foo")
		assert.Nil(t, got)
		assert.EqualError(t, err, `no open merge request available for "foo" owned by @zemzale`)
	})
}

func Test_GetOpenMRForBranch(t *testing.T) {
	baseRepo := glrepo.NewWithHost("foo", "bar", "gitlab.com")

	testCases := []struct {
		name   string
		input  string
		mrs    []*gitlab.MergeRequest
		expect *gitlab.MergeRequest
	}{
		{
			name: "one-match",
			mrs: []*gitlab.MergeRequest{
				{
					IID: 1,
					Author: &gitlab.BasicUser{
						Username: "profclems",
					},
				},
			},
			expect: &gitlab.MergeRequest{
				IID: 1,
				Author: &gitlab.BasicUser{
					Username: "profclems",
				},
			},
		},
		{
			name:  "owner-match",
			input: "maxice8:foo",
			mrs: []*gitlab.MergeRequest{
				{
					IID: 1,
					Author: &gitlab.BasicUser{
						Username: "profclems",
					},
				},
				{
					IID: 2,
					Author: &gitlab.BasicUser{
						Username: "maxice8",
					},
				},
			},
			expect: &gitlab.MergeRequest{
				IID: 2,
				Author: &gitlab.BasicUser{
					Username: "maxice8",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			api.ListMRs = func(_ *gitlab.Client, _ interface{}, _ *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
				return tC.mrs, nil
			}

			got, err := GetOpenMRForBranch(&gitlab.Client{}, baseRepo, tC.input)
			assert.NoError(t, err)

			assert.Equal(t, tC.expect.IID, got.IID)
			assert.Equal(t, tC.expect.Author.Username, got.Author.Username)
		})
	}
}

func Test_GetOpenMRForBranchPrompt(t *testing.T) {
	baseRepo := glrepo.NewWithHost("foo", "bar", "gitlab.com")

	api.ListMRs = func(_ *gitlab.Client, _ interface{}, _ *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
		return []*gitlab.MergeRequest{
			{
				IID: 1,
				Author: &gitlab.BasicUser{
					Username: "profclems",
				},
			},
			{
				IID: 2,
				Author: &gitlab.BasicUser{
					Username: "maxice8",
				},
			},
		}, nil
	}

	t.Run("success", func(t *testing.T) {
		as, restoreAsk := prompt.InitAskStubber()
		defer restoreAsk()

		as.Stub([]*prompt.QuestionStub{
			{
				Name:  "mr",
				Value: "!1 (foo) by @profclems",
			},
		})

		got, err := GetOpenMRForBranch(&gitlab.Client{}, baseRepo, "foo")
		assert.NoError(t, err)

		assert.Equal(t, 1, got.IID)
		assert.Equal(t, "profclems", got.Author.Username)
	})

	t.Run("error", func(t *testing.T) {
		as, restoreAsk := prompt.InitAskStubber()
		defer restoreAsk()

		as.Stub([]*prompt.QuestionStub{
			{
				Name:  "mr",
				Value: errors.New("prompt failed"),
			},
		})

		got, err := GetOpenMRForBranch(&gitlab.Client{}, baseRepo, "foo")
		assert.Nil(t, got)
		assert.EqualError(t, err, "a merge request must be picked: prompt failed")
	})
}

package mrutils

import (
	"errors"
	"fmt"
	"testing"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
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
	streams, _, _, _ := iostreams.Test()
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := DisplayMR(streams.Color(), tC.mr)
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

func Test_getMRForBranchFails(t *testing.T) {
	baseRepo := glrepo.NewWithHost("foo", "bar", "gitlab.com")

	t.Run("API-call-failed", func(t *testing.T) {
		api.ListMRs = func(_ *gitlab.Client, _ interface{}, _ *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
			return nil, errors.New("API call failed")
		}

		got, err := getMRForBranch(&gitlab.Client{}, baseRepo, "foo", "opened")
		assert.Nil(t, got)
		assert.EqualError(t, err, `failed to get open merge request for "foo": API call failed`)
	})

	t.Run("no-return", func(t *testing.T) {
		api.ListMRs = func(_ *gitlab.Client, _ interface{}, _ *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
			return []*gitlab.MergeRequest{}, nil
		}

		got, err := getMRForBranch(&gitlab.Client{}, baseRepo, "foo", "opened")
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

		got, err := getMRForBranch(&gitlab.Client{}, baseRepo, "zemzale:foo", "opened")
		assert.Nil(t, got)
		assert.EqualError(t, err, `no open merge request available for "foo" owned by @zemzale`)
	})
}

func Test_getMRForBranch(t *testing.T) {
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

			got, err := getMRForBranch(&gitlab.Client{}, baseRepo, tC.input, "opened")
			assert.NoError(t, err)

			assert.Equal(t, tC.expect.IID, got.IID)
			assert.Equal(t, tC.expect.Author.Username, got.Author.Username)
		})
	}
}

func Test_getMRForBranchPrompt(t *testing.T) {
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

		got, err := getMRForBranch(&gitlab.Client{}, baseRepo, "foo", "opened")
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

		got, err := getMRForBranch(&gitlab.Client{}, baseRepo, "foo", "opened")
		assert.Nil(t, got)
		assert.EqualError(t, err, "a merge request must be picked: prompt failed")
	})
}

func Test_MRFromArgsWithOpts(t *testing.T) {
	// Mock cmdutils.Factory object that can be modified as required to perform certain functions
	f := &cmdutils.Factory{
		HttpClient: func() (*gitlab.Client, error) { return &gitlab.Client{}, nil },
		BaseRepo:   func() (glrepo.Interface, error) { return glrepo.New("foo", "bar"), nil },
		Branch:     func() (string, error) { return "trunk", nil },
	}

	t.Run("success", func(t *testing.T) {
		t.Run("via-ID", func(t *testing.T) {
			f := *f

			api.GetMR = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.GetMergeRequestsOptions) (*gitlab.MergeRequest, error) {
				return &gitlab.MergeRequest{
					IID:          2,
					Title:        "test mr",
					SourceBranch: "trunk",
				}, nil
			}

			expectedRepo, err := f.BaseRepo()
			if err != nil {
				t.Skipf("failed to get base repo: %s", err)
			}

			gotMR, gotRepo, err := MRFromArgs(&f, []string{"2"}, "")
			assert.NoError(t, err)

			assert.Equal(t, expectedRepo.FullName(), gotRepo.FullName())

			assert.Equal(t, 2, gotMR.IID)
			assert.Equal(t, "test mr", gotMR.Title)
			assert.Equal(t, "trunk", gotMR.SourceBranch)
		})
		t.Run("via-name", func(t *testing.T) {
			f := *f

			getMRForBranch = func(apiClient *gitlab.Client, baseRepo glrepo.Interface, arg string, state string) (*gitlab.MergeRequest, error) {
				return &gitlab.MergeRequest{
					IID:          2,
					Title:        "test mr",
					SourceBranch: "trunk",
				}, nil
			}

			api.GetMR = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.GetMergeRequestsOptions) (*gitlab.MergeRequest, error) {
				return &gitlab.MergeRequest{
					IID:          2,
					Title:        "test mr",
					SourceBranch: "trunk",
				}, nil
			}

			expectedRepo, err := f.BaseRepo()
			if err != nil {
				t.Skipf("failed to get base repo: %s", err)
			}

			gotMR, gotRepo, err := MRFromArgs(&f, []string{"foo"}, "")
			assert.NoError(t, err)

			assert.Equal(t, expectedRepo.FullName(), gotRepo.FullName())

			assert.Equal(t, 2, gotMR.IID)
			assert.Equal(t, "test mr", gotMR.Title)
			assert.Equal(t, "trunk", gotMR.SourceBranch)
		})
	})

	t.Run("fail", func(t *testing.T) {
		t.Run("HttpClient", func(t *testing.T) {
			f := *f

			f.HttpClient = func() (*gitlab.Client, error) { return nil, errors.New("failed to create HttpClient") }

			gotMR, gotRepo, err := MRFromArgs(&f, []string{}, "")
			assert.Nil(t, gotMR)
			assert.Nil(t, gotRepo)
			assert.EqualError(t, err, "failed to create HttpClient")
		})
		t.Run("BaseRepo", func(t *testing.T) {
			f := *f

			f.BaseRepo = func() (glrepo.Interface, error) { return nil, errors.New("failed to create glrepo.Interface") }

			gotMR, gotRepo, err := MRFromArgs(&f, []string{}, "")
			assert.Nil(t, gotMR)
			assert.Nil(t, gotRepo)
			assert.EqualError(t, err, "failed to create glrepo.Interface")
		})
		t.Run("Branch", func(t *testing.T) {
			f := *f

			f.Branch = func() (string, error) { return "", errors.New("failed to get Branch") }

			gotMR, gotRepo, err := MRFromArgs(&f, []string{}, "")
			assert.Nil(t, gotMR)
			assert.Nil(t, gotRepo)
			assert.EqualError(t, err, "failed to get Branch")
		})
		t.Run("Invalid-MR-ID", func(t *testing.T) {
			f := *f

			gotMR, gotRepo, err := MRFromArgs(&f, []string{"0"}, "")
			assert.Nil(t, gotMR)
			assert.Nil(t, gotRepo)
			assert.EqualError(t, err, "invalid merge request ID provided")
		})
		t.Run("invalid-name", func(t *testing.T) {
			f := *f

			getMRForBranch = func(apiClient *gitlab.Client, baseRepo glrepo.Interface, arg string, state string) (*gitlab.MergeRequest, error) {
				return nil, fmt.Errorf("no merge requests from branch %q", arg)
			}

			gotMR, gotRepo, err := MRFromArgs(&f, []string{"foo"}, "")
			assert.Nil(t, gotMR)
			assert.Nil(t, gotRepo)
			assert.EqualError(t, err, `no merge requests from branch "foo"`)

		})
		t.Run("api.GetMR", func(t *testing.T) {
			f := *f

			api.GetMR = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.GetMergeRequestsOptions) (*gitlab.MergeRequest, error) {
				return nil, errors.New("API call failed")
			}

			gotMR, gotRepo, err := MRFromArgs(&f, []string{"2"}, "")
			assert.Nil(t, gotMR)
			assert.Nil(t, gotRepo)
			assert.EqualError(t, err, "failed to get merge request 2: API call failed")
		})
	})
}

func Test_DisplayAllMRs(t *testing.T) {
	streams, _, _, _ := iostreams.Test()
	mrs := []*gitlab.MergeRequest{
		{
			IID:          1,
			State:        "opened",
			Title:        "add tests",
			TargetBranch: "trunk",
			SourceBranch: "new-tests",
		},
		{
			IID:          2,
			State:        "merged",
			Title:        "fix bug",
			TargetBranch: "trunk",
			SourceBranch: "new-feature",
		},
		{
			IID:          1,
			State:        "closed",
			Title:        "add new feature",
			TargetBranch: "trunk",
			SourceBranch: "new-tests",
		},
	}

	expected := `!1	add tests	(trunk) ← (new-tests)
!2	fix bug	(trunk) ← (new-feature)
!1	add new feature	(trunk) ← (new-tests)
`

	got := DisplayAllMRs(streams, mrs, "unused")
	assert.Equal(t, expected, got)
}

package unsubscribe

import (
	"fmt"
	"testing"
	"time"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/google/shlex"
	"github.com/stretchr/testify/require"

	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "")
}

func TestNewCmdUnsubscribe(t *testing.T) {
	t.Parallel()
	defer config.StubConfig(`---
hosts:
  gitlab.com:
    username: monalisa
    token: OTOKEN
`, "")()

	io, _, stdout, stderr := iostreams.Test()
	stubFactory, _ := cmdtest.StubFactoryWithConfig("")
	stubFactory.IO = io
	stubFactory.IO.IsaTTY = true
	stubFactory.IO.IsErrTTY = true

	oldUnsubscribeMR := api.UnsubscribeFromMR
	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	api.UnsubscribeFromMR = func(client *gitlab.Client, projectID interface{}, mrID int, opts gitlab.RequestOptionFunc) (*gitlab.MergeRequest, error) {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "expected_err" || mrID == 0 {
			return nil, fmt.Errorf("error expected")
		}
		repo, err := stubFactory.BaseRepo()
		if err != nil {
			return nil, err
		}
		return &gitlab.MergeRequest{
			ID:          1,
			IID:         1,
			Title:       "mrtitile",
			Labels:      gitlab.Labels{"bug", "test"},
			State:       "opened",
			Description: "mrbody",
			Subscribed:  false,
			Author: &gitlab.BasicUser{
				ID:       1,
				Name:     "John Dev Wick",
				Username: "jdwick",
			},
			WebURL:    "https://" + repo.RepoHost() + "/" + repo.FullName() + "/-/merge_requests/1",
			CreatedAt: &timer,
		}, nil
	}

	api.GetMR = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.GetMergeRequestsOptions) (*gitlab.MergeRequest, error) {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "expected_err" {
			return nil, fmt.Errorf("error expected")
		}
		repo, err := stubFactory.BaseRepo()
		if err != nil {
			return nil, err
		}
		return &gitlab.MergeRequest{
			ID:          mrID,
			IID:         mrID,
			Title:       "mrTitle",
			Labels:      gitlab.Labels{"test", "bug"},
			State:       "opened",
			Description: "mrBody",
			Subscribed:  true,
			Author: &gitlab.BasicUser{
				ID:       mrID,
				Name:     "John Dev Wick",
				Username: "jdwick",
			},
			WebURL: fmt.Sprintf("https://%s/%s/-/merge_requests/%d", repo.RepoHost(), repo.FullName(), mrID),
		}, nil
	}

	testCases := []struct {
		Name        string
		Issue       string
		ExpectedMsg []string
		wantErr     bool
	}{
		{
			Name:        "Issue Exists",
			Issue:       "1",
			ExpectedMsg: []string{"- Unsubscribing from Merge Request !1", "✓ You have successfully unsubscribed from merge request !1"},
		},
		{
			Name:        "Issue on another repo",
			Issue:       "1 -R profclems/glab",
			ExpectedMsg: []string{"- Unsubscribing from Merge Request !1", "✓ You have successfully unsubscribed from merge request !1"},
		},
		{
			Name:        "Issue Does Not Exist",
			Issue:       "0",
			ExpectedMsg: []string{"- Unsubscribing from Merge Request !0", "error expected"},
			wantErr:     true,
		},
	}

	cmd := NewCmdUnsubscribe(stubFactory)
	cmd.Flags().StringP("repo", "R", "", "")

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			argv, err := shlex.Split(tc.Issue)
			if err != nil {
				t.Fatal(err)
			}
			cmd.SetArgs(argv)
			_, err = cmd.ExecuteC()
			if err != nil {
				if tc.wantErr {
					require.Error(t, err)
					return
				} else {
					t.Fatal(err)
				}
			}

			out := stripansi.Strip(stdout.String())

			for _, msg := range tc.ExpectedMsg {
				assert.Contains(t, out, msg)
				assert.Equal(t, "", stderr.String())
			}
		})
	}

	api.UnsubscribeFromMR = oldUnsubscribeMR
}

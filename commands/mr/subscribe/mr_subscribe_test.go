package subscribe

import (
	"fmt"
	"testing"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "")
}

func TestNewCmdSubscribe(t *testing.T) {
	t.Parallel()
	defer config.StubConfig(`---
hosts:
  gitlab.com:
    username: monalisa
    token: OTOKEN
`, "")()
	stubFactory, _ := cmdtest.StubFactoryWithConfig("")

	oldSubscribeMR := api.SubscribeToMR
	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	api.SubscribeToMR = func(client *gitlab.Client, projectID interface{}, mrID int, opts gitlab.RequestOptionFunc) (*gitlab.MergeRequest, error) {
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
			Author: &gitlab.BasicUser{
				ID:       1,
				Name:     "John Dev Wick",
				Username: "jdwick",
			},
			WebURL:    "https://" + repo.RepoHost() + "/" + repo.FullName() + "/-/merge_requests/1",
			CreatedAt: &timer,
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
			ExpectedMsg: []string{"- Subscribing to merge request !1", "✓ You have successfully subscribed to merge request !1"},
		},
		{
			Name:  "Issue on another repo",
			Issue: "1 -R profclems/glab",
			ExpectedMsg: []string{"- Subscribing to merge request !1",
				"✓ You have successfully subscribed to merge request !1",
				"https://gitlab.com/profclems/glab/-/merge_requests/1\n",
			},
		},
		{
			Name:        "Issue Does Not Exist",
			Issue:       "0",
			ExpectedMsg: []string{"- Subscribing to merge request !0", "error expected"},
			wantErr:     true,
		},
	}

	cmd := NewCmdSubscribe(stubFactory)
	cmdutils.EnableRepoOverride(cmd, stubFactory)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			output, err := cmdtest.RunCommand(cmd, tc.Issue)

			if tc.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			out := stripansi.Strip(output.String())

			for _, msg := range tc.ExpectedMsg {
				assert.Contains(t, out, msg)
			}
		})
	}

	api.SubscribeToMR = oldSubscribeMR
}

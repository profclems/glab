package close

import (
	"bytes"
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
	"testing"
	"time"

	"github.com/profclems/glab/commands/cmdtest"
	"github.com/stretchr/testify/assert"
)

func Test_issueClose(t *testing.T) {
	t.Parallel()

	oldUpdateIssue := api.UpdateIssue
	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	api.UpdateIssue = func(client *gitlab.Client, projectID interface{}, issueID int, opts *gitlab.UpdateIssueOptions) (*gitlab.Issue, error) {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "expected_err" || issueID == 0 {
			return nil, fmt.Errorf("error expected")
		}
		return &gitlab.Issue{
			ID:          issueID,
			IID:         issueID,
			State:       "closed",
			Description: "Dummy description for issue " + string(rune(issueID)),
			Author: &gitlab.IssueAuthor{
				ID:       1,
				Name:     "John Dev Wick",
				Username: "jdwick",
			},
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
			ExpectedMsg: []string{"Closing Issue...", "Issue #1 closed"},
		},
		{
			Name:        "Issue Does Not Exist",
			Issue:       "0",
			ExpectedMsg: []string{"Closing Issue", "404 Not found"},
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var stderr bytes.Buffer
			var stdout bytes.Buffer
			cmd := NewCmdClose(cmdtest.StubFactory())
			cmd.SetArgs([]string{tc.Issue})
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)

			_, err := cmd.ExecuteC()
			if tc.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			out := stripansi.Strip(stdout.String())
			//outErr := stripansi.Strip(stderr.String())

			for _, msg := range tc.ExpectedMsg {
				assert.Contains(t, out, msg)
			}
		})
	}

	api.UpdateIssue = oldUpdateIssue
}

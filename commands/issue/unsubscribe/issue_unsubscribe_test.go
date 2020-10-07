package unsubscribe

import (
	"fmt"
	"testing"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
)

func TestNewCmdUnsubscribe(t *testing.T) {
	t.Parallel()

	oldUnsubscribeIssue := api.UnsubscribeFromIssue
	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	api.UnsubscribeFromIssue = func(client *gitlab.Client, projectID interface{}, issueID int, opts gitlab.RequestOptionFunc) (*gitlab.Issue, error) {
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
			ExpectedMsg: []string{"- Unsubscribing from Issue #1", "✔ Unsubscribed from issue #1"},
		},
		{
			Name:        "Issue on another repo",
			Issue:       "1 -R profclems/glab",
			ExpectedMsg: []string{"- Unsubscribing from Issue #1", "✔ Unsubscribed from issue #1"},
		},
		{
			Name:        "Issue Does Not Exist",
			Issue:       "0",
			ExpectedMsg: []string{"- Unsubscribing from Issue #0", "error expected"},
			wantErr:     true,
		},
	}

	cmd := NewCmdUnsubscribe(cmdtest.StubFactory("https://gitlab.com/glab-cli/test"))
	cmd.Flags().StringP("repo", "R", "", "")

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

	api.UnsubscribeFromIssue = oldUnsubscribeIssue
}

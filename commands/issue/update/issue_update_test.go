package update

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/profclems/glab/internal/utils"

	"github.com/google/shlex"

	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
)

func TestNewCmdUpdate(t *testing.T) {
	t.Parallel()

	oldUpdateIssue := api.UpdateIssue
	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	testIssue := &gitlab.Issue{
		ID:               1,
		IID:              1,
		State:            "closed",
		Labels:           gitlab.Labels{"bug, test, removeable-label"},
		Description:      "Dummy description for issue 1",
		DiscussionLocked: false,
		Author: &gitlab.IssueAuthor{
			ID:       1,
			Name:     "John Dev Wick",
			Username: "jdwick",
		},
		CreatedAt: &timer,
	}
	api.UpdateIssue = func(client *gitlab.Client, projectID interface{}, issueID int, opts *gitlab.UpdateIssueOptions) (*gitlab.Issue, error) {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "expected_err" || issueID != testIssue.ID {
			return nil, fmt.Errorf("error expected")
		}
		if *opts.Title != "" {
			testIssue.Title = *opts.Title
		}
		if *opts.Description != "" {
			testIssue.Description = *opts.Description
		}
		if opts.AddLabels != nil {
			testIssue.Labels = opts.AddLabels
		}
		return testIssue, nil
	}

	testCases := []struct {
		Name        string
		Issue       string
		ExpectedMsg []string
		wantErr     bool
	}{
		{
			Name:  "Issue Exists",
			Issue: `1 -t "New Title" -d "A new description" --lock-discussion -l newLabel --unlabel bug`,
			ExpectedMsg: []string{
				"- Updating issue #1",
				"✓ updated title to \"New Title\"",
				"✓ locked discussion",
				"✓ added labels newLabel",
				"✓ removed labels bug",
				"#1 New Title",
			},
		},
		{
			Name:  "Issue Exists on different repo",
			Issue: `1 -R glab_cli/test`,
			ExpectedMsg: []string{
				"- Updating issue #1",
				"✓ updated title to \"New Title\"",
				"✓ locked discussion",
				"✓ added labels newLabel",
				"✓ removed labels bug",
				"#1 New Title",
			},
		},
		{
			Name:        "Issue Does Not Exist",
			Issue:       "0",
			ExpectedMsg: []string{"- Updating issue #0", "error expected"},
			wantErr:     true,
		},
	}

	io, _, stdout, stderr := utils.IOTest()
	f := cmdtest.StubFactory("https://gitlab.com/glab-cli/test")
	f.IO = io
	f.IO.IsaTTY = true
	f.IO.IsErrTTY = true

	cmd := NewCmdUpdate(f)
	cmd.Flags().StringP("repo", "R", "", "")

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			args, _ := shlex.Split(tc.Issue)
			cmd.SetArgs(args)
			cmd.SetOut(ioutil.Discard)
			cmd.SetErr(ioutil.Discard)

			_, err := cmd.ExecuteC()
			if tc.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			out := stripansi.Strip(stdout.String())
			outErr := stripansi.Strip(stderr.String())

			for _, msg := range tc.ExpectedMsg {
				assert.Contains(t, out, msg)
				assert.Contains(t, outErr, "")
			}
		})
	}

	api.UpdateIssue = oldUpdateIssue
}

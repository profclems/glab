package delete

import (
	"fmt"
	"strings"
	"testing"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "issue_delete_test")
}

func TestNewCmdDelete(t *testing.T) {
	t.Parallel()
	oldDeleteMR := api.DeleteMR

	api.DeleteIssue = func(client *gitlab.Client, projectID interface{}, issueID int) error {
		if projectID == "" || projectID == "NAMESPACE/WRONG_REPO" || projectID == "expected_err" || issueID == 0 {
			return fmt.Errorf("error expected")
		}
		return nil
	}
	api.GetIssue = func(client *gitlab.Client, projectID interface{}, issueID int) (*gitlab.Issue, error) {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "expected_err" {
			return nil, fmt.Errorf("error expected")
		}
		return &gitlab.Issue{
			IID: issueID,
		}, nil
	}

	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		errMsg     string
		assertFunc func(*testing.T, string, string)
	}{
		{
			name:    "delete",
			args:    []string{"0", "-R", "NAMESPACE/WRONG_REPO"},
			wantErr: true,
		},
		{
			name:    "id exists",
			args:    []string{"1"},
			wantErr: false,
			assertFunc: func(t *testing.T, out string, err string) {
				assert.Contains(t, err, "✓ Issue Deleted\n")
			},
		},
		{
			name:    "delete on different repo",
			args:    []string{"12", "-R", "profclems/glab"},
			wantErr: false,
			assertFunc: func(t *testing.T, out string, stderr string) {
				assert.Contains(t, stderr, "✓ Issue Deleted\n")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			io, _, stdout, stderr := iostreams.Test()
			f := cmdtest.StubFactory("")
			f.IO = io
			f.IO.IsaTTY = true
			f.IO.IsErrTTY = true

			cmd := NewCmdDelete(f)

			cmd.Flags().StringP("repo", "R", "", "")

			cli := strings.Join(tt.args, " ")
			t.Log(cli)
			_, err := cmdtest.RunCommand(cmd, cli)
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.assertFunc(t, stdout.String(), stderr.String())
			} else {
				assert.NotNil(t, err)
			}
		})
	}

	api.DeleteMR = oldDeleteMR
}

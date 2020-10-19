package delete

import (
	"fmt"
	"strings"
	"testing"

	"github.com/profclems/glab/internal/utils"

	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

// TODO: test by mocking the appropriate api function
func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "mr_delete_test")
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

			assertFunc: func(t *testing.T, out string, err string) {
				assert.Contains(t, err, "error expected")
			},
		},
		{
			name:    "id exists",
			args:    []string{"1"},
			wantErr: false,
			assertFunc: func(t *testing.T, out string, err string) {
				assert.Contains(t, out, "✓ Issue Deleted\n")
			},
		},
		{
			name:    "delete on different repo",
			args:    []string{"12", "-R", "profclems/glab"},
			wantErr: false,
			assertFunc: func(t *testing.T, out string, err string) {
				assert.Contains(t, out, "✓ Issue Deleted\n")
			},
		},
	}

	io, _, stdout, stderr := utils.IOTest()
	f := cmdtest.StubFactory("")
	f.IO = io
	f.IO.IsaTTY = true
	f.IO.IsErrTTY = true

	cmd := NewCmdDelete(f)

	cmd.Flags().StringP("repo", "R", "", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cli := strings.Join(tt.args, " ")
			t.Log(cli)
			_, err := cmdtest.RunCommand(cmd, cli)
			if !tt.wantErr {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				stderr.WriteString(err.Error()) // write err to stderr
			}

			out := stripansi.Strip(stdout.String())
			outErr := stripansi.Strip(stderr.String())

			tt.assertFunc(t, out, outErr)
			assert.Contains(t, outErr, tt.errMsg)
			stderr.Reset()
			stdout.Reset()
		})
	}

	api.DeleteMR = oldDeleteMR
}

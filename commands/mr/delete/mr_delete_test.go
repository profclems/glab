package delete

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/shlex"

	"github.com/profclems/glab/internal/utils"

	"github.com/profclems/glab/internal/config"

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

func Test_deleteMergeRequest(t *testing.T) {
	defer config.StubConfig(`---
hosts:
  gitlab.com:
    username: monalisa
    token: OTOKEN
`, "")()
	io, _, stdout, stderr := utils.IOTest()
	stubFactory, _ := cmdtest.StubFactoryWithConfig("")
	stubFactory.IO = io
	stubFactory.IO.IsaTTY = true
	stubFactory.IO.IsErrTTY = true
	oldDeleteMR := api.DeleteMR

	api.DeleteMR = func(client *gitlab.Client, projectID interface{}, mrID int) error {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "expected_err" || mrID == 0 {
			return fmt.Errorf("error expected")
		}
		return nil
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
			Author: &gitlab.BasicUser{
				ID:       mrID,
				Name:     "John Dev Wick",
				Username: "jdwick",
			},
			WebURL: fmt.Sprintf("https://%s/%s/-/merge_requests/%d", repo.RepoHost(), repo.FullName(), mrID),
		}, nil
	}

	api.ListMRs = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
		return []*gitlab.MergeRequest{}, nil
	}

	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		errMsg     string
		assertFunc func(*testing.T, string, string, error)
	}{
		{
			name:    "delete",
			args:    []string{"0"},
			wantErr: true,

			assertFunc: func(t *testing.T, out, outErr string, err error) {
				assert.Equal(t, "invalid merge request ID provided", err.Error())
			},
		},
		{
			name:    "id exists",
			args:    []string{"1"},
			wantErr: false,
			assertFunc: func(t *testing.T, out, outErr string, err error) {
				assert.Contains(t, out, "- Deleting Merge Request !1\n")
				assert.Contains(t, out, "✔ Merge request !1 deleted\n")
			},
		},
		{
			name:    "delete on different repo",
			args:    []string{"1", "-R", "profclems/glab"},
			wantErr: false,
			assertFunc: func(t *testing.T, out, outErr string, err error) {
				assert.Contains(t, out, "- Deleting Merge Request !1\n")
				assert.Contains(t, out, "✔ Merge request !1 deleted\n")
			},
		},
		{
			name:    "delete no args",
			wantErr: true,
			assertFunc: func(t *testing.T, out, outErr string, err error) {
				assert.Equal(t, `no open merge request available for "master"`, err.Error())
			},
		},
	}

	cmd := NewCmdDelete(cmdtest.StubFactory(""))

	cmd.Flags().StringP("repo", "R", "", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cli := strings.Join(tt.args, " ")
			t.Log(cli)
			argv, err := shlex.Split(cli)
			if err != nil {
				t.Fatal(err)
			}
			cmd.SetArgs(argv)
			_, err = cmd.ExecuteC()
			if !tt.wantErr {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}

			out := stripansi.Strip(stdout.String())
			outErr := stripansi.Strip(stderr.String())
			t.Log(outErr)

			tt.assertFunc(t, out, outErr, err)
			assert.Contains(t, outErr, tt.errMsg)
		})
	}

	api.DeleteMR = oldDeleteMR
}

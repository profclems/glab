package delete

import (
	"fmt"
	"strings"
	"testing"

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
	t.Parallel()
	stubFactory, _ := cmdtest.StubFactoryWithConfig("")
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
		assertFunc func(t *testing.T, out string)
	}{
		{
			name:    "delete",
			args:    []string{"0"},
			wantErr: true,

			assertFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "invalid merge request ID provided")
			},
		},
		{
			name:    "id exists",
			args:    []string{"1"},
			wantErr: false,
			assertFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "- Deleting Merge Request !1\n")
				assert.Contains(t, out, "✔ Merge request !1 deleted\n")
			},
		},
		{
			name:    "delete on different repo",
			args:    []string{"1", "-R", "profclems/glab"},
			wantErr: false,
			assertFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "- Deleting Merge Request !1\n")
				assert.Contains(t, out, "✔ Merge request !1 deleted\n")
			},
		},
		{
			name:    "delete no args",
			wantErr: true,
			assertFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "no open merge request availabe for \"master\"")
			},
		},
	}

	cmd := NewCmdDelete(cmdtest.StubFactory(""))

	cmd.Flags().StringP("repo", "R", "", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cli := strings.Join(tt.args, " ")
			t.Log(cli)
			output, err := cmdtest.RunCommand(cmd, cli)
			if !tt.wantErr {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}

			out := stripansi.Strip(output.String())
			outErr := stripansi.Strip(output.Stderr())

			tt.assertFunc(t, out)
			assert.Contains(t, outErr, tt.errMsg)
		})
	}

	api.DeleteMR = oldDeleteMR
}

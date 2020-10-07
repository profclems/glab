package create

import (
	"fmt"
	"testing"

	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
)

func TestNewCmdCreate(t *testing.T) {
	api.CreateIssueBoard = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateIssueBoardOptions) (*gitlab.IssueBoard, error) {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "NS/WRONG_REPO" {
			return nil, fmt.Errorf("error expected")
		}

		return &gitlab.IssueBoard{
			ID:        11,
			Name:      *opts.Name,
			Project:   &gitlab.Project{PathWithNamespace: projectID.(string)},
			Milestone: nil,
			Lists:     nil,
		}, nil
	}
	tests := []struct {
		name    string
		arg     string
		want    string
		wantErr bool
	}{
		{
			name: "Name passed as arg",
			arg:  `"Test"`,
			want: `✓ Board created: "Test"`,
		},
		{
			name: "Name passed in name flag",
			arg:  `--name "Test"`,
			want: `✓ Board created: "Test"`,
		},
		{
			name:    "WRONG_REPO",
			arg:     `"Test" -R NS/WRONG_REPO`,
			wantErr: true,
		},
	}

	cmd := NewCmdCreate(cmdtest.StubFactory("https://gitlab.com/glab-cli/test"))
	cmd.Flags().StringP("repo", "R", "", "")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := cmdtest.RunCommand(cmd, tc.arg)
			if tc.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			out := stripansi.Strip(output.String())

			assert.Contains(t, out, tc.want)

		})
	}
}

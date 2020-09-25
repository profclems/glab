package archive

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_repoArchive(t *testing.T) {
	t.Parallel()
	repo := copyTestRepo(t)

	type argFlags struct {
		format string
		sha    string
		repo   string
		dest   string
	}

	tests := []struct {
		name    string
		args    argFlags
		wantMsg []string
		wantErr bool
	}{
		{
			name:    "Has invalid format",
			args:    argFlags{"asp", "master", "glab-cli/test", "test"},
			wantMsg: []string{"format must be one of"},
			wantErr: true,
		},
		{
			name:    "Has valid format",
			args:    argFlags{"zip", "master", "glab-cli/test", "test"},
			wantMsg: []string{"Cloning...", "Complete... test.zip"},
		},
		{
			name:    "Repo is invalid",
			args:    argFlags{"zip", "master", "glab-cli/testzz", "test"},
			wantMsg: []string{"404 Project Not Found"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(glabBinaryPath, "repo", "archive", tt.args.repo, tt.args.dest, "--format", tt.args.format, "--sha", tt.args.sha)
			cmd.Dir = repo
			b, err := cmd.CombinedOutput()
			if err != nil && !tt.wantErr {
				t.Log(string(b))
				t.Fatal(err)
			}
			out := string(b)
			t.Log(out)

			for _, msg := range tt.wantMsg {
				assert.Contains(t, out, msg)
			}
		})
	}
}

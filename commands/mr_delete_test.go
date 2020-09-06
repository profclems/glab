package commands

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func Test_deleteMergeRequest(t *testing.T) {
	t.Parallel()
	repo := copyTestRepo(t)
	var cmd *exec.Cmd
	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		assertFunc func(t *testing.T, out string)
	}{
		{
			name:    "delete",
			args:    []string{"0"},
			wantErr: true,
			assertFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "404 Not Found")
			},
		},
		{
			name:    "delete no args",
			wantErr: true,
			assertFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "accepts 1 arg(s), received 0")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd = exec.Command(glabBinaryPath, tt.args...)
			cmd.Dir = repo
			b, err := cmd.CombinedOutput()
			if err != nil && !tt.wantErr {
				t.Log(string(b))
				t.Error(err)
			}
			out := string(b)
			t.Log(out)
		})
	}
}

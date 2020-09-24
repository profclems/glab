package mr

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_mrNoteCreate(t *testing.T) {
	repo := copyTestRepo(t)
	var cmd *exec.Cmd

	tests := []struct {
		name          string
		args          []string
		want          bool
		assertionFunc func(t *testing.T, out string)
	}{
		{
			name: "Has -m flag",
			args: []string{"225", "-m", "Some test note"},
			assertionFunc: func(t *testing.T, out string) {
				require.Contains(t, out, "https://gitlab.com/glab-cli/test/merge_requests/1#note_")
			},
		},
		{
			name: "Has no flag",
			args: []string{"225"},
			assertionFunc: func(t *testing.T, out string) {
				require.Contains(t, out, "aborted... Note has an empty message")
			},
		},
		{
			name: "With --repo flag",
			args: []string{"225", "-m", "Some test note", "-R", "glab-cli/test"},
			assertionFunc: func(t *testing.T, out string) {
				require.Contains(t, out, "https://gitlab.com/glab-cli/test/merge_requests/1#note_")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd = exec.Command(glabBinaryPath, append([]string{"mr", "note"}, tt.args...)...)
			cmd.Dir = repo

			b, err := cmd.CombinedOutput()
			if err != nil {
				t.Log(string(b))
				t.Fatal(err)
			}
		})
	}
}

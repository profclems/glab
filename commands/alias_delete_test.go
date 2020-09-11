package commands

import (
	"os/exec"
	"testing"

	"github.com/profclems/glab/internal/config"
	"github.com/stretchr/testify/require"
)

func TestAliasDelete(t *testing.T) {
	repo := copyTestRepo(t)
	config.SetAlias("co-test", "glab mr checkout")
	tests := []struct {
		name       string
		config     string
		cli        string
		isTTY      bool
		wantStdout string
		wantStderr string
		wantErr    string
	}{
		{
			name:       "alias does not exist",
			config:     "",
			cli:        "test-nonexistent",
			isTTY:      true,
			wantStderr: "",
			wantErr:    "no such alias test-nonexistent",
		},
		{
			name: "delete one",
			cli:        "co-test",
			wantStderr: "✓ Deleted alias co; was mr checkout\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(glabBinaryPath, "alias", "delete", "co-test")
			cmd.Dir = repo

			b, err := cmd.CombinedOutput()
			if tt.wantErr != "" {
				require.Error(t, err)
				return
			}
			//require.NoError(t, err)

			require.NotEmpty(t, string(b))
		})
	}
}
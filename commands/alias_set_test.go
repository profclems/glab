package commands

import (
	"github.com/profclems/glab/internal/utils"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AliasSetCmd(t *testing.T) {
	t.Parallel()
	repo := copyTestRepo(t)
	var cmd *exec.Cmd

	tests := []struct{
		Name string
		args []string
		wantErr bool
		assertFunc func(t *testing.T, out string)
	}{
		{
			Name: "Alias name is a command name",
			args: []string{"mr", "'mr list'"},
			wantErr: true,
			assertFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "could not create alias: \"mr\" is already a glab command")
			},
		},
		{
			Name: "Is valid",
			args: []string{"mrl", "'mr list'"},
			assertFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, utils.GreenCheck()+" Alias added")
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			cmd = exec.Command(glabBinaryPath, append([]string{"alias", "set"}, test.args...)...)
			cmd.Dir = repo

			b, err := cmd.CombinedOutput()
			if err != nil && !test.wantErr {
				t.Log(string(b))
				t.Fatal(err)
			}
			out := string(b)
			test.assertFunc(t, out)
		})
	}
}
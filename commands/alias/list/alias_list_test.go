package list

import (
	"bytes"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"io/ioutil"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAliasList(t *testing.T) {
	tests := []struct {
		name       string
		config     string
		wantStdout string
		wantStderr string
	}{
		{
			name:       "empty",
			config:     "",
			wantStdout: "",
			wantStderr: "no aliases configured\n",
		},
		{
			name: "some",
			config: heredoc.Doc(`
				aliases:
				  co: mr checkout
				  gc: "!glab mr create -f \"$@\" | pbcopy"
			`),
			wantStdout: "co\tmr checkout                     \ngc\t!glab mr create -f \"$@\" | pbcopy",
			wantStderr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: change underlying config implementation so Write is not
			// automatically called when editing aliases in-memory
			defer config.StubWriteConfig(ioutil.Discard, ioutil.Discard)()

			cfg := config.NewFromString(tt.config)

			var stderr bytes.Buffer
			var stdout bytes.Buffer

			factoryConf := &cmdutils.Factory{
				Config: func() (config.Config, error) {
					return cfg, nil
				},
			}

			cmd := NewCmdList(factoryConf, nil)
			cmd.SetArgs([]string{})

			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)

			_, err := cmd.ExecuteC()
			require.NoError(t, err)

			assert.Equal(t, tt.wantStdout, stdout.String())
			assert.Equal(t, tt.wantStderr, stderr.String())
		})
	}
}

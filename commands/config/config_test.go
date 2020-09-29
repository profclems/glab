package config

import (
	"bytes"
	"errors"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/profclems/glab/internal/config"
	"github.com/stretchr/testify/require"
)

type configStub map[string]string

func (c configStub) Local() (*config.LocalConfig, error) {
	return nil, nil
}

func (c configStub) WriteAll() error {
	c["_written"] = "true"
	return nil
}

func genKey(host, key string) string {
	if host != "" {
		return host + ":" + key
	}
	return key
}

func (c configStub) Get(host, key string) (string, error) {
	val, _, err := c.GetWithSource(host, key)
	return val, err
}

func (c configStub) GetWithSource(host, key string) (string, string, error) {
	if v, found := c[genKey(host, key)]; found {
		return v, "(memory)", nil
	}
	return "", "", errors.New("not found")
}

func (c configStub) Set(host, key, value string) error {
	c[genKey(host, key)] = value
	return nil
}

func (c configStub) Aliases() (*config.AliasConfig, error) {
	return nil, nil
}

func (c configStub) Hosts() ([]string, error) {
	return nil, nil
}

func (c configStub) UnsetHost(hostname string) {
}

func (c configStub) Write() error {
	c["_written"] = "true"
	return nil
}

func TestConfigGet(t *testing.T) {
	tests := []struct {
		name   string
		config configStub
		args   []string
		stdout string
		stderr string
	}{
		{
			name: "get key",
			config: configStub{
				"editor": "ed",
			},
			args:   []string{"editor"},
			stdout: "ed\n",
			stderr: "",
		},
		{
			name: "get key scoped by host",
			config: configStub{
				"editor":            "ed",
				"gitlab.com:editor": "vim",
			},
			args:   []string{"editor", "-h", "gitlab.com"},
			stdout: "vim\n",
			stderr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			var stdout bytes.Buffer
			f := &cmdutils.Factory{
				Config: func() (config.Config, error) {
					return tt.config, nil
				},
			}

			cmd := NewCmdConfigGet(f)
			cmd.Flags().BoolP("help", "x", false, "")
			cmd.SetArgs(tt.args)
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)

			_, err := cmd.ExecuteC()
			require.NoError(t, err)

			assert.Equal(t, tt.stdout, stdout.String())
			assert.Equal(t, tt.stderr, stderr.String())
			assert.Equal(t, "", tt.config["_written"])
		})
	}
}

func TestConfigSet(t *testing.T) {
	tests := []struct {
		name      string
		config    configStub
		args      []string
		expectKey string
		stdout    string
		stderr    string
	}{
		{
			name:      "set key",
			config:    configStub{},
			args:      []string{"editor", "vim"},
			expectKey: "editor",
			stdout:    "",
			stderr:    "",
		},
		{
			name:      "set key scoped by host",
			config:    configStub{},
			args:      []string{"editor", "vim", "-h", "gitlab.com"},
			expectKey: "gitlab.com:editor",
			stdout:    "",
			stderr:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			var stdout bytes.Buffer
			f := &cmdutils.Factory{
				Config: func() (config.Config, error) {
					return tt.config, nil
				},
			}

			cmd := NewCmdConfigSet(f)
			cmd.Flags().BoolP("help", "x", false, "")
			cmd.SetArgs(append(tt.args, "-g"))
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)

			_, err := cmd.ExecuteC()
			require.NoError(t, err)

			assert.Equal(t, tt.stdout, stdout.String())
			assert.Equal(t, tt.stderr, stderr.String())
			assert.Equal(t, "vim", tt.config[tt.expectKey])
			assert.Equal(t, "true", tt.config["_written"])
		})
	}
}

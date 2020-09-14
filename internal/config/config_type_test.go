package config

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsSelfHosted(t *testing.T) {
	type args struct {
		h string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "self_hosted",
			args: args{h: "gitlab.example.com"},
			want: true,
		},
		{
			name: "gitlab.com",
			args: args{h: "gitlab.com"},
			want: true,
		},
		{
			name: "is a gitlab.com subdomain",
			args: args{h: "example.gitlab.com"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSelfHosted(tt.args.h); got != tt.want {
				t.Errorf("IsSelfHosted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileConfig_Set(t *testing.T) {
	mainBuf := bytes.Buffer{}
	hostsBuf := bytes.Buffer{}
	defer StubWriteConfig(&mainBuf, &hostsBuf)()

	c := NewBlankConfig()
	assert.NoError(t, c.Set("", "editor", "nano"))
	assert.NoError(t, c.Set("gitlab.com", "git_protocol", "ssh"))
	assert.NoError(t, c.Set("example.com", "editor", "vim"))
	assert.NoError(t, c.Set("gitlab.com", "user", "hubot"))
	assert.NoError(t, c.Write())

	expected := "# What protocol to use when performing git operations. Supported values: ssh, https\ngit_protocol: https\n# What editor gh should run when creating issues, pull requests, etc. If blank, will refer to environment.\neditor: nano\n# When to interactively prompt. This is a global config that cannot be overriden by hostname. Supported values: enabled, disabled\nprompt: enabled\n# Aliases allow you to create nicknames for gh commands\naliases:\n    co: pr checkout\n"
	assert.Equal(t, expected, mainBuf.String())
	assert.Equal(t, `gitlab.com:
    git_protocol: ssh
    user: hubot
gitlab.example.com:
    editor: vim
`, hostsBuf.String())
}

func Test_defaultConfig(t *testing.T) {
	mainBuf := bytes.Buffer{}
	hostsBuf := bytes.Buffer{}
	defer StubWriteConfig(&mainBuf, &hostsBuf)()

	cfg := NewBlankConfig()
	assert.NoError(t, cfg.Write())

	expected := "# What protocol to use when performing git operations. Supported values: ssh, https\ngit_protocol: https\n# What editor gh should run when creating issues, pull requests, etc. If blank, will refer to environment.\neditor:\n# When to interactively prompt. This is a global config that cannot be overriden by hostname. Supported values: enabled, disabled\nprompt: enabled\n# Aliases allow you to create nicknames for gh commands\naliases:\n    co: pr checkout\n"
	assert.Equal(t, expected, mainBuf.String())
	assert.Equal(t, "", hostsBuf.String())

	proto, err := cfg.Get("", "git_protocol")
	assert.Nil(t, err)
	assert.Equal(t, "https", proto)

	editor, err := cfg.Get("", "editor")
	assert.Nil(t, err)
	assert.Equal(t, "", editor)

	aliases, err := cfg.Aliases()
	assert.Nil(t, err)
	assert.Equal(t, len(aliases.All()), 1)
	expansion, _ := aliases.Get("co")
	assert.Equal(t, expansion, "mr checkout")
}
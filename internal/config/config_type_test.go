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
			want: false,
		},
		{
			name: "is a gitlab.com subdomain",
			args: args{h: "example.gitlab.com"},
			want: false,
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
	assert.NoError(t, c.Set("gitlab.com", "username", "hubot"))
	assert.NoError(t, c.Write())

	expected := "# What protocol to use when performing git operations. Supported values: ssh, https\ngit_protocol: ssh\n# What editor glab should run when creating issues, merge requests, etc.  This is a global config that cannot be overridden by hostname.\neditor: nano\n# What browser glab should run when opening links. This is a global config that cannot be overridden by hostname.\nbrowser:\n# Git remote alias which glab should use when fetching the remote url. This can be overridden by hostname\nremote_alias: origin\n# Set your desired markdown renderer style. Available options are [dark, light, notty] or set a custom style. Refer to https://github.com/charmbracelet/glamour#styles\nglamour_style: dark\n# Allow glab to automatically check for updates and notify you when there are new updates\ncheck_update: false\n# configuration specific for gitlab instances\nhosts:\n    gitlab.com:\n        # What protocol to use to access the api endpoint. Supported values: http, https\n        protocol: https\n        # Your GitLab access token. Get an access token at https://gitlab.com/profile/personal_access_tokens\n        token:\n        git_protocol: ssh\n        username: hubot\n    example.com:\n        editor: vim\n"
	assert.Equal(t, expected, mainBuf.String())
	assert.Equal(t, `gitlab.com:
    git_protocol: ssh
    username: hubot
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

	expected := "# What protocol to use when performing git operations. Supported values: ssh, https\ngit_protocol: ssh\n# What editor glab should run when creating issues, merge requests, etc.  This is a global config that cannot be overridden by hostname.\neditor:\n# What browser glab should run when opening links. This is a global config that cannot be overridden by hostname.\nbrowser:\n# Git remote alias which glab should use when fetching the remote url. This can be overridden by hostname\nremote_alias: origin\n# Set your desired markdown renderer style. Available options are [dark, light, notty] or set a custom style. Refer to https://github.com/charmbracelet/glamour#styles\nglamour_style: dark\n# Allow glab to automatically check for updates and notify you when there are new updates\ncheck_update: false\n# configuration specific for gitlab instances\nhosts:\n    gitlab.com:\n        # What protocol to use to access the api endpoint. Supported values: http, https\n        protocol: https\n        # Your GitLab access token. Get an access token at https://gitlab.com/profile/personal_access_tokens\n        token:\n"
	assert.Equal(t, expected, mainBuf.String())
	assert.Equal(t, "", hostsBuf.String())

	proto, err := cfg.Get("", "git_protocol")
	assert.Nil(t, err)
	assert.Equal(t, "ssh", proto)

	editor, err := cfg.Get("", "editor")
	assert.Nil(t, err)
	assert.Equal(t, "", editor)

	aliases, err := cfg.Aliases()
	assert.Nil(t, err)
	assert.Equal(t, len(aliases.All()), 2)
	expansion, _ := aliases.Get("co")
	assert.Equal(t, expansion, "mr checkout")
}

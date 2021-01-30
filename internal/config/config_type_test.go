package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fileConfig_Set(t *testing.T) {
	mainBuf := bytes.Buffer{}
	aliasesBuf := bytes.Buffer{}
	defer StubWriteConfig(&mainBuf, &aliasesBuf)()

	c := NewBlankConfig()
	assert.NoError(t, c.Set("", "editor", "nano"))
	assert.NoError(t, c.Set("gitlab.com", "git_protocol", "ssh"))
	assert.NoError(t, c.Set("example.com", "editor", "vim"))
	assert.NoError(t, c.Set("gitlab.com", "username", "hubot"))
	assert.NoError(t, c.WriteAll())
	//a, _ := c.Aliases()
	//assert.NoError(t, a.Set("co", "mr checkout"))
	//assert.NoError(t, a.Write())

	expected := "# What protocol to use when performing git operations. Supported values: ssh, https\ngit_protocol: ssh\n# What editor glab should run when creating issues, merge requests, etc.  This is a global config that cannot be overridden by hostname.\neditor: nano\n# What browser glab should run when opening links. This is a global config that cannot be overridden by hostname.\nbrowser:\n# Set your desired markdown renderer style. Available options are [dark, light, notty] or set a custom style. Refer to https://github.com/charmbracelet/glamour#styles\nglamour_style: dark\n# Allow glab to automatically check for updates and notify you when there are new updates\ncheck_update: false\n# configuration specific for gitlab instances\nhosts:\n    gitlab.com:\n        # What protocol to use to access the api endpoint. Supported values: http, https\n        api_protocol: api_host\n        # Configure host for api endpoint, defaults to the host itself\n        https: token\n        # Your GitLab access token. Get an access token at https://gitlab.com/profile/personal_access_tokens\n        '': git_protocol\n        ssh: username\n    example.com:\n        editor: vim\n"
	assert.Equal(t, expected, mainBuf.String())
	assert.Equal(t, `ci: pipeline ci
co: mr checkout
`, aliasesBuf.String())
}

func Test_defaultConfig(t *testing.T) {
	mainBuf := bytes.Buffer{}
	hostsBuf := bytes.Buffer{}
	defer StubWriteConfig(&mainBuf, &hostsBuf)()

	cfg := NewBlankConfig()
	assert.NoError(t, cfg.Write())

	expected := "# What protocol to use when performing git operations. Supported values: ssh, https\ngit_protocol: ssh\n# What editor glab should run when creating issues, merge requests, etc.  This is a global config that cannot be overridden by hostname.\neditor:\n# What browser glab should run when opening links. This is a global config that cannot be overridden by hostname.\nbrowser:\n# Set your desired markdown renderer style. Available options are [dark, light, notty] or set a custom style. Refer to https://github.com/charmbracelet/glamour#styles\nglamour_style: dark\n# Allow glab to automatically check for updates and notify you when there are new updates\ncheck_update: false\n# configuration specific for gitlab instances\nhosts:\n    gitlab.com:\n        # What protocol to use to access the api endpoint. Supported values: http, https\n        api_protocol: api_host\n        # Configure host for api endpoint, defaults to the host itself\n        https: token\n        # Your GitLab access token. Get an access token at https://gitlab.com/profile/personal_access_tokens\n"
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

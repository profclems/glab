package config

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func eq(t *testing.T, got interface{}, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func Test_parseConfig(t *testing.T) {
	defer StubConfig(`---
hosts:
  gitlab.com:
    username: monalisa
    token: OTOKEN
aliases:
`, "")()
	// prevent using env variable for test
	envToken := os.Getenv("GITLAB_TOKEN")
	if envToken != "" {
		_ = os.Setenv("GITLAB_TOKEN", "")
	}
	config, err := ParseConfig("config.yml")
	eq(t, err, nil)
	username, err := config.Get("gitlab.com", "username")
	eq(t, err, nil)
	eq(t, username, "monalisa")
	token, err := config.Get("gitlab.com", "token")
	eq(t, err, nil)
	eq(t, token, "OTOKEN")
	if envToken != "" {
		_ = os.Setenv("GITLAB_TOKEN", "")
	}
}

func Test_parseConfig_multipleHosts(t *testing.T) {
	defer StubConfig(`---
hosts:
  gitlab.example.com:
    username: wrongusername
    token: NOTTHIS
  gitlab.com:
    username: monalisa
    token: OTOKEN
`, "")()
	// prevent using env variable for test
	envToken := os.Getenv("GITLAB_TOKEN")
	if envToken != "" {
		_ = os.Setenv("GITLAB_TOKEN", "")
	}
	config, err := ParseConfig("config.yml")
	eq(t, err, nil)
	username, err := config.Get("gitlab.com", "username")
	eq(t, err, nil)
	eq(t, username, "monalisa")
	token, err := config.Get("gitlab.com", "token")
	eq(t, err, nil)
	eq(t, token, "OTOKEN")
	if envToken != "" {
		_ = os.Setenv("GITLAB_TOKEN", "")
	}
}

func Test_parseConfig_AliasesFile(t *testing.T) {
	defer StubConfig("", `---
gitlab.com:
  username: monalisa
  token: OTOKEN
`)()
	// prevent using env variable for test
	envToken := os.Getenv("GITLAB_TOKEN")
	if envToken != "" {
		_ = os.Setenv("GITLAB_TOKEN", "")
	}
	config, err := ParseConfig("config.yml")
	eq(t, err, nil)
	username, err := config.Get("gitlab.com", "username")
	eq(t, err, nil)
	eq(t, username, "monalisa")
	token, err := config.Get("gitlab.com", "token")
	eq(t, err, nil)
	eq(t, token, "OTOKEN")
	if envToken != "" {
		_ = os.Setenv("GITLAB_TOKEN", "")
	}
}

func Test_parseConfig_hostFallback(t *testing.T) {
	defer StubConfig(`---
git_protocol: ssh
`, `---
gitlab.com:
    username: monalisa
    token: OTOKEN
gitlab.example.com:
    username: wrongusername
    token: NOTTHIS
    git_protocol: https
`)()
	config, err := ParseConfig("config.yml")
	eq(t, err, nil)
	val, err := config.Get("gitlab.example.com", "git_protocol")
	eq(t, err, nil)
	eq(t, val, "https")
	val, err = config.Get("gitlab.com", "git_protocol")
	eq(t, err, nil)
	eq(t, val, "ssh")
	val, err = config.Get("nonexist.io", "git_protocol")
	eq(t, err, nil)
	eq(t, val, "ssh")
}

func Test_ParseConfig_migrateConfig(t *testing.T) {
	defer StubConfig(`---
gitlab.com:
  - username: keiyuri
    token: 123456
`, `ci: pipeline ci
co: mr checkout
`)()

	mainBuf := bytes.Buffer{}
	aliasesBuf := bytes.Buffer{}
	defer StubWriteConfig(&mainBuf, &aliasesBuf)()
	defer StubBackupConfig()()

	_, err := ParseConfig("config.yml")
	assert.Nil(t, err)

	expectedMain := "# What protocol to use when performing git operations. Supported values: ssh, https\ngit_protocol: ssh\n# What editor glab should run when creating issues, merge requests, etc.  This is a global config that cannot be overridden by hostname.\neditor:\n# What browser glab should run when opening links. This is a global config that cannot be overridden by hostname.\nbrowser:\n# Git remote alias which glab should use when fetching the remote url. This can be overridden by hostname\nremote_alias: origin\n# Set your desired markdown renderer style. Available options are [dark, light, notty] or set a custom style. Refer to https://github.com/charmbracelet/glamour#styles\nglamour_style: dark\n# Allow glab to automatically check for updates and notify you when there are new updates\ncheck_update: false\n# configuration specific for gitlab instances\nhosts:\n    gitlab.com:\n        # What protocol to use to access the api endpoint. Supported values: http, https\n        protocol: https\n        # Your GitLab access token. Get an access token at https://gitlab.com/profile/personal_access_tokens\n        token: 123456\n        username: keiyuri\n"
	expectedAliases := `ci: pipeline ci
co: mr checkout
`

	assert.Equal(t, expectedMain, mainBuf.String())
	assert.Equal(t, expectedAliases, aliasesBuf.String())
}

func Test_parseConfigFile(t *testing.T) {
	tests := []struct {
		contents string
		wantsErr bool
	}{
		{
			contents: "",
			wantsErr: true,
		},
		{
			contents: " ",
			wantsErr: false,
		},
		{
			contents: "\n",
			wantsErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("contents: %q", tt.contents), func(t *testing.T) {
			defer StubConfig(tt.contents, "")()
			_, yamlRoot, err := parseConfigFile("config.yml")
			if tt.wantsErr != (err != nil) {
				t.Fatalf("got error: %v", err)
			}
			if tt.wantsErr {
				return
			}
			assert.Equal(t, yaml.MappingNode, yamlRoot.Content[0].Kind)
			assert.Equal(t, 0, len(yamlRoot.Content[0].Content))
		})
	}
}

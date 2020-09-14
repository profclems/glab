package config

import (
	"bytes"
	"fmt"
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
    user: monalisa
    token: OTOKEN
`, "")()
	config, err := ParseConfig("config.yml")
	eq(t, err, nil)
	user, err := config.Get("gitlab.com", "user")
	eq(t, err, nil)
	eq(t, user, "monalisa")
	token, err := config.Get("gitlab.com", "token")
	eq(t, err, nil)
	eq(t, token, "OTOKEN")
}

func Test_parseConfig_multipleHosts(t *testing.T) {
	defer StubConfig(`---
hosts:
  gitlab.example.com:
    user: wronguser
    token: NOTTHIS
  gitlab.com:
    user: monalisa
    token: OTOKEN
`, "")()
	config, err := ParseConfig("config.yml")
	eq(t, err, nil)
	user, err := config.Get("gitlab.com", "user")
	eq(t, err, nil)
	eq(t, user, "monalisa")
	token, err := config.Get("gitlab.com", "token")
	eq(t, err, nil)
	eq(t, token, "OTOKEN")
}

func Test_parseConfig_hostsFile(t *testing.T) {
	defer StubConfig("", `---
gitlab.com:
  user: monalisa
  token: OTOKEN
`)()
	config, err := ParseConfig("config.yml")
	eq(t, err, nil)
	user, err := config.Get("gitlab.com", "user")
	eq(t, err, nil)
	eq(t, user, "monalisa")
	token, err := config.Get("gitlab.com", "token")
	eq(t, err, nil)
	eq(t, token, "OTOKEN")
}

func Test_parseConfig_hostFallback(t *testing.T) {
	defer StubConfig(`---
git_protocol: ssh
`, `---
gitlab.com:
    user: monalisa
    token: OTOKEN
example.com:
    user: wronguser
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
  - user: keiyuri
    token: 123456
`, "")()

	mainBuf := bytes.Buffer{}
	hostsBuf := bytes.Buffer{}
	defer StubWriteConfig(&mainBuf, &hostsBuf)()
	defer StubBackupConfig()()

	_, err := ParseConfig("config.yml")
	assert.Nil(t, err)

	expectedMain := "# What protocol to use when performing git operations. Supported values: ssh, https\ngit_protocol: https\n# What editor gh should run when creating issues, pull requests, etc. If blank, will refer to environment.\neditor:\n# When to interactively prompt. This is a global config that cannot be overriden by hostname. Supported values: enabled, disabled\nprompt: enabled\n# Aliases allow you to create nicknames for gh commands\naliases:\n    co: pr checkout\n"
	expectedHosts := `gitlab.com:
    user: keiyuri
    token: "123456"
`

	assert.Equal(t, expectedMain, mainBuf.String())
	assert.Equal(t, expectedHosts, hostsBuf.String())
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

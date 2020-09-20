package config

import (
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

func Test_parseConfig_Hosts(t *testing.T) {
	defer StubConfig(`---
hosts:
  gitlab.com:
    username: monalisa
    token: OTOKEN
`, `
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

func Test_parseConfig_AliasesFile(t *testing.T) {
	defer StubConfig("", `---
ci: pipeline ci
co: mr checkout
`)()
	config, err := ParseConfig("aliases.yml")
	eq(t, err, nil)
	aliases, err := config.Aliases()
	eq(t, err, nil)
	a, isAlias := aliases.Get("ci")
	eq(t, isAlias, true)
	eq(t, a, "pipeline ci")
	b, isAlias := aliases.Get("co")
	eq(t, isAlias, true)
	eq(t, b, "mr checkout")
	eq(t, len(aliases.All()), 2)
}

func Test_parseConfig_hostFallback(t *testing.T) {
	defer StubConfig(`---
git_protocol: ssh
hosts:
  gitlab.com:
    username: monalisa
    token: OTOKEN
  gitlab.example.com:
    username: wrongusername
    token: NOTTHIS
    git_protocol: https
`, `
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

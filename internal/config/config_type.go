package config

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"errors"
	"gopkg.in/yaml.v3"
)

const (
	defaultGitProtocol 	= "ssh"
	defaultGlamourStyle = "dark"
	defaultHostname    	= "gitlab.com"
)

// This interface describes interacting with some persistent configuration for glab.
type Config interface {
	Get(string, string) (string, error)
	GetWithSource(string, string) (string, string, error)
	Set(string, string, string) error
	UnsetHost(string)
	Hosts() ([]string, error)
	Aliases() (*AliasConfig, error)
	CheckWriteable(string, string) error
	Write() error
}

type NotFoundError struct {
	error
}

type HostConfig struct {
	ConfigMap
	Host string
}

type LocalConfig struct {
	ConfigMap
}

// This type implements a low-level get/set config that is backed by an in-memory tree of Yaml
// nodes. It allows us to interact with a yaml-based config programmatically, preserving any
// comments that were present when the yaml was parsed.
type ConfigMap struct {
	Root *yaml.Node
}

// Default returns the host name of the default GitLab instance
func Default() string {
	return defaultHostname
}

// IsEnterprise reports whether a non-normalized host name looks like a GHE instance
func IsSelfHosted(h string) bool {
	return NormalizeHostname(h) != defaultHostname
}

// NormalizeHostname returns the canonical host name of a GitLab instance
// Taking cover in case GitLab allows subdomains on gitlab.com https://gitlab.com/gitlab-org/gitlab/-/issues/26703
func NormalizeHostname(h string) string {
	hostname := strings.ToLower(h)
	if strings.HasSuffix(hostname, "."+defaultHostname) {
		return defaultHostname
	}
	return hostname
}

func (cm *ConfigMap) Empty() bool {
	return cm.Root == nil || len(cm.Root.Content) == 0
}

func (cm *ConfigMap) GetStringValue(key string) (string, error) {
	entry, err := cm.FindEntry(key)
	if err != nil {
		return "", err
	}
	return entry.ValueNode.Value, nil
}

func (cm *ConfigMap) SetStringValue(key, value string) error {
	entry, err := cm.FindEntry(key)

	var notFound *NotFoundError

	valueNode := entry.ValueNode

	if err != nil && errors.As(err, &notFound) {
		keyNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: key,
		}
		valueNode = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!str",
			Value: "",
		}

		cm.Root.Content = append(cm.Root.Content, keyNode, valueNode)
	} else if err != nil {
		return err
	}

	valueNode.Value = value

	return nil
}

type ConfigEntry struct {
	KeyNode   *yaml.Node
	ValueNode *yaml.Node
	Index     int
}

func (cm *ConfigMap) FindEntry(key string) (ce *ConfigEntry, err error) {
	err = nil

	ce = &ConfigEntry{}

	topLevelKeys := cm.Root.Content
	for i, v := range topLevelKeys {
		if v.Value == key {
			ce.KeyNode = v
			ce.Index = i
			if i+1 < len(topLevelKeys) {
				ce.ValueNode = topLevelKeys[i+1]
			}
			return
		}
	}

	return ce, &NotFoundError{errors.New("not found")}
}

func (cm *ConfigMap) RemoveEntry(key string) {
	var newContent []*yaml.Node

	content := cm.Root.Content
	for i := 0; i < len(content); i++ {
		if content[i].Value == key {
			i++ // skip the next node which is this key's value
		} else {
			newContent = append(newContent, content[i])
		}
	}

	cm.Root.Content = newContent
}

func NewConfig(root *yaml.Node) Config {
	return &fileConfig{
		ConfigMap:    ConfigMap{Root: root.Content[0]},
		documentRoot: root,
	}
}

// NewFromString initializes a Config from a yaml string
func NewFromString(str string) Config {
	root, err := parseConfigData([]byte(str))
	if err != nil {
		panic(err)
	}
	return NewConfig(root)
}

// NewBlankConfig initializes a config file pre-populated with comments and default values
func NewBlankConfig() Config {
	return NewConfig(NewBlankRoot())
}

func NewBlankRoot() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{
						HeadComment: "What protocol to use when performing git operations. Supported values: ssh, https",
						Kind:        yaml.ScalarNode,
						Value:       "git_protocol",
					},
					{
						Kind:  yaml.ScalarNode,
						Value: defaultGitProtocol,
					},
					{
						HeadComment: "What editor glab should run when creating issues, merge requests, etc.  This is a global config that cannot be overridden by hostname.",
						Kind:        yaml.ScalarNode,
						Value:       "editor",
					},
					{
						Kind:  yaml.ScalarNode,
						Value: "",
					},
					{
						HeadComment: "What browser glab should run when opening links. This is a global config that cannot be overridden by hostname.",
						Kind:        yaml.ScalarNode,
						Value:       "browser",
					},
					{
						Kind:  yaml.ScalarNode,
						Value: "",
					},
					{
						HeadComment: "Git remote alias which glab should use when fetching the remote url. This can be overridden by hostname",
						Kind:        yaml.ScalarNode,
						Value:       "remote_alias",
					},
					{
						Kind:  yaml.ScalarNode,
						Value: "origin",
					},
					{
						HeadComment: "Set your desired markdown renderer style. Available options are [dark, light, notty] or set a custom style. Refer to https://github.com/charmbracelet/glamour#styles",
						Kind:        yaml.ScalarNode,
						Value:       "glamour_style",
					},
					{
						Kind:  yaml.ScalarNode,
						Value: defaultGlamourStyle,
					},
					{
						HeadComment: "configuration specific for gitlab instances",
						Kind:        yaml.ScalarNode,
						Value:       "hosts",
					},
					{
						Kind: yaml.MappingNode,
						Content: []*yaml.Node{
							{
								Kind:  yaml.ScalarNode,
								Value: "gitlab.com",
							},
							{
								Kind: yaml.MappingNode,
								Content: []*yaml.Node{
									{
										HeadComment: "What protocol to use to access the api endpoint. Supported values: http, https",
										Kind:  yaml.ScalarNode,
										Value: "protocol",
									},
									{
										Kind:  yaml.ScalarNode,
										Value: "https",
									},
									{
										HeadComment: "Your GitLab access token. Get an access token at https://gitlab.com/profile/personal_access_tokens",
										Kind:  yaml.ScalarNode,
										Value: "token",
									},
									{
										Kind:  yaml.ScalarNode,
										Value: "",
									},
								},
							},
						},
					},

					{
						HeadComment: "Aliases allow you to create nicknames for glab commands. Supports shell-executable aliases that may not be glab commands",
						Kind:        yaml.ScalarNode,
						Value:       "aliases",
					},
					{
						Kind: yaml.MappingNode,
						Content: []*yaml.Node{
							{
								Kind:  yaml.ScalarNode,
								Value: "ci",
							},
							{
								Kind:  yaml.ScalarNode,
								Value: "pipeline ci",
							},
							{
								Kind:  yaml.ScalarNode,
								Value: "co",
							},
							{
								Kind:  yaml.ScalarNode,
								Value: "mr checkout",
							},
						},
					},
				},
			},
		},
	}
}

// This type implements a Config interface and represents a config file on disk.
type fileConfig struct {
	ConfigMap
	documentRoot *yaml.Node
}

func (c *fileConfig) Root() *yaml.Node {
	return c.ConfigMap.Root
}

func (c *fileConfig) Get(hostname, key string) (string, error) {
	var env string
	envEq := EnvKeyEquivalence(key)
	for  _, e := range envEq {
		if val := os.Getenv(e); val != "" {
			env = val
			break
		}
	}
	if env != "" {
		return env, nil
	}
	key = ConfigKeyEquivalence(key)
	val, _, err := c.GetWithSource(hostname, key)
	return val, err
}

func (c *fileConfig) GetWithSource(hostname, key string) (string, string, error) {
	if hostname != "" {
		var notFound *NotFoundError

		hostCfg, err := c.configForHost(hostname)
		if err != nil && !errors.As(err, &notFound) {
			return "", "", err
		}

		var hostValue string
		if hostCfg != nil {
			hostValue, err = hostCfg.GetStringValue(key)
			if err != nil && !errors.As(err, &notFound) {
				return "", "", err
			}
		}

		if hostValue != "" {
			// TODO: avoid hard-coding this
			return hostValue, "~/config/glab-cli/config.yml", nil
		}
	}

	// TODO: avoid hard-coding this
	defaultSource := "~/config/glab-cli/config.yml"

	value, err := c.GetStringValue(key)

	var notFound *NotFoundError

	if err != nil && errors.As(err, &notFound) {
		return defaultFor(key), defaultSource, nil
	} else if err != nil {
		return "", defaultSource, err
	}

	if value == "" {
		return defaultFor(key), defaultSource, nil
	}

	return value, defaultSource, nil
}

func (c *fileConfig) Set(hostname, key, value string) error {
	key = ConfigKeyEquivalence(key)
	if hostname == "" {
		return c.SetStringValue(key, value)
	} else {
		hostCfg, err := c.configForHost(hostname)
		var notFound *NotFoundError
		if errors.As(err, &notFound) {
			hostCfg = c.makeConfigForHost(hostname)
		} else if err != nil {
			return err
		}
		return hostCfg.SetStringValue(key, value)
	}
}

func (c *fileConfig) UnsetHost(hostname string) {
	if hostname == "" {
		return
	}

	hostsEntry, err := c.FindEntry("hosts")
	if err != nil {
		return
	}

	cm := ConfigMap{hostsEntry.ValueNode}
	cm.RemoveEntry(hostname)
}

func (c *fileConfig) CheckWriteable(hostname, key string) error {
	// TODO: check filesystem permissions
	return nil
}

func (c *fileConfig) Write() error {
	mainData := yaml.Node{Kind: yaml.MappingNode}
	aliasesData := yaml.Node{Kind: yaml.MappingNode}

	nodes := c.documentRoot.Content[0].Content
	for i := 0; i < len(nodes)-1; i += 2 {
		if nodes[i].Value == "aliases" {
			aliasesData.Content = append(aliasesData.Content, nodes[i+1].Content...)
		} else {
			mainData.Content = append(mainData.Content, nodes[i], nodes[i+1])
		}
	}

	mainBytes, err := yaml.Marshal(&mainData)
	if err != nil {
		return err
	}

	filename := ConfigFile()
	err = WriteConfigFile(filename, yamlNormalize(mainBytes))
	if err != nil {
		return err
	}

	aliasesBytes, err := yaml.Marshal(&aliasesData)
	if err != nil {
		return err
	}

	return WriteConfigFile(aliasesConfigFile(filename), yamlNormalize(aliasesBytes))
}

func yamlNormalize(b []byte) []byte {
	if bytes.Equal(b, []byte("{}\n")) {
		return []byte{}
	}
	return b
}

func (c *fileConfig) Aliases() (*AliasConfig, error) {
	// The complexity here is for dealing with either a missing or empty aliases key. It's something
	// we'll likely want for other config sections at some point.
	entry, err := c.FindEntry("aliases")
	var nfe *NotFoundError
	notFound := errors.As(err, &nfe)
	if err != nil && !notFound {
		return nil, err
	}

	var toInsert []*yaml.Node

	keyNode := entry.KeyNode
	valueNode := entry.ValueNode

	if keyNode == nil {
		keyNode = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "aliases",
		}
		toInsert = append(toInsert, keyNode)
	}

	if valueNode == nil || valueNode.Kind != yaml.MappingNode {
		valueNode = &yaml.Node{
			Kind:  yaml.MappingNode,
			Value: "",
		}
		toInsert = append(toInsert, valueNode)
	}

	if len(toInsert) > 0 {
		var newContent []*yaml.Node
		if notFound {
			newContent = append(c.Root().Content, keyNode, valueNode)
		} else {
			for i := 0; i < len(c.Root().Content); i++ {
				if i == entry.Index {
					newContent = append(newContent, keyNode, valueNode)
					i++
				} else {
					newContent = append(newContent, c.Root().Content[i])
				}
			}
		}
		c.Root().Content = newContent
	}

	return &AliasConfig{
		Parent:    c,
		ConfigMap: ConfigMap{Root: valueNode},
	}, nil
}

func (c *fileConfig) hostEntries() ([]*HostConfig, error) {
	entry, err := c.FindEntry("hosts")
	if err != nil {
		return nil, fmt.Errorf("could not find hosts config: %w", err)
	}

	hostConfigs, err := c.parseHosts(entry.ValueNode)
	if err != nil {
		return nil, fmt.Errorf("could not parse hosts config: %w", err)
	}

	return hostConfigs, nil
}

// Hosts returns a list of all known hostnames configured in hosts.yml
func (c *fileConfig) Hosts() ([]string, error) {
	entries, err := c.hostEntries()
	if err != nil {
		return nil, err
	}

	var hostnames []string
	for _, entry := range entries {
		hostnames = append(hostnames, entry.Host)
	}

	sort.SliceStable(hostnames, func(i, j int) bool { return hostnames[i] == Default() })

	return hostnames, nil
}

func (c *fileConfig) makeConfigForHost(hostname string) *HostConfig {
	hostRoot := &yaml.Node{Kind: yaml.MappingNode}
	hostCfg := &HostConfig{
		Host:      hostname,
		ConfigMap: ConfigMap{Root: hostRoot},
	}

	var notFound *NotFoundError
	hostsEntry, err := c.FindEntry("hosts")
	if errors.As(err, &notFound) {
		hostsEntry.KeyNode = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "hosts",
		}
		hostsEntry.ValueNode = &yaml.Node{Kind: yaml.MappingNode}
		root := c.Root()
		root.Content = append(root.Content, hostsEntry.KeyNode, hostsEntry.ValueNode)
	} else if err != nil {
		panic(err)
	}

	hostsEntry.ValueNode.Content = append(hostsEntry.ValueNode.Content,
		&yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: hostname,
		}, hostRoot)

	return hostCfg
}

func (c *fileConfig) parseHosts(hostsEntry *yaml.Node) ([]*HostConfig, error) {
	var hostConfigs []*HostConfig

	for i := 0; i < len(hostsEntry.Content)-1; i = i + 2 {
		hostname := hostsEntry.Content[i].Value
		hostRoot := hostsEntry.Content[i+1]
		hostConfig := HostConfig{
			ConfigMap: ConfigMap{Root: hostRoot},
			Host:      hostname,
		}
		hostConfigs = append(hostConfigs, &hostConfig)
	}

	if len(hostConfigs) == 0 {
		return nil, errors.New("could not find any host configurations")
	}

	return hostConfigs, nil
}

func defaultFor(key string) string {
	key = strings.ToLower(key)
	// we only have a set default for one setting right now
	switch key {
	case "gitlab_host", "gitlab_uri":
		return defaultHostname
	case "git_protocol":
		return defaultGitProtocol
	default:
		return ""
	}
}

// ConfigKeyEquivalence returns the equivalent key that's actually used in the config file
func ConfigKeyEquivalence(key string) string {
	key = strings.ToLower(key)
	// we only have a set default for one setting right now
	switch key {
	case "gitlab_host", "gitlab_uri":
		return "host"
	case "gitlab_token", "oauth_token":
		return "token"
	case "git_remote_url_var", "git_remote_alias", "remote_alias", "remote_nickname", "git_remote_nickname":
		return "remote_alias"
	default:
		return key
	}
}

// EnvKeyEquivalence returns the equivalent key that's used for environment variables
func EnvKeyEquivalence(key string) []string {
	key = strings.ToLower(key)
	// we only have a set default for one setting right now
	switch key {
	case "host":
		return []string{"GITLAB_HOST", "GITLAB_URI"}
	case "token":
		return []string{"GITLAB_TOKEN", "OAUTH_TOKEN"}
	case "remote_alias":
		return []string{"GIT_REMOTE_URL_VAR", "GIT_REMOTE_ALIAS", "REMOTE_ALIAS", "REMOTE_NICKNAME", "GIT_REMOTE_NICKNAME"}
	default:
		return []string{strings.ToUpper(key)}
	}
}

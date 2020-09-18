package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

var cachedConfig Config
var configError error

func ConfigDir() string {
	var dir string
	if UseGlobalConfig {
		usrHome := os.Getenv("XDG_CONFIG_HOME")
		if usrHome == "" {
			usrHome = os.Getenv("HOME")
			if usrHome == "" {
				usrHome, _ = homedir.Expand("~/.config")
			} else {
				usrHome = filepath.Join(usrHome, ".config")
			}
		}
		dir = filepath.Join(usrHome, "glab-cli")
	} else {
		dir = ".glab-cli/config"
	}
	return dir
}

func ConfigFile() string {
	return path.Join(ConfigDir(), "config.yml")
}

func Init() (Config, error) {
	if cachedConfig != nil || configError != nil {
		return cachedConfig, configError
	}
	cachedConfig, configError = ParseDefaultConfig()

	if os.IsNotExist(configError) {
		UseGlobalConfigDefaultValue := UseGlobalConfig
		UseGlobalConfig = true
		if err := cachedConfig.WriteAll(); err != nil {
			return nil, err
		}
		UseGlobalConfig = UseGlobalConfigDefaultValue
		configError = nil
	}
	return cachedConfig, configError
}

func ParseDefaultConfig() (Config, error) {
	return ParseConfig(ConfigFile())
}

var ReadConfigFile = func(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, pathError(err)
	}

	return data, nil
}

var WriteConfigFile = func(filename string, data []byte) error {
	err := os.MkdirAll(path.Dir(filename), 0771)
	if err != nil {
		return pathError(err)
	}
	_, err = ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = WriteFile(filename, data, 0600)
	return err
}

var BackupConfigFile = func(filename string) error {
	return os.Rename(filename, filename+".bak")
}

func parseConfigFile(filename string) ([]byte, *yaml.Node, error) {
	data, err := ReadConfigFile(filename)
	if err != nil {
		return nil, nil, err
	}

	root, err := parseConfigData(data)
	if err != nil {
		return nil, nil, err
	}
	return data, root, err
}

func parseConfigData(data []byte) (*yaml.Node, error) {
	var root yaml.Node
	err := yaml.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}

	if len(root.Content) == 0 {
		return &yaml.Node{
			Kind:    yaml.DocumentNode,
			Content: []*yaml.Node{{Kind: yaml.MappingNode}},
		}, nil
	}
	if root.Content[0].Kind != yaml.MappingNode {
		return &root, fmt.Errorf("expected a top level map")
	}
	return &root, nil
}

func isLegacy(root *yaml.Node) bool {
	for _, v := range root.Content[0].Content {
		if v.Value == "gitlab.com" {
			return true
		}
	}

	return false
}

func migrateConfig(filename string) error {
	b, err := ReadConfigFile(filename)
	if err != nil {
		return err
	}

	var hosts map[string][]yaml.Node
	err = yaml.Unmarshal(b, &hosts)
	if err != nil {
		return fmt.Errorf("error decoding legacy format: %w", err)
	}

	cfg := NewBlankConfig()
	for hostname, entries := range hosts {
		if len(entries) < 1 {
			continue
		}
		mapContent := entries[0].Content
		for i := 0; i < len(mapContent)-1; i += 2 {
			if err := cfg.Set(hostname, mapContent[i].Value, mapContent[i+1].Value); err != nil {
				return err
			}
		}
	}

	err = BackupConfigFile(filename)
	if err != nil {
		return fmt.Errorf("failed to back up existing config: %w", err)
	}

	return cfg.Write()
}

func ParseConfig(filename string) (Config, error) {
	_, root, err := parseConfigFile(filename)
	var confError error
	if err != nil {
		if os.IsNotExist(err) {
			root = NewBlankRoot()
			confError = err
		} else {
			return nil, err
		}
	}
	if isLegacy(root) {

		err = migrateConfig(filename)
		if err != nil {
			return nil, fmt.Errorf("error migrating legacy config: %w", err)
		}

		_, root, err = parseConfigFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to reparse migrated config: %w", err)
		}
	} else {
		if _, aliasesRoot, err := parseConfigFile(aliasesConfigFile()); err == nil {
			if len(aliasesRoot.Content[0].Content) > 0 {
				newContent := []*yaml.Node{
					{Value: "aliases"},
					aliasesRoot.Content[0],
				}
				restContent := root.Content[0].Content
				root.Content[0].Content = append(newContent, restContent...)
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	return NewConfig(root), confError
}

func pathError(err error) error {
	var pathError *os.PathError
	if errors.As(err, &pathError) && errors.Is(pathError.Err, syscall.ENOTDIR) {
		if p := findRegularFile(pathError.Path); p != "" {
			return fmt.Errorf("remove or rename regular file `%s` (must be a directory)", p)
		}

	}
	return err
}

func findRegularFile(p string) string {
	for {
		if s, err := os.Stat(p); err == nil && s.Mode().IsRegular() {
			return p
		}
		newPath := path.Dir(p)
		if newPath == p || newPath == "/" || newPath == "." {
			break
		}
		p = newPath
	}
	return ""
}

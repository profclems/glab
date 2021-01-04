package config

import (
	"errors"
	"fmt"
	"path"

	"gopkg.in/yaml.v3"
)

type LocalConfig struct {
	ConfigMap
	Parent Config
}

// LocalConfigDir returns the local config path in map
// which must be joined for complete path
var LocalConfigDir = func() []string {
	return []string{".glab-cli", "config"}
}

// LocalConfigFile returns the config file name with full path
var LocalConfigFile = func() string {
	configFile := append(LocalConfigDir(), "config.yml")
	return path.Join(configFile...)
}

func (a *LocalConfig) Get(key string) (string, bool) {
	key = ConfigKeyEquivalence(key)
	if a.Empty() {
		return "", false
	}
	value, _ := a.GetStringValue(key)

	return value, value != ""
}

func (a *LocalConfig) Set(key, value string) error {
	key = ConfigKeyEquivalence(key)
	err := a.SetStringValue(key, value)
	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	err = a.Write()
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (a *LocalConfig) Delete(key string) error {
	a.RemoveEntry(key)

	return a.Write()
}

func (a *LocalConfig) Write() error {
	// Check if it's a git repository
	if !CheckPathExists(".git") {
		return errors.New("not a git repository")
	}

	localConfigBytes, err := yaml.Marshal(a.ConfigMap.Root)
	if err != nil {
		return err
	}
	err = WriteConfigFile(LocalConfigFile(), yamlNormalize(localConfigBytes))

	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Append local config dir if not already ignored in the .gitignore file
	if !CheckFileHasLine(".gitignore", LocalConfigDir()[0]) {
		err := ReadAndAppend(".gitignore", LocalConfigDir()[0]+"\n")
		if err != nil {
			return fmt.Errorf("failed to write file to .gitignore: %w", err)
		}
	}

	return nil
}

func (a *LocalConfig) All() map[string]string {
	out := map[string]string{}

	if a.Empty() {
		return out
	}

	for i := 0; i < len(a.Root.Content)-1; i += 2 {
		key := a.Root.Content[i].Value
		value := a.Root.Content[i+1].Value
		out[key] = value
	}

	return out
}

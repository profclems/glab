package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"path"
)

func aliasesConfigFile() string {
	return path.Join(globalPathDir, "aliases.yml")
}

type AliasConfig struct {
	ConfigMap
	Parent Config
}

func (a *AliasConfig) Get(alias string) (string, bool) {
	if a.Empty() {
		return "", false
	}
	value, _ := a.GetStringValue(alias)

	return value, value != ""
}

func (a *AliasConfig) Set(alias, expansion string) error {
	err := a.SetStringValue(alias, expansion)
	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	err = a.Write()
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (a *AliasConfig) Delete(alias string) error {
	a.RemoveEntry(alias)

	return a.Write()
}

func (a *AliasConfig) Write() error {
	aliasesBytes, err := yaml.Marshal(a.ConfigMap.Root)
	if err != nil {
		return err
	}
	err = WriteConfigFile(aliasesConfigFile(), yamlNormalize(aliasesBytes))

	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

func (a *AliasConfig) All() map[string]string {
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

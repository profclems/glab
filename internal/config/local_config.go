package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"path"
)

type LocalConfig struct {
	ConfigMap
	Parent Config
}


func localConfigFile() (conf string) {
	useGlobalConfigDefaultValue := UseGlobalConfig
	UseGlobalConfig = false
	conf = path.Join(ConfigFile())
	UseGlobalConfig = useGlobalConfigDefaultValue
	fmt.Println(conf)
	return
}


func (a *LocalConfig) Get(key string) (string, bool) {
	if a.Empty() {
		return "", false
	}
	value, _ := a.GetStringValue(key)

	return value, value != ""
}

func (a *LocalConfig) Set(key, value string) error {
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
	localConfigBytes, err := yaml.Marshal(a.ConfigMap.Root)
	if err != nil {
		return err
	}
	err = WriteConfigFile(localConfigFile(), yamlNormalize(localConfigBytes))

	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
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
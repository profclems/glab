package config

import (
	"fmt"
	"path"
)

func hostsConfigFile(filename string) string {
	return path.Join(path.Dir(filename), "hosts.yml")
}

func (c *fileConfig) configForHost(hostname string) (*HostConfig, error) {
	hosts, err := c.hostEntries()
	if err != nil {
		return nil, fmt.Errorf("failed to parse hosts config: %w", err)
	}

	for _, hc := range hosts {
		if hc.Host == hostname {
			return hc, nil
		}
	}
	return nil, &NotFoundError{fmt.Errorf("could not find config entry for %q", hostname)}
}

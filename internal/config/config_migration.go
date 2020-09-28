package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// MigrateOldConfig migrates the old aliases and configuration
func MigrateOldConfig() error {
	// migrate config directory first
	err := migrateGlobalConfigDir()
	if err != nil {
		return err
	}

	// initialise new configuration
	cfg, err := Init()
	if err != nil {
		return err
	}
	// and migrate the alias file
	// Note that this uses the new config directory so it's important to run migrateGlobalConfigDir() first
	if err := migrateOldAliasFile(cfg); err != nil {
		return err
	}
	// migrate old global configs
	if err := migrateUserConfigs(ConfigDir(), cfg, true); err != nil {
		return err
	}
	// migrate old local configs
	if err := migrateUserConfigs(".glab-cli/config", cfg, false); err != nil {
		return err
	}
	return nil
}

// getAllOldAliases retrieves all of the aliases in the old aliases.format file.
func getAllOldAliases(aliasFile string) map[string]string {
	if !CheckFileExists(aliasFile) {
		return nil
	}

	contents, err := ioutil.ReadFile(aliasFile)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(contents), "\n")
	if len(lines) == 0 {
		return nil
	}

	aliasMap := make(map[string]string)

	for _, line := range lines {
		if line != "" {
			aliasSplit := strings.SplitN(line, ":", 2)
			aliasMap[aliasSplit[0]] = aliasSplit[1]
		}
	}

	return aliasMap
}

func getOldGlobalConfigDir() (string, error) {
	usrHome, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(usrHome, ".glab-cli", "config"), nil
}

func migrateGlobalConfigDir() error {
	// check if xdg directory exists, bail if so.
	newConfigDir := ConfigDir()
	if CheckPathExists(newConfigDir) {
		return nil
	}

	// check if old config dir exists, or there's nothing to migrate.
	oldConfigDir, err := getOldGlobalConfigDir()
	if err != nil {
		return err
	}
	if !CheckPathExists(oldConfigDir) {
		return nil
	}

	// do the migration
	log.Println("- Migrating config dir to XDG_CONFIG_HOME.")
	// First make sure parent directory exists
	if !CheckPathExists(filepath.Join(newConfigDir, "..")) {
		if err := os.MkdirAll(filepath.Join(newConfigDir, ".."), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create new parent config dir: %v", err)
		}
	}

	if err := os.Rename(oldConfigDir, newConfigDir); err != nil {
		return fmt.Errorf("failed to move config dir: %v", err)
	}

	// cleanup: remove parent directory tree of oldConfigDir if empty
	return os.Remove(filepath.Join(oldConfigDir, ".."))
}

// migrateOldAliasFile gets the aliases in the old aliases.format and inserts in the new aliases.yml
// Note that this uses the new config directory so it's important to run migrateGlobalConfigDir() first
func migrateOldAliasFile(cfg Config) error {
	oldAliasFile := filepath.Join(ConfigDir(), "aliases.format")
	if CheckFileExists(oldAliasFile) {
		log.Println("- Migrating aliases")
		// Get aliases in the old alias file
		OldAliases := getAllOldAliases(oldAliasFile)

		// Get the new alias file
		newAliasCfg, err := cfg.Aliases()
		if err != nil {
			return err
		}

		// insert the old aliases in the new alias file
		for name, command := range OldAliases {
			err = newAliasCfg.Set(name, command)
			if err != nil {
				return err
			}
		}

		// backup the old alias file
		return BackupConfigFile(oldAliasFile)
	}
	return nil
}

// migrateUserConfigs gets the config in the old config (.env) and insert into the new config file
// Note that this uses the new config directory so it's important to run migrateGlobalConfigDir() first
func migrateUserConfigs(filePath string, cfg Config, isGlobal bool) error {
	oldConfigFile := filepath.Join(filePath, ".env")
	if CheckFileExists(oldConfigFile) {
		log.Println("- Migrating configuration")
		data, _ := ioutil.ReadFile(oldConfigFile)
		file := string(data)
		temp := strings.Split(file, "\n")

		var host string
		var token string
		var schema string
		var err error

		for _, item := range temp {
			item = strings.TrimSpace(item)
			if item != "" {
				env := strings.SplitN(item, "=", 2)
				if len(env) == 2 {
					// skip if config is hostname or token
					if isGlobal && (env[0] == "GITLAB_URI" || env[0] == "GITLAB_TOKEN") {
						if env[0] == "GITLAB_URI" {
							host = env[1]
						} else if env[0] == "GITLAB_TOKEN" {
							token = env[1]
						}
						continue
					}
					cfg, err = writeConfig(cfg, env[0], env[1], isGlobal)
					if err != nil {
						return err
					}
				}
			}
		}

		if host != "" {
			h, err := url.Parse(host)
			if err == nil {
				host = h.Hostname()
				schema = h.Scheme
			}
			err = cfg.Set(host, "api_protocol", schema)
			if err != nil {
				return err
			}
		}
		if token != "" {
			err = cfg.Set(host, "token", token)
			if err != nil {
				return err
			}
		}

		err = cfg.Write()
		if err != nil {
			return err
		}
		// backup the old alias file
		return BackupConfigFile(oldConfigFile)
	}

	return nil
}

func writeConfig(cfg Config, key, value string, isGlobal bool) (nCfg Config, err error) {
	nCfg = cfg
	if !isGlobal {
		lCfg, _ := cfg.Local()
		err = lCfg.Set(key, value)
	} else {
		err = nCfg.Set("", key, value)
	}
	return
}

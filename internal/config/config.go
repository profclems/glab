package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/profclems/glab/internal/utils"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	// UseGlobalConfig : use the global configuration file
	globalPathDir           = ""
	configFileFileParentDir = ".glab-cli"
	configFile              = configFileFileParentDir + "/config/.env"
	globalConfigFile        = configFile
	aliasFile               = configFileFileParentDir + "/config/aliases.yml"
)

func getXdgGlobalConfigDir() (string, error) {
	cfgDir := os.Getenv("XDG_CONFIG_HOME")
	if cfgDir == "" {
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			homeDir, _ = homedir.Dir()
		}
		if homeDir != "" {
			cfgDir = filepath.Join(homeDir, ".config/")
		}
	}

	if cfgDir == "" {
		return "", fmt.Errorf("could not determine XDG_CONFIG_HOME directory")
	}

	return filepath.Join(cfgDir, "glab-cli"), nil
}

// SetGlobalPathDir sets the directory for the global config file
func SetGlobalPathDir() error {
	err := migrateGlobalConfigDir()
	if err != nil {
		return err
	}

	globalPathDir, err = getXdgGlobalConfigDir()
	if err != nil {
		return err
	}
	if oldCfg, _ := getOldGlobalConfigDir(); !CheckPathExists(globalPathDir) && CheckPathExists(oldCfg) {
		// Migration apparently failed, use old dir.
		globalPathDir = oldCfg
	}

	globalConfigFile = filepath.Join(globalPathDir, filepath.Base(configFile))
	aliasFile = filepath.Join(globalPathDir, "aliases.yml")
	if err := migrateOldAliasFile(); err != nil {
		return err
	}
	return nil
}

// GetAllAliases retrieves all of the aliases in the old aliases.yml file.
func GetAllOldAliases() map[string]string {
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

// PromptAndSetEnv : prompts user for value and writes to config
func Prompt(question, defaultVal string) (envVal string, err error) {
	envVal = utils.AskQuestionWithInput(question, defaultVal, false)
	return
}

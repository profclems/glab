package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/profclems/glab/internal/manip"

	"github.com/google/renameio"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
)

var (
	// UseGlobalConfig : use the global configuration file
	UseGlobalConfig         bool
	globalPathDir           = ""
	configFileFileParentDir = ".glab-cli"
	configFile              = configFileFileParentDir + "/config/.env"
	globalConfigFile        = configFile
	aliasFile               = configFileFileParentDir + "/config/aliases.format"
)

func getOldGlobalConfigDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(usr.HomeDir, ".glab-cli", "config")
}

func getXdgGlobalConfigDir() string {
	cfgDir := os.Getenv("XDG_CONFIG_HOME")
	if cfgDir == "" {
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			if usr, err := user.Current(); err == nil {
				homeDir = usr.HomeDir
			}
		}
		if homeDir != "" {
			cfgDir = filepath.Join(homeDir, ".config/")
		}
	}

	if cfgDir == "" {
		log.Fatal("Could not determine XDG_CONFIG_HOME directory.")
	}

	return filepath.Join(cfgDir, "glab-cli")
}

func migrateGlobalConfigDir() {
	// check if xdg directory exists, bail if so.
	newConfigDir := getXdgGlobalConfigDir()
	if CheckPathExists(newConfigDir) {
		return
	}

	// check if old config dir exists, or there's nothing to migrate.
	oldConfigDir := getOldGlobalConfigDir()
	if !CheckPathExists(oldConfigDir) {
		return
	}

	// do the migration
	log.Println("Migrating config dir to XDG_CONFIG_HOME.")

	// First make sure parent directory exists
	if !CheckPathExists(filepath.Join(newConfigDir, "..")) {
		if err := os.MkdirAll(filepath.Join(newConfigDir, ".."), os.ModePerm); err != nil {
			fmt.Println("Failed to create new parent config dir:", err)
		}
	}

	if err := os.Rename(oldConfigDir, newConfigDir); err != nil {
		fmt.Println("Failed to move config dir:", err)
	}

	// cleanup: remove parent directory tree of oldConfigDir if empty
	_ = os.Remove(filepath.Join(oldConfigDir, ".."))
}

func SetGlobalPathDir() {
	migrateGlobalConfigDir()

	globalPathDir = getXdgGlobalConfigDir()

	if !CheckPathExists(globalPathDir) && CheckPathExists(getOldGlobalConfigDir()) {
		// Migration apparently failed, use old dir.
		globalPathDir = getOldGlobalConfigDir()
	}

	globalConfigFile = filepath.Join(globalPathDir, filepath.Base(configFile))
	aliasFile = filepath.Join(globalPathDir, "aliases.format")
}

// GetEnv : returns env variable value
func GetEnv(key string) string {
	if key != "" {
		env := os.Getenv(key) //first get user defined variable via export

		if env == "" {
			env = GetKeyValueInFile(configFile, key) //Find variable in local env
			if env == "NOTFOUND" || env == "OK" {
				env = GetKeyValueInFile(globalConfigFile, key) //Find variable in global env
				if env == "NOTFOUND" || env == "OK" {
					//log.Fatal("Configuration not set for ", key)
					return ""
				}
			}
		}
		return env
	}
	return ""
}

// SetEnv : sets env variable
func SetEnv(key, value string) {
	cFile := configFile
	if UseGlobalConfig {
		cFile = globalConfigFile
	}

	defer InvalidateEnvCacheForFile(cFile)

	data, err := ioutil.ReadFile(cFile)
	if err != nil && !os.IsNotExist(err) {
		log.Println("Failed to read/update env config:", err)
		return
	}

	file := string(data)
	temp := strings.Split(file, "\n")
	newData := ""
	keyExists := false
	newConfig := key + "=" + (value) + "\n"
	for _, item := range temp {
		if item == "" {
			continue
		}

		env := strings.Split(item, "=")
		if env[0] == key {
			newData += newConfig
			keyExists = true
		} else {
			newData += item + "\n"
		}
	}
	if !keyExists {
		newData += newConfig
	}
	_ = os.MkdirAll(filepath.Join(cFile, ".."), 0755)
	if err = renameio.WriteFile(cFile, []byte(newData), 0666); err != nil {
		log.Println("Failed to update config file:", err)
		return
	}

	if !UseGlobalConfig && !CheckFileHasLine(".gitignore", configFileFileParentDir) {
		ReadAndAppend(".gitignore", configFileFileParentDir+"\n")
	}
}

// SetAlias sets an alias for a command
func SetAlias(name string, command string) error {
	if !CheckFileExists(aliasFile) {
		aliasDir := filepath.Join(aliasFile, "..")
		if !CheckPathExists(aliasDir) {
			errDir := os.MkdirAll(aliasDir, 0700)
			if errDir != nil {
				return errDir
			}
		}
		f, err := os.Create(aliasFile)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	contents, err := ioutil.ReadFile(aliasFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(contents), "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = []string{}
	}
	set := false

	for i, line := range lines {
		aliasSplit := strings.SplitN(line, ":", 2)
		if aliasSplit[0] == name {
			lines[i] = name + ":" + command
			set = true
			break
		}
	}

	if !set {
		lines = append(lines, name+":"+command)
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(aliasFile, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetAllAliases retrieves all of the aliases.
func GetAllAliases() map[string]string {
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

// GetAlias retrieves the command for an alias
func GetAlias(name string) string {
	if !CheckFileExists(aliasFile) {
		return ""
	}

	contents, err := ioutil.ReadFile(aliasFile)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		aliasSplit := strings.SplitN(line, ":", 2)
		if aliasSplit[0] == name {
			return aliasSplit[1]
		}
	}

	return ""
}

// DeleteAlias deletes an alias
func DeleteAlias(name string) error {
	if !CheckFileExists(aliasFile) {
		return errors.New("No aliases are currently set")
	}

	contents, err := ioutil.ReadFile(aliasFile)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(contents), "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = []string{}
	}
	deleted := false

	for i, line := range lines {
		aliasSplit := strings.SplitN(line, ":", 2)
		if aliasSplit[0] == name {
			copy(lines[i:], lines[i+1:])
			lines[len(lines)-1] = ""
			lines = lines[:len(lines)-1]
			deleted = true
			break
		}
	}

	if !deleted {
		return errors.New("That alias does not exist")
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(aliasFile, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func GetRepo() string {
	gitRemoteVar, err := gitconfig.Entire("remote." + GetEnv("GIT_REMOTE_URL_VAR") + ".url")
	if err != nil {
		log.Fatal("Could not find remote url for gitlab. Run git config init")
	}
	repoBaseUrl := strings.Trim(GetEnv("GITLAB_URI"), "/ ")
	repoBaseUrl = strings.TrimPrefix(repoBaseUrl, "https://")
	repoBaseUrl = strings.TrimPrefix(repoBaseUrl, "http://")
	repo := strings.TrimSuffix(gitRemoteVar, ".git")
	repo = strings.TrimPrefix(repo, repoBaseUrl)
	repo = strings.TrimPrefix(repo, "https://"+repoBaseUrl)
	repo = strings.TrimPrefix(repo, "http://"+repoBaseUrl)
	repo = strings.TrimPrefix(repo, "git@"+repoBaseUrl+":")
	return strings.Trim(repo, "/")
}

func readAndSetEnv(question, env string) string {
	envDefVal := GetEnv(env)
	envVal := manip.AskQuestionWithInput(question, envDefVal, false)
	SetEnv(env, envVal)
	return envVal
}

func Set(cmd *cobra.Command, args []string) {
	var isUpdated bool
	if b, _ := cmd.Flags().GetBool("global"); b {
		UseGlobalConfig = true
	}
	if b, _ := cmd.Flags().GetString("token"); b != "" {
		SetEnv("GITLAB_TOKEN", b)
		isUpdated = true
	}
	if b, _ := cmd.Flags().GetString("url"); b != "" {
		SetEnv("GITLAB_URI", b)
		isUpdated = true
	}
	if b, _ := cmd.Flags().GetString("remote-var"); b != "" {
		SetEnv("GIT_REMOTE_URL_VAR", b)
		isUpdated = true
	}
	if b, _ := cmd.Flags().GetString("pid"); b != "" {
		SetEnv("GITLAB_PROJECT_ID", b)
		isUpdated = true
	}
	if !isUpdated {
		readAndSetEnv(fmt.Sprintf("Enter default Gitlab Host (Current Value: %s): ", GetEnv("GITLAB_URI")), "GITLAB_URI")
		readAndSetEnv("Enter default Gitlab Token: ", "GITLAB_TOKEN")
		readAndSetEnv(fmt.Sprintf("Enter Git remote url variable (Current Value: %s): ", GetEnv("GIT_REMOTE_URL_VAR")), "GIT_REMOTE_URL_VAR")
		isUpdated = true
	}
	if isUpdated {
		fmt.Println(aurora.Green("Environment variable(s) updated"))
	}
}

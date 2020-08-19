package config

import (
	"bufio"
	"errors"
	"fmt"
	"glab/internal/manip"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
)

var (
	// UseGlobalConfig : use the global configuration file
	UseGlobalConfig         bool
	globalPathDir           = ""
	configFileFileParentDir = ".glab-cli"
	configFileFileDir       = configFileFileParentDir + "/config"
	configFile              = configFileFileDir + "/.env"
	globalConfigFile        = configFile
	aliasFile               = configFileFileDir + "/aliases.txt"
)

func SetGlobalPathDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	globalPathDir = usr.HomeDir
	globalConfigFile = globalPathDir + "/" + globalConfigFile
	return globalPathDir
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
	cFileFileParentDir := configFileFileParentDir
	cFileDir := configFileFileDir
	if UseGlobalConfig {
		cFileFileParentDir = globalPathDir + "/" + cFileFileParentDir
		cFileDir = globalPathDir + "/" + configFileFileDir
		cFile = globalConfigFile
	}
	data, _ := ioutil.ReadFile(cFile)

	file := string(data)
	line := 0
	temp := strings.Split(file, "\n")
	newData := ""
	keyExists := false
	newConfig := key + "=" + (value) + "\n"
	for _, item := range temp {
		//fmt.Println("[",line,"]",item)
		env := strings.Split(item, "=")
		justString := fmt.Sprint(item)
		if env[0] == key {
			newData += newConfig
			keyExists = true
		} else {
			newData += justString + "\n"
		}
		line++
	}
	if !keyExists {
		newData += newConfig
	}
	_ = os.Mkdir(cFileFileParentDir, 0700)
	_ = os.Mkdir(cFileDir, 0700)
	f, _ := os.Create(cFile) // Create a writer
	w := bufio.NewWriter(f)
	_, _ = w.WriteString(strings.Trim(newData, "\n"))
	_ = w.Flush()
	if GetKeyValueInFile(".gitignore", configFileFileParentDir) == "NOTFOUND" {
		ReadAndAppend(".gitignore", configFileFileParentDir)
	}
}

// SetAlias sets an alias for a command
func SetAlias(name string, command string) {
	if !CheckFileExists(aliasFile) {
		_, err := os.Stat(configFileFileDir)
		if os.IsNotExist(err) {
			errDir := os.MkdirAll(configFileFileDir, 0700)
			if errDir != nil {
				log.Fatalln(err)
			}
		}
		f, err := os.Create(aliasFile)
		if err != nil {
			log.Fatalln(err)
		}

		err = f.Close()
		if err != nil {
			log.Println(err)
		}
	}

	contents, err := ioutil.ReadFile(aliasFile)
	if err != nil {
		log.Fatalln(err)
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
		log.Fatalln(err)
	}
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

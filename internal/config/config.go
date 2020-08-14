package config

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
	"glab/internal/manip"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

var (
	// UseGlobalConfig : use the global configuration file
	UseGlobalConfig         bool
	globalPathDir           = ""
	configFileFileParentDir = ".glab-cli"
	configFileFileDir       = configFileFileParentDir + "/config"
	configFile              = configFileFileDir + "/.env"
	globalConfigFile        = configFile
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

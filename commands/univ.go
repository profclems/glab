package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/tcnksm/go-gitconfig"
	"github.com/xanzy/go-gitlab"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// UseGlobalConfig : use the global configuration file
	UseGlobalConfig         bool
	globalPathDir        	= ""
	configFileFileParentDir = ".glab-cli"
	configFileFileDir       = configFileFileParentDir + "/config"
	configFile              = configFileFileDir + "/.env"
	globalConfigFile        = configFile
)

func SetGlobalPathDir() string  {
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	globalPathDir = usr.HomeDir
	globalConfigFile  = globalPathDir + "/" + globalConfigFile
	return globalPathDir
}

func getRepo() string {
	gitlab, err := gitconfig.Entire("remote."+GetEnv("GIT_REMOTE_URL_VAR")+".url")
	if err != nil {
		log.Fatal("Could not find remote url for gitlab")
	}
	repoBaseUrl := strings.Trim(GetEnv("GITLAB_URI"), "/ ")
	repoBaseUrl = strings.TrimPrefix(repoBaseUrl, "https://")
	repoBaseUrl = strings.TrimPrefix(repoBaseUrl, "http://")
	repo :=  strings.TrimSuffix(gitlab, ".git")
	repo = strings.TrimPrefix(repo, repoBaseUrl)
	repo = strings.TrimPrefix(repo, "https://"+repoBaseUrl)
	repo = strings.TrimPrefix(repo, "http://"+repoBaseUrl)
	repo = strings.TrimPrefix(repo, "git@"+repoBaseUrl+":")
	return strings.Trim(repo, "/")
}

// InitGitlabClient : creates client
func InitGitlabClient() (*gitlab.Client, string)  {
	git, err := gitlab.NewClient(GetEnv("GITLAB_TOKEN"), gitlab.WithBaseURL(strings.TrimRight(GetEnv("GITLAB_URI"),"/") + "/api/v4"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return git, getRepo()
}

func VariableExists(key string) string  {
	return GetKeyValueInFile(configFile, key)
}

// GetEnv : returns env variable value
func GetEnv(key string) string {
	env := os.Getenv(key)

	if len(env) == 0 {
		env = GetKeyValueInFile(configFile, key)
		if env == "NOTFOUND" || env == "OK" {
			env = GetKeyValueInFile(globalConfigFile, key)
			if env == "NOTFOUND" || env == "OK" {
				log.Fatal("Configuration not set for ", key)
			}
		}
	}
	return env
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

// ReadAndAppend : appends string to file
func ReadAndAppend(file, text string) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("\n" + text)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// ReplaceNonAlphaNumericChars : Replaces non alpha-numeric values with provided char/string
func ReplaceNonAlphaNumericChars(words, replaceWith string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	newStr := reg.ReplaceAllString(words, replaceWith)
	return newStr
}

// GetKeyValueInFile : returns env variable value
func GetKeyValueInFile(filePath, key string) string {
	data, _ := ioutil.ReadFile(filePath)

	file := string(data)
	line := 0
	temp := strings.Split(file, "\n")
	for _, item := range temp {
		//fmt.Println("[",line,"]",item)
		env := strings.Split(item, "=")
		if env[0] == key {
			if len(env) > 1 {
				return env[1]
			}
			return "OK"
		}
		line++
	}
	return "NOTFOUND"
}

// CommandExists : checks if string is available in the defined commands
func CommandExists(mapArr map[string]func(map[string]string, map[int]string), key string) bool {
	if _, ok := mapArr[key]; ok {
		return true
	}
	return false
}

// CommandArgExists : checks if string is available in the defined command flags
func CommandArgExists(mapArr map[string]string, key string) bool {
	if _, ok := mapArr[key]; ok {
		return true
	}
	return false
}

func stringToInt(str string) int {
	strInt, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return strInt
}

// TimeAgo is ...
func TimeAgo(timeVal time.Time) string {
	//now := time.Now().Format(time.RFC3339)
	layout := "2006-01-02T15:04:05.000Z"
	then, _ := time.Parse(layout, timeVal.Format("2006-01-02T15:04:05.000Z"))
	totalSeconds := time.Since(then).Seconds()
	if totalSeconds < 60 {
		if totalSeconds < 1 {
			totalSeconds = 0
		}
		return fmt.Sprint(totalSeconds, "secs ago")
	} else if totalSeconds >= 60 && totalSeconds < (60*60) {
		return fmt.Sprint(math.Round(totalSeconds/60), "mins ago")
	} else if totalSeconds >= (60*60) && totalSeconds < (60*3600) {
		return fmt.Sprint(math.Round(totalSeconds/(60*60)), "hrs ago")
	} else if totalSeconds >= (60*3600) && totalSeconds < (60*60*3600) {
		return fmt.Sprint(math.Round(totalSeconds/(60*3600)), "days ago")
	}
	return ""
}

// MakeRequest is ...
func MakeRequest(payload, url, method string) map[string]interface{} {

	url = GetEnv("GITLAB_URI") + "/api/v4/" + url
	var reader io.Reader
	if payload != "" && payload != "{}" {
		reader = bytes.NewReader([]byte(payload))
	}

	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	client := &http.Client{}
	request.Header.Set("PRIVATE-TOKEN", GetEnv("GITLAB_TOKEN"))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	m := make(map[string]interface{})
	m["responseCode"] = resp.StatusCode
	m["responseMessage"] = bodyString

	return m
}

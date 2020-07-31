package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	UseGlobalConfig bool
	GlobalPathDir,_ 		= filepath.Abs(filepath.Dir(os.Args[0]))
	ConfigFileFileParentDir = ".glab-cli"
	ConfigFileFileDir       = ConfigFileFileParentDir+"/config"
	ConfigFile              = ConfigFileFileDir +"/.env"
	GlobalConfigFile        = GlobalPathDir+"/"+ConfigFileFileDir +"/.env"
)

func GetEnv(key string) string {
	env := os.Getenv(key)
	// load .env file

	if len(env) == 0 {
		env = GetKeyValueInFile(ConfigFile, key)
		if env == "NOTFOUND" || env == "OK" {
			env = GetKeyValueInFile(GlobalConfigFile, key)
			if env == "NOTFOUND" || env == "OK" {
				log.Fatal("Configuration not set for ", key)
			}
		}
	}
	return  env
}

func SetEnv(key, value string) {
	cFile := ConfigFile
	cFileFileParentDir := ConfigFileFileParentDir
	cFileDir := ConfigFileFileDir
	if UseGlobalConfig {
		cFileFileParentDir = GlobalPathDir + "/" + cFileFileParentDir
		cFileDir = GlobalPathDir + "/" + ConfigFileFileDir
		cFile = GlobalConfigFile
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
	if GetKeyValueInFile(".gitignore",ConfigFileFileParentDir) == "NOTFOUND" {
		ReadAndAppend(".gitignore", ConfigFileFileParentDir)
	}
}

func ReadAndAppend(file, text string) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("\n"+text)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func ReplaceNonAlphaNumericChars(words, replaceWith string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	newStr := reg.ReplaceAllString(words, replaceWith)
	return newStr
}

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
			} else {
				return "OK"
			}
		}
		line++
	}
	return "NOTFOUND"
}

func CommandExists(mapArr map[string]func(map[string]string, map[int]string), key string) bool {
	if _, ok := mapArr[key]; ok {
		return true
	} else {
		return false
	}
}

func CommandArgExists(mapArr map[string]string, key string) bool {
	if _, ok := mapArr[key]; ok {
		return true
	} else {
		return false
	}
}

func TimeAgo(timeVal interface{}) string {
	//now := time.Now().Format(time.RFC3339)
	layout := "2006-01-02T15:04:05.000Z"
	then, _ := time.Parse(layout, timeVal.(string))
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

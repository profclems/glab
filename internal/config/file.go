package config

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// glab environment cache: <file: <key: value>>
var envCache map[string]map[string]string

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

func readConfig(filePath string) map[string]string {
	var config = make(map[string]string)
	data, _ := ioutil.ReadFile(filePath)
	file := string(data)
	temp := strings.Split(file, "\n")
	for _, item := range temp {
		//fmt.Println("[",line,"]",item)
		env := strings.Split(item, "=")
		if len(env) > 1 {
			config[env[0]] = env[1]
		}
	}
	return config
}

// GetKeyValueInFile : returns env variable value
func GetKeyValueInFile(filePath, key string) string {
	configCache, okConfig := envCache[filePath]
	if !okConfig {
		configCache = readConfig(filePath)
		if envCache == nil {
			envCache = make(map[string]map[string]string)
		}
		envCache[filePath] = configCache
	}

	if cachedEnv, okEnv := configCache[key]; okEnv {
		if cachedEnv == "" {
			cachedEnv = "OK"
		}
		return cachedEnv
	}
	return "NOTFOUND"
}

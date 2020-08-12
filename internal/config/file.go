package config

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

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
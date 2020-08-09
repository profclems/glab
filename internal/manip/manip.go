package manip

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"io/ioutil"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

/*
func VariableExists(key string, global bool) string {
	return GetKeyValueInFile(config.ConfigFile, key)
}
 */

func AskQuestionWithSelect(question, defaultVal string, isRequired bool)  {
	color := ""
	prompt := &survey.Select{
		Message: "Choose a color:",
		Options: []string{"red", "blue", "green"},
	}
	survey.AskOne(prompt, &color)
}

func AskQuestionWithInput(question, defaultVal string, isRequired bool) string  {
	str := ""
	prompt := &survey.Input{
		Message: question,
	}
	if isRequired {
			_ = survey.AskOne(prompt, &str, survey.WithValidator(survey.Required))
	} else {
		_ = survey.AskOne(prompt, &str)
	}
	str = strings.TrimSuffix(str, "\n")
	if str == "" && defaultVal != "" {
		return defaultVal
	}
	return str
}

func AskQuestionMultiline(question string, defaultVal string) string  {
	str := ""
	prompt := &survey.Multiline{
		Message: question,
	}
	_ = survey.AskOne(prompt, &str)
	str = strings.TrimSuffix(str, "\n")
	if str == "" && defaultVal != "" {
		return defaultVal
	}
	return str
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
	newStr := reg.ReplaceAllString(strings.Trim(words," "), replaceWith)
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

func StringToInt(str string) int {
	strInt, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return strInt
}

// TruncateString truncate a string by the specified length (n)
func TruncateStrings(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for !utf8.ValidString(s[:n]) {
		n--
	}
	return s[:n]+"..."
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
		return fmt.Sprint(math.Round(totalSeconds), "secs ago")
	} else if totalSeconds >= 60 && totalSeconds < (60*60) {
		return fmt.Sprint(math.Round(totalSeconds/60), "mins ago")
	} else if totalSeconds >= (60*60) && totalSeconds < (60*3600) {
		return fmt.Sprint(math.Round(totalSeconds/(60*60)), "hrs ago")
	} else if totalSeconds >= (60*3600) && totalSeconds < (60*60*3600) {
		return fmt.Sprint(math.Round(totalSeconds/(60*3600)), "days ago")
	}
	return ""
}

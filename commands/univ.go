package commands

import (
	"bytes"
	_ "bytes"
	"fmt"
	_ "fmt"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	_ "io/ioutil"
	"log"
	"math"
	"net/http"
	_ "net/http"
	"os"
	"regexp"
	"time"
)

func GetEnv(key string) string {

	// load .env file
	err := godotenv.Load("./config/.env")

	if err != nil {
		log.Fatalf("Error loading .env file.\nRun glab config init to set up your environment")
	}

	return os.Getenv(key)
}

func SetEnv(key, value string) string {

	// load .env file
	env, _ := godotenv.Unmarshal(key + "=" + value)
	err := godotenv.Write(env, "./.env")

	if err != nil {
		log.Fatalf("Error writing .env file")
	}

	return value
}

func ReplaceNonAlphaNumericChars(words, replaceWith string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	newStr := reg.ReplaceAllString(words, replaceWith)
	return newStr
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

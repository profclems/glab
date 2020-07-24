package commands

import (
	"bytes"
	_ "bytes"
	"encoding/json"
	"fmt"
	_ "fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	_ "io/ioutil"
	"log"
	"net/http"
	_ "net/http"
	"os"
)

func GetEnv(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func SetEnv(key string, value string) string {

	// load .env file
	env, err := godotenv.Unmarshal(key+"="+value)
	err = godotenv.Write(env, "./.env")

	if err != nil {
		log.Fatalf("Error writing .env file")
	}

	return value
}

func CommandExists(mapArr map[string]func(map[string]string), key string) bool {
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

func MakeRequest(payload string, url string, method string) []string {

	url = GetEnv("GITLAB_URI")+"/api/"+GetEnv("API_VERSION")+"/"+url
	reader := bytes.NewReader([]byte(payload))
	request, err := http.NewRequest(method, url, reader)
	if err != nil{
		log.Fatal("Error: ", err)
	}
	client := &http.Client{}
	request.Header.Set("PRIVATE-TOKEN", GetEnv("GITLAB_TOKEN"))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(request)
	if err != nil{
		log.Fatal("Error: ", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		fmt.Println()
	} else {
		fmt.Println("An error occurred connecting to remote host")
		fmt.Print(resp.StatusCode, ": ", bodyString)
	}
	var arr []string
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(bodyString), &m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)
	return arr
}


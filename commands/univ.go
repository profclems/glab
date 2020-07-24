package commands

import (
	_ "bytes"
	"encoding/json"
	"fmt"
	_ "fmt"
	"github.com/joho/godotenv"
	_ "io/ioutil"
	"log"
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
	/*
	url = GetEnv("GITLAB_URI")+"/api/"+GetEnv("API_VERSION")+"/"+url
	reader := bytes.NewReader([]byte(payload))
	request, err := http.NewRequest(method, url, reader)
	if err != nil{
		log.Fatal("Error: ", err)
	}
	client := &http.Client{}
	request.Header.Set("PRIVATE-TOKEN", GetEnv("GITLAB_TOKEN"))
	request.Header.Set("Content-Type", "application/json")
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
		fmt.Println(bodyString)
	} else {
		fmt.Println("An error occurred connecting to remote host")
		fmt.Print(resp.StatusCode, ": ", bodyString)
	}
	 */
	testData := `{"id":68960731,"iid":3,"project_id":20131402,"title":"thdf","description":"djhkjf","state":"opened","created_at":"2020-07-24T19:57:12.481Z","updated_at":"2020-07-24T19:57:12.481Z","closed_at":null,"closed_by":null,"labels":[],"milestone":null,"assignees":[],"author":{"id":5568402,"name":"Clement Sam","username":"profclems","state":"active","avatar_url":"https://assets.gitlab-static.net/uploads/-/system/user/avatar/5568402/avatar.png","web_url":"https://gitlab.com/profclems"},"assignee":null,"user_notes_count":0,"merge_requests_count":0,"upvotes":0,"downvotes":0,"due_date":null,"confidential":false,"discussion_locked":null,"web_url":"https://gitlab.com/profclems/glab/-/issues/3","time_stats":{"time_estimate":0,"total_time_spent":0,"human_time_estimate":null,"human_total_time_spent":null},"task_completion_status":{"count":0,"completed_count":0},"weight":null,"blocking_issues_count":null,"has_tasks":false,"_links":{"self":"https://gitlab.com/api/v4/projects/20131402/issues/3","notes":"https://gitlab.com/api/v4/projects/20131402/issues/3/notes","award_emoji":"https://gitlab.com/api/v4/projects/20131402/issues/3/award_emoji","project":"https://gitlab.com/api/v4/projects/20131402"},"references":{"short":"#3","relative":"#3","full":"profclems/glab#3"},"subscribed":true,"moved_to_id":null}`
	var arr []string
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(testData), &m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)
	return arr
}


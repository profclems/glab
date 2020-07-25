package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	. "github.com/logrusorgru/aurora"
	"log"
	"net/url"
	"os"
	"strings"
)

type IssueInfo struct {
	Title   string  `json:"title"`
	Name string `json:"description"`
	IssueId string `json:"iid"`
	State string `json:"state"`
}

func DisplayIssue(hm map[string]interface{})  {
	duration := TimeAgo(hm["created_at"])
	if hm["state"] == "opened" {
		fmt.Println(Green(fmt.Sprint("#",hm["iid"])), hm["title"], Magenta(duration))
	} else {
		fmt.Println(Red(fmt.Sprint("#",hm["iid"])), hm["title"], Magenta(duration))
	}
}

func CreateIssue(map[string]string)  {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(Cyan("Title"+"\n"+"-> "))
	issueTitle, _ := reader.ReadString('\n')
	issueTitle = strings.Replace(issueTitle, "\n", "", -1)
	fmt.Println()
	fmt.Println(Cyan("Enter Issue Description"), Yellow("[info: Type `exit` to close]"))
	var issueDescription string
	for {
		fmt.Print(Cyan("-> "))
		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)
		if strings.Compare("exit", input) == 0 {
			break
		}
		issueDescription += "\n"+input

	}
	params := url.Values{}
	params.Add("title", issueTitle)
	params.Add("description", issueDescription)
	reqBody := params.Encode()
	fmt.Println(Yellow("Creating Issue {"+issueTitle+"}..."))
	resp := MakeRequest(reqBody,"projects/"+GetEnv("GITLAB_PROJECT_ID")+"/issues","POST")

	if resp["responseCode"]==201 {
		bodyString := resp["responseMessage"]

		fmt.Println(Green("Issue created successfully"))
		if _, ok := bodyString.(string); ok {
			/* act on str */
			m := make(map[string]interface{})
			err := json.Unmarshal([]byte(bodyString.(string)), &m)
			if err != nil {
				log.Fatal(err)
			}
			DisplayIssue(m)
		} else {
			/* not string */
		}
	}
}

func ListIssues(cmdArgs map[string]string)  {
	var queryStrings = "?state="
	if CommandArgExists(cmdArgs, "all") {
		queryStrings = ""
	} else if CommandArgExists(cmdArgs, "closed") {
		queryStrings += "closed"
	} else {
		queryStrings += "opened"
	}
	resp := MakeRequest("{}","projects/"+GetEnv("GITLAB_PROJECT_ID")+"/issues"+queryStrings,"GET")
	//fmt.Println(resp)
	if resp["responseCode"]==200 {
		bodyString := resp["responseMessage"]
		if _, ok := bodyString.(string); ok {
			/* act on str */
			var m []interface{}
			err := json.Unmarshal([]byte(bodyString.(string)), &m)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println()
			for i:=0; i < len(m); i++ {
				hm := m[i].(map[string]interface{})
				DisplayIssue(hm)
			}
		} else {
			/* not string */
		}
	}
}

func ExecIssue(cmdArgs map[string]string)  {
	commandList := map[string]func(map[string]string) {
		"create" : CreateIssue,
		"list" : ListIssues,
	}
	if CommandArgExists(cmdArgs, "create") {
		commandList["create"](cmdArgs)
	} else {
		commandList["list"](cmdArgs)
	}
}

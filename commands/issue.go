package commands

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type IssueInfo struct {
	Id   int64  `json:"title"`
	Name string `json:"description"`
}

func CreateIssue(map[string]string)  {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Issue Details")
	fmt.Println("--------------------------------")
	fmt.Print("Title"+"\n"+"-> ")
	issueTitle, _ := reader.ReadString('\n')
	issueTitle = strings.Replace(issueTitle, "\n", "", -1)
	fmt.Println()
	fmt.Println("Enter Issue Description")
	fmt.Println("info: Type `exit` to close")
	var issueDescription string
	for {
		fmt.Print("-> ")
		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)

		if strings.Compare("exit", input) == 0 {
			break
		}
		issueDescription += "\n"+input

	}
	reqBody := fmt.Sprintf("{\"title\":\"%s\",\"description\":\"%s\"}",issueTitle, interface{}(url.ParseRequestURI(issueDescription)))
	fmt.Println(reqBody)
	MakeRequest(reqBody,"projects/"+GetEnv("GITLAB_PROJECT_ID")+"/issues","POST")
}

func ExecIssue(cmdArgs map[string]string)  {
	commandList := map[string]func(map[string]string) {
		"create" : CreateIssue,
	}
	commandList["create"](cmdArgs)
	/*
	urls := map[string]string {
		"contributions"
	}
	*/

	//MakeRequest(`{}`,"projects/20131402/issues/1","GET")
}

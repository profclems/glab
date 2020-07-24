package commands

import (
	"bufio"
	"fmt"
	. "github.com/logrusorgru/aurora"
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

	fmt.Print(Cyan("Title"+"\n"+"-> "))
	issueTitle, _ := reader.ReadString('\n')
	issueTitle = strings.Replace(issueTitle, "\n", "", -1)
	fmt.Println()
	fmt.Println(Cyan("Enter Issue Description"), Yellow("[info: Type `exit` to close]"))
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
	params := url.Values{}
	params.Add("title", issueTitle)
	params.Add("description", issueDescription)
	reqBody := params.Encode()
	fmt.Println(Green("Creating Issue {"+issueTitle+"}..."))
	MakeRequest(reqBody,"projects/"+GetEnv("GITLAB_PROJECT_ID")+"/issues","POST")
}

func ExecIssue(cmdArgs map[string]string)  {
	commandList := map[string]func(map[string]string) {
		"create" : CreateIssue,
	}
	if CommandArgExists(cmdArgs, "create") {
		commandList["create"](cmdArgs)
	}
}

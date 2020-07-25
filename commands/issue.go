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

func DisplayIssue(hm map[string]interface{})  {
	duration := TimeAgo(hm["created_at"])
	if hm["state"] == "opened" {
		fmt.Println(Green(fmt.Sprint("#",hm["iid"])), hm["title"], Magenta(duration))
	} else {
		fmt.Println(Red(fmt.Sprint("#",hm["iid"])), hm["title"], Magenta(duration))
	}
}

func CreateIssue(cmdArgs map[string]string, _ map[int]string)  {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(Cyan("Title"+"\n"+"-> "))
	issueTitle, _ := reader.ReadString('\n')
	var issueLabel string
	if !CommandArgExists(cmdArgs, "labels") || !CommandArgExists(cmdArgs, "label") {
		fmt.Print(Cyan("Label(s) [Comma Separated]"+"\n"+"-> "))
		issueLabel, _ = reader.ReadString('\n')
	} else {
		issueLabel = strings.Trim(cmdArgs["label"],"[] ")
	}
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
	fmt.Print(Cyan("Due Date"))
	fmt.Print(Yellow("(Format: YYYY-MM-DD)"+"\n"+"-> "))
	issueDue, _ := reader.ReadString('\n')
	params := url.Values{}
	params.Add("title", issueTitle)
	params.Add("description", issueDescription)
	params.Add("labels", issueLabel)
	params.Add("due_date", issueDue)
	if CommandArgExists(cmdArgs, "confidential") {
		params.Add("confidential", "true")
	}
	if CommandArgExists(cmdArgs, "weight") {
		params.Add("weight", cmdArgs["weight"])
	}
	if CommandArgExists(cmdArgs, "mr") {
		params.Add("merge_request_to_resolve_discussions_of", cmdArgs["mr"])
	}
	if CommandArgExists(cmdArgs, "milestone") {
		params.Add("milestone_id", cmdArgs["milestone"])
	}
	if CommandArgExists(cmdArgs, "epic") {
		params.Add("epic_id", cmdArgs["epic"])
	}
	if CommandArgExists(cmdArgs, "assigns") {
		params.Add("epic_id", cmdArgs["epic"])
		assignId := cmdArgs["assigns"]
		arrIds := strings.Split(strings.Trim(assignId,"[] "), ",")
		for _, i2 := range arrIds {
			params.Add("assignee_ids[]", i2)
		}
	}

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

func ListIssues(cmdArgs map[string]string, _ map[int]string)  {
	var queryStrings = "state="
	if CommandArgExists(cmdArgs, "all") {
		queryStrings = ""
	} else if CommandArgExists(cmdArgs, "closed") {
		queryStrings += "closed&"
	} else {
		queryStrings += "opened&"
	}
	if CommandArgExists(cmdArgs, "label") || CommandArgExists(cmdArgs, "labels")  {
		queryStrings += "labels="+cmdArgs["label"]+"&"
	}
	if CommandArgExists(cmdArgs, "milestone")  {
		queryStrings += "milestone="+cmdArgs["milestone"]+"&"
	}
	if CommandArgExists(cmdArgs, "confidential")  {
		queryStrings += "confidential="+cmdArgs["confidential"]
	}
	queryStrings = strings.Trim(queryStrings,"& ")
	if len(queryStrings) > 0 {
		queryStrings = "?"+queryStrings
	}
	fmt.Println(queryStrings)
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
	} else {
		fmt.Println(resp)
	}
}

func DeleteIssue(cmdArgs map[string]string, arrFlags map[int]string)  {
	issueId := strings.Trim(arrFlags[1]," ")
	if CommandArgExists(cmdArgs, issueId) {
		arrIds := strings.Split(strings.Trim(issueId,"[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Deleting Issue #"+i2)
			queryStrings := "/"+i2
			resp := MakeRequest("{}","projects/"+GetEnv("GITLAB_PROJECT_ID")+"/issues"+queryStrings,"DELETE")
			if resp["responseCode"]==204 {
				bodyString := resp["responseMessage"]
				fmt.Println(bodyString)
				fmt.Println(Green("Issue Deleted Successfully"))
			} else if resp["responseCode"]==404 {
				fmt.Println(Red("Issue does not exist"))
			} else {
				fmt.Println(Red("Could not complete request."))
			}
			fmt.Println()
		}
	} else {
		fmt.Println(Red("Invalid command"))
		fmt.Println("Usage: glab issue delete <issue-id>")
	}
}

func ExecIssue(cmdArgs map[string]string, arrCmd map[int]string)  {
	commandList := map[interface{}]func(map[string]string,map[int]string) {
		"create" : CreateIssue,
		"list" : ListIssues,
		"delete" : DeleteIssue,
	}
	commandList[arrCmd[0]](cmdArgs, arrCmd)
}

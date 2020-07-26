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

func DisplayMergeRequest(hm map[string]interface{})  {
	duration := TimeAgo(hm["created_at"])
	if hm["state"] == "opened" {
		fmt.Println(Green(fmt.Sprint("#",hm["iid"])), hm["title"], Magenta(duration))
	} else {
		fmt.Println(Red(fmt.Sprint("#",hm["iid"])), hm["title"], Magenta(duration))
	}
}

func CreateMergeRequest(cmdArgs map[string]string, _ map[int]string)  {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(Cyan("Title"+"\n"+"-> "))
	mergeTitle, _ := reader.ReadString('\n')
	var mergeLabel string
	var sourceBranch string
	var targetBranch string
	if !CommandArgExists(cmdArgs, "labels") {
		fmt.Print(Cyan("Label(s) [Comma Separated]"+"\n"+"-> "))
		mergeLabel, _ = reader.ReadString('\n')
		mergeLabel = strings.Replace(mergeLabel, "\n", "", -1)
	} else {
		mergeLabel = strings.Trim(cmdArgs["labels"],"[] ")
	}
	mergeTitle = strings.Replace(mergeTitle, "\n", "", -1)
	fmt.Println()
	fmt.Println(Cyan("Enter Merge Request Description"), Yellow("[info: Type `exit` to close]"))
	var mergeDescription string
	for {
		fmt.Print(Cyan("-> "))
		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)
		if strings.Compare("exit", input) == 0 {
			break
		}
		mergeDescription += "\n"+input

	}
	if !CommandArgExists(cmdArgs, "source") {
		fmt.Print(Cyan("Source Branch"))
		fmt.Print(Yellow("-> "))
		sourceBranch, _ = reader.ReadString('\n')
	} else {
		sourceBranch = strings.Trim(cmdArgs["source"],"[] ")
	}
	if !CommandArgExists(cmdArgs, "target") {
		fmt.Print(Cyan("Target Branch"))
		fmt.Print(Yellow("-> "))
		targetBranch, _ = reader.ReadString('\n')
	} else {
		targetBranch = strings.Trim(cmdArgs["target"],"[] ")
	}
	targetBranch = strings.Replace(targetBranch, "\n", "", -1)
	sourceBranch = strings.Replace(targetBranch, "\n", "", -1)
	params := url.Values{}
	params.Add("title", mergeTitle)
	params.Add("description", mergeDescription)
	params.Add("labels", mergeLabel)
	params.Add("source_branch", sourceBranch)
	params.Add("target_branch", targetBranch)
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
	if CommandArgExists(cmdArgs, "allow-collaboration") {
		params.Add("allow_collaboration", cmdArgs["allow-collaboration"])
	}
	if CommandArgExists(cmdArgs, "remove-source-branch") {
		params.Add("remove_source_branch", cmdArgs["remove-source-branch"])
	}
	if CommandArgExists(cmdArgs, "target-project") {
		params.Add("target_project_id", cmdArgs["target-project"])
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
	fmt.Println(Yellow("Creating Merge Request {"+mergeTitle+"}..."))
	resp := MakeRequest(reqBody,"projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests","POST")

	if resp["responseCode"]==201 {
		bodyString := resp["responseMessage"]

		fmt.Println(Green("Merge Request created successfully"))
		if _, ok := bodyString.(string); ok {
			/* act on str */
			m := make(map[string]interface{})
			err := json.Unmarshal([]byte(bodyString.(string)), &m)
			if err != nil {
				log.Fatal(err)
			}
			DisplayMergeRequest(m)
		} else {
			/* not string */
		}
	} else {
		fmt.Println(resp["responseCode"], resp["responseMessage"])
	}
}

func ListMergeRequests(cmdArgs map[string]string, _ map[int]string)  {
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
	resp := MakeRequest("{}","projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests"+queryStrings,"GET")
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
				DisplayMergeRequest(hm)
			}
		} else {
			/* not string */
		}
	} else {
		fmt.Println(resp)
	}
}

func DeleteMergeRequest(cmdArgs map[string]string, arrFlags map[int]string)  {
	mergeId := strings.Trim(arrFlags[1]," ")
	if CommandArgExists(cmdArgs, mergeId) {
		arrIds := strings.Split(strings.Trim(mergeId,"[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Deleting Merge Request #"+i2)
			queryStrings := "/"+i2
			resp := MakeRequest("{}","projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests"+queryStrings,"DELETE")
			if resp["responseCode"]==204 {
				bodyString := resp["responseMessage"]
				fmt.Println(bodyString)
				fmt.Println(Green("Merge Request Deleted Successfully"))
			} else if resp["responseCode"]==404 {
				fmt.Println(Red("Merge Request does not exist"))
			} else {
				fmt.Println(Red("Could not complete request."))
			}
			fmt.Println()
		}
	} else {
		fmt.Println(Red("Invalid command"))
		fmt.Println("Usage: glab merge delete <merge-id>")
	}
}

func subscribeMergeRequest(cmdArgs map[string]string, arrFlags map[int]string)  {
	mergeId := strings.Trim(arrFlags[1]," ")
	if CommandArgExists(cmdArgs, mergeId) {
		arrIds := strings.Split(strings.Trim(mergeId,"[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Subscribing Merge Request #"+i2)
			queryStrings := "/"+i2+"/subscribe"
			resp := MakeRequest("{}","projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests"+queryStrings,"POST")
			if resp["responseCode"]==204 {
				bodyString := resp["responseMessage"]
				fmt.Println(bodyString)
				fmt.Println(Green("You have successfully subscribe to merge request with id #"+i2))
			} else if resp["responseCode"]==404 {
				fmt.Println(Red("Merge Request does not exist"))
			} else {
				fmt.Println(Red("Could not complete request."))
			}
			fmt.Println()
		}
	} else {
		fmt.Println(Red("Invalid command"))
		fmt.Println("Usage: glab merge delete <merge-id>")
	}
}

func unSubscribeMergeRequest(cmdArgs map[string]string, arrFlags map[int]string)  {
	mergeId := strings.Trim(arrFlags[1]," ")
	if CommandArgExists(cmdArgs, mergeId) {
		arrIds := strings.Split(strings.Trim(mergeId,"[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Unsubscribing Merge Request #"+i2)
			queryStrings := "/"+i2+"/unsubscribe"
			resp := MakeRequest("{}","projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests"+queryStrings,"POST")
			if resp["responseCode"]==204 {
				bodyString := resp["responseMessage"]
				fmt.Println(bodyString)
				fmt.Println(Green("You have successfully unsubscribed to merge request with id #"+i2))
			} else if resp["responseCode"]==404 {
				fmt.Println(Red("Merge Request does not exist"))
			} else {
				fmt.Println(Red("Could not complete request."))
			}
			fmt.Println()
		}
	} else {
		fmt.Println(Red("Invalid command"))
		fmt.Println("Usage: glab mr unsubscribe <merge-id>")
	}
}

func ExecMergeRequest(cmdArgs map[string]string, arrCmd map[int]string)  {
	commandList := map[interface{}]func(map[string]string,map[int]string) {
		"create" : CreateMergeRequest,
		"list" : ListMergeRequests,
		"delete" : DeleteMergeRequest,
		"subscribe" : subscribeMergeRequest,
	}
	commandList[arrCmd[0]](cmdArgs, arrCmd)
}

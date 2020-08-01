package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/logrusorgru/aurora"
	"log"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
)

func displayMergeRequest(hm map[string]interface{}) {
	duration := TimeAgo(hm["created_at"])
	if hm["state"] == "opened" {
		fmt.Println(aurora.Green(fmt.Sprint("#", hm["iid"])), hm["title"], aurora.Cyan("("+hm["source_branch"].(string)+")"), aurora.Magenta(duration))
	} else {
		fmt.Println(aurora.Red(fmt.Sprint("#", hm["iid"])), hm["title"], aurora.Cyan("("+hm["source_branch"].(string)+")"), aurora.Magenta(duration))
	}
}

func DisplayMultipleMergeRequests(m []interface{}) {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()
	if len(m) > 0 {
		fmt.Printf("Showing merge requests %d of %d on %s\n\n", len(m), len(m), GetEnv("GITLAB_REPO"))
		for i := 0; i < len(m); i++ {
			hm := m[i].(map[string]interface{})
			labels := hm["labels"]
			duration := TimeAgo(hm["created_at"])
			if hm["state"] == "opened" {
				_, _ = fmt.Fprintln(w, aurora.Green(fmt.Sprint(" #", hm["iid"])), "\t", hm["title"], "\t", aurora.Magenta(labels), "\t", aurora.Cyan("("+hm["source_branch"].(string)+")"), aurora.Magenta(duration))
			} else {
				_, _ = fmt.Fprintln(w, aurora.Red(fmt.Sprint(" #", hm["iid"])), "\t", hm["title"], "\t", aurora.Magenta(labels), "\t", aurora.Cyan("("+hm["source_branch"].(string)+")"), aurora.Magenta(duration))
			}
		}
	} else {
		fmt.Println("No merge requests available on " + GetEnv("GITLAB_REPO"))
	}
}

func createMergeRequest(cmdArgs map[string]string, _ map[int]string) {
	reader := bufio.NewReader(os.Stdin)
	var sourceBranch string
	var targetBranch string
	var mergeTitle string
	var mergeLabel string
	var mergeDescription string
	if !CommandArgExists(cmdArgs, "title") {
		fmt.Print(aurora.Cyan("Title" + "\n" + "-> "))
		mergeTitle, _ = reader.ReadString('\n')
	} else {
		mergeTitle = strings.Trim(cmdArgs["title"], " ")
	}
	if !CommandArgExists(cmdArgs, "label") {
		fmt.Print(aurora.Cyan("Label(s) [Comma Separated]" + "\n" + "-> "))
		mergeLabel, _ = reader.ReadString('\n')
	} else {
		mergeLabel = strings.Trim(cmdArgs["label"], "[] ")
	}
	mergeTitle = strings.ReplaceAll(mergeTitle, "\n", "")
	fmt.Println()
	if !CommandArgExists(cmdArgs, "description") {
		fmt.Println(aurora.Cyan("Description"), aurora.Yellow("[info: Type `exit` to close]"))
		for {
			fmt.Print(aurora.Cyan("-> "))
			input, _ := reader.ReadString('\n')
			// convert CRLF to LF
			input = strings.ReplaceAll(input, "\n", "")
			if strings.Compare("exit", input) == 0 {
				break
			}
			mergeDescription += "\n" + input
		}
	} else {
		mergeDescription = strings.Trim(cmdArgs["description"], " ")
	}
	if !CommandArgExists(cmdArgs, "source") {
		if CommandArgExists(cmdArgs, "create-branch") {
			sourceBranch = ReplaceNonAlphaNumericChars(mergeTitle, "-")
		} else {
			fmt.Print(aurora.Cyan("Source Branch"))
			fmt.Print(aurora.Yellow("-> "))
			sourceBranch, _ = reader.ReadString('\n')
		}
	} else {
		sourceBranch = strings.Trim(cmdArgs["source"], "[] ")
	}
	if !CommandArgExists(cmdArgs, "target") {
		fmt.Print(aurora.Cyan("Target Branch"))
		fmt.Print(aurora.Yellow("-> "))
		targetBranch, _ = reader.ReadString('\n')
	} else {
		targetBranch = strings.Trim(cmdArgs["target"], "[] ")
	}
	targetBranch = strings.ReplaceAll(targetBranch, "\n", "")
	sourceBranch = strings.ReplaceAll(sourceBranch, "\n", "")
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
		arrIds := strings.Split(strings.Trim(assignId, "[] "), ",")
		for _, i2 := range arrIds {
			params.Add("assignee_ids[]", i2)
		}
	}

	if CommandArgExists(cmdArgs, "create-branch") {
		minParams := url.Values{}
		minParams.Add("branch", sourceBranch)
		minParams.Add("ref", targetBranch)
		MakeRequest(minParams.Encode(), "projects/"+GetEnv("GITLAB_PROJECT_ID")+"/repository/branches", "POST")
	}

	reqBody := params.Encode()
	fmt.Println(aurora.Yellow("Creating Merge Request {" + mergeTitle + "}..."))
	resp := MakeRequest(reqBody, "projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests", "POST")

	if resp["responseCode"] == 201 {
		bodyString := resp["responseMessage"]

		fmt.Println(aurora.Green("Merge Request created successfully"))
		if _, ok := bodyString.(string); ok {
			/* act on str */
			m := make(map[string]interface{})
			err := json.Unmarshal([]byte(bodyString.(string)), &m)
			if err != nil {
				log.Fatal(err)
			}
			displayMergeRequest(m)
			fmt.Println()
		}
	} else {
		fmt.Println(resp["responseCode"], resp["responseMessage"])
	}
}

func acceptMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeId := strings.Trim(arrFlags[1], " ")
	params := url.Values{}
	if CommandArgExists(cmdArgs, "message") {
		params.Add("merge_commit_message", cmdArgs["message"])
	}
	if CommandArgExists(cmdArgs, "squash-message") {
		params.Add("squash_commit_message", cmdArgs["squash-message"])
	}
	if CommandArgExists(cmdArgs, "squash") {
		params.Add("squash", cmdArgs["squash"])
	}
	if CommandArgExists(cmdArgs, "remove-source-branch") {
		params.Add("should_remove_source_branch", cmdArgs["remove-source-branch"])
	}
	if CommandArgExists(cmdArgs, "when-pipeline-succeed") {
		params.Add("merge_when_pipeline_succeed", cmdArgs["when-pipeline-succeed"])
	}
	if CommandArgExists(cmdArgs, "sha") {
		params.Add("sha", cmdArgs["sha"])
	}

	reqBody := params.Encode()
	fmt.Println(aurora.Yellow("Accepting Merge Request #" + mergeId + "..."))
	resp := MakeRequest(reqBody, "projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests/"+mergeId+"/merge", "PUT")

	if resp["responseCode"] == 200 {
		bodyString := resp["responseMessage"]

		fmt.Println(aurora.Green("Merge Request accepted successfully"))
		if _, ok := bodyString.(string); ok {
			/* act on str */
			m := make(map[string]interface{})
			err := json.Unmarshal([]byte(bodyString.(string)), &m)
			if err != nil {
				log.Fatal(err)
			}
			displayMergeRequest(m)
			fmt.Println()
		}
	} else if resp["responseCode"] == 405 {
		fmt.Println("Merge request cannot be merged")
	} else if resp["responseCode"] == 401 {
		fmt.Println("You don't have enough permission to accept this merge request")
	} else if resp["responseCode"] == 406 {
		fmt.Println("Branch cannot be merged. There are merge conflicts.")
	} else {
		fmt.Println(resp["responseCode"], resp["responseMessage"])
	}
}

func listMergeRequests(cmdArgs map[string]string, _ map[int]string) {
	var queryStrings = "state="
	if CommandArgExists(cmdArgs, "all") {
		queryStrings = ""
	} else if CommandArgExists(cmdArgs, "closed") {
		queryStrings += "closed&"
	} else {
		queryStrings += "opened&"
	}
	if CommandArgExists(cmdArgs, "label") || CommandArgExists(cmdArgs, "labels") {
		queryStrings += "labels=" + cmdArgs["label"] + "&"
	}
	if CommandArgExists(cmdArgs, "milestone") {
		queryStrings += "milestone=" + cmdArgs["milestone"] + "&"
	}
	if CommandArgExists(cmdArgs, "confidential") {
		queryStrings += "confidential=" + cmdArgs["confidential"]
	}
	queryStrings = strings.Trim(queryStrings, "& ")
	if len(queryStrings) > 0 {
		queryStrings = "?" + queryStrings
	}
	resp := MakeRequest("{}", "projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests"+queryStrings, "GET")
	//fmt.Println(resp)
	if resp["responseCode"] == 200 {
		bodyString := resp["responseMessage"]
		if _, ok := bodyString.(string); ok {
			/* act on str */
			var m []interface{}
			err := json.Unmarshal([]byte(bodyString.(string)), &m)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println()
			DisplayMultipleMergeRequests(m)
			fmt.Println()
		}
	} else {
		fmt.Println(resp)
	}
}

func issuesRelatedMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	var queryStrings = "state="
	mergeId := strings.Trim(arrFlags[1], " ")
	if CommandArgExists(cmdArgs, "all") {
		queryStrings = ""
	} else if CommandArgExists(cmdArgs, "closed") {
		queryStrings += "closed&"
	} else {
		queryStrings += "opened&"
	}
	if CommandArgExists(cmdArgs, "label") || CommandArgExists(cmdArgs, "labels") {
		queryStrings += "labels=" + cmdArgs["label"] + "&"
	}
	if CommandArgExists(cmdArgs, "milestone") {
		queryStrings += "milestone=" + cmdArgs["milestone"] + "&"
	}
	if CommandArgExists(cmdArgs, "confidential") {
		queryStrings += "confidential=" + cmdArgs["confidential"]
	}
	queryStrings = strings.Trim(queryStrings, "& ")
	if len(queryStrings) > 0 {
		queryStrings = "?" + queryStrings
	}
	resp := MakeRequest("{}", "projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests/"+mergeId+"/closes_issues"+queryStrings, "GET")
	//fmt.Println(resp)
	if resp["responseCode"] == 200 {
		bodyString := resp["responseMessage"]
		if _, ok := bodyString.(string); ok {
			/* act on str */
			var m []interface{}
			err := json.Unmarshal([]byte(bodyString.(string)), &m)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println()
			displayMultipleIssues(m)
			fmt.Println()
		}
	} else {
		fmt.Println(resp)
	}
}

func deleteMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeId := strings.Trim(arrFlags[1], " ")
	if CommandArgExists(cmdArgs, mergeId) {
		arrIds := strings.Split(strings.Trim(mergeId, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Deleting Merge Request #" + i2)
			queryStrings := "/" + i2
			resp := MakeRequest("{}", "projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests"+queryStrings, "DELETE")
			if resp["responseCode"] == 204 {
				bodyString := resp["responseMessage"]
				fmt.Println(bodyString)
				fmt.Println(aurora.Green("Merge Request Deleted Successfully"))
			} else if resp["responseCode"] == 404 {
				fmt.Println(aurora.Red("Merge Request does not exist"))
			} else {
				fmt.Println(aurora.Red("Could not complete request."))
			}
			fmt.Println()
		}
	} else {
		fmt.Println(aurora.Red("Invalid command"))
		fmt.Println("Usage: glab merge delete <merge-id>")
	}
}

func subscribeMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeId := strings.Trim(arrFlags[1], " ")
	if CommandArgExists(cmdArgs, mergeId) {
		arrIds := strings.Split(strings.Trim(mergeId, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Subscribing Merge Request #" + i2)
			queryStrings := "/" + i2 + "/subscribe"
			resp := MakeRequest("{}", "projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests"+queryStrings, "POST")
			if resp["responseCode"] == 204 {
				bodyString := resp["responseMessage"]
				fmt.Println(bodyString)
				fmt.Println(aurora.Green("You have successfully subscribe to merge request with id #" + i2))
			} else if resp["responseCode"] == 404 {
				fmt.Println(aurora.Red("Merge Request does not exist"))
			} else {
				fmt.Println(aurora.Red("Could not complete request."))
			}
			fmt.Println()
		}
	} else {
		fmt.Println(aurora.Red("Invalid command"))
		fmt.Println("Usage: glab merge delete <merge-id>")
	}
}

func unsubscribeMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeId := strings.Trim(arrFlags[1], " ")
	if CommandArgExists(cmdArgs, mergeId) {
		arrIds := strings.Split(strings.Trim(mergeId, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Unsubscribing Merge Request #" + i2)
			queryStrings := "/" + i2 + "/unsubscribe"
			resp := MakeRequest("{}", "projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests"+queryStrings, "POST")
			if resp["responseCode"] == 204 {
				bodyString := resp["responseMessage"]
				fmt.Println(bodyString)
				fmt.Println(aurora.Green("You have successfully unsubscribed to merge request with id #" + i2))
			} else if resp["responseCode"] == 404 {
				fmt.Println(aurora.Red("Merge Request does not exist"))
			} else {
				fmt.Println(aurora.Red("Could not complete request."))
			}
			fmt.Println()
		}
	} else {
		fmt.Println(aurora.Red("Invalid command"))
		fmt.Println("Usage: glab mr unsubscribe <merge-id>")
	}
}

func changeMergeRequestState(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeId := strings.Trim(arrFlags[1], " ")
	if CommandArgExists(cmdArgs, mergeId) {
		reqType := arrFlags[0]
		params := url.Values{}
		mergeMessage := ""
		if reqType == "close" {
			params.Add("state_event", "close")
			mergeMessage = "closed"
		} else {
			params.Add("state_event", "reopen")
			mergeMessage = "opened"
		}
		arrIds := strings.Split(strings.Trim(mergeId, "[] "), ",")
		reqBody := params.Encode()
		for _, i2 := range arrIds {
			fmt.Println("...")
			resp := MakeRequest(reqBody, "projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests/"+i2, "PUT")
			if resp["responseCode"] == 200 {
				fmt.Println(aurora.Green("You have successfully " + mergeMessage + " to merge request with id #" + i2))
			} else if resp["responseCode"] == 404 {
				fmt.Println(aurora.Red("Merge Request does not exist"))
			} else {
				fmt.Println("Could not complete request")
				fmt.Println(resp["responseCode"], resp["responseMessage"])
			}
			fmt.Println()
		}
	} else {
		fmt.Println(aurora.Red("Invalid command"))
		fmt.Println("Usage: glab mr <state> <merge-id>")
	}
}

// ExecMergeRequest is exported
func ExecMergeRequest(cmdArgs map[string]string, arrCmd map[int]string) {
	commandList := map[interface{}]func(map[string]string, map[int]string){
		"create":      createMergeRequest,
		"list":        listMergeRequests,
		"ls":          listMergeRequests,
		"delete":      deleteMergeRequest,
		"subscribe":   subscribeMergeRequest,
		"unsubscribe": unsubscribeMergeRequest,
		"accept":      acceptMergeRequest,
		"merge":       acceptMergeRequest,
		"close":       changeMergeRequestState,
		"reopen":      changeMergeRequestState,
		"issues":      issuesRelatedMergeRequest,
	}
	if _, ok := commandList[arrCmd[0]]; ok {
		commandList[arrCmd[0]](cmdArgs, arrCmd)
	} else {
		fmt.Println(arrCmd[0]+":", "Invalid Command")
	}
}

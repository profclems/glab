package commands

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/xanzy/go-gitlab"
	"log"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
)

func displayMergeRequest(hm *gitlab.MergeRequest) {
	duration := TimeAgo(*hm.CreatedAt)
	if hm.State == "opened" {
		fmt.Println(aurora.Green(fmt.Sprint("#", hm.IID)), hm.Title, aurora.Cyan("("+hm.SourceBranch+")"), aurora.Magenta(duration))
	} else {
		fmt.Println(aurora.Red(fmt.Sprint("#", hm.IID)), hm.Title, aurora.Cyan("("+hm.SourceBranch+")"), aurora.Magenta(duration))
	}
}

func displayAllMergeRequests(m []*gitlab.MergeRequest) {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()
	if len(m) > 0 {
		fmt.Printf("Showing mergeRequests %d of %d on %s\n\n", len(m), len(m), getRepo())
		for _, issue := range m {
			labels := issue.Labels
			duration := TimeAgo(*issue.CreatedAt)
			if issue.State == "opened" {
				_, _ = fmt.Fprintln(w, aurora.Green(fmt.Sprint("#", issue.IID)), "\t", issue.Title, "\t", aurora.Magenta(labels), "\t", aurora.Magenta(duration))
			} else {
				_, _ = fmt.Fprintln(w, aurora.Red(fmt.Sprint("#", issue.IID)), "\t", issue.Title, "\t", aurora.Magenta(labels), "\t", aurora.Magenta(duration))
			}
		}
	} else {
		fmt.Println("No Merge Requests available on " + getRepo())
	}
}

func createMergeRequest(cmdArgs map[string]string, _ map[int]string) {
	l := &gitlab.CreateMergeRequestOptions{}
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
	l.Title = gitlab.String(mergeTitle)
	l.Description = gitlab.String(mergeDescription)
	l.Labels = &gitlab.Labels{mergeLabel}
	l.SourceBranch = gitlab.String(sourceBranch)
	l.TargetBranch = gitlab.String(targetBranch)
	if CommandArgExists(cmdArgs, "milestone") {
		l.MilestoneID = gitlab.Int(stringToInt(cmdArgs["milestone"]))
	}
	if CommandArgExists(cmdArgs, "allow-collaboration") {
		l.AllowCollaboration = gitlab.Bool(true)
	}
	if CommandArgExists(cmdArgs, "remove-source-branch") {
		l.RemoveSourceBranch = gitlab.Bool(true)
	}
	if CommandArgExists(cmdArgs, "target-project") {
		l.TargetProjectID = gitlab.Int(stringToInt(cmdArgs["target-project"]))
	}
	if CommandArgExists(cmdArgs, "assigns") {
		assignID := cmdArgs["assigns"]
		arrIds := strings.Split(strings.Trim(assignID, "[] "), ",")
		var t2 []int

		for _, i := range arrIds {
			j := stringToInt(i)
			t2 = append(t2, j)
		}
		l.AssigneeIDs = t2
	}

	git, repo := InitGitlabClient()
	if CommandArgExists(cmdArgs, "create-branch") {
		lb := &gitlab.CreateBranchOptions{
			Branch: gitlab.String(sourceBranch),
			Ref:    gitlab.String(targetBranch),
		}
		fmt.Println("Creating related branch...")
		branch, resp, _ := git.Branches.CreateBranch(repo, lb)
		if resp.StatusCode == 201 {
			fmt.Println("Branch created: ", branch.WebURL)
		} else {
			fmt.Println("Error creating branch: ", resp.Status)
		}

	}

	mr, _, err := git.MergeRequests.CreateMergeRequest(repo, l)
	if err != nil {
		log.Fatal(err)
	}
	displayMergeRequest(mr)
}

func acceptMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")
	params := url.Values{}
	l := &gitlab.AcceptMergeRequestOptions{}
	if CommandArgExists(cmdArgs, "message") {
		l.MergeCommitMessage = gitlab.String(cmdArgs["message"])
	}
	if CommandArgExists(cmdArgs, "squash-message") {
		l.SquashCommitMessage = gitlab.String(cmdArgs["squash-message"])
	}
	if CommandArgExists(cmdArgs, "squash") {
		l.Squash = gitlab.Bool(true)
	}
	if CommandArgExists(cmdArgs, "remove-source-branch") {
		l.ShouldRemoveSourceBranch = gitlab.Bool(true)
	}
	if CommandArgExists(cmdArgs, "when-pipeline-succeed") {
		l.MergeWhenPipelineSucceeds = gitlab.Bool(true)
	}
	if CommandArgExists(cmdArgs, "sha") {
		params.Add("sha", cmdArgs["sha"])
		l.SHA = gitlab.String(cmdArgs["sha"])
	}
	git, repo := InitGitlabClient()
	mr, resp, err := git.MergeRequests.AcceptMergeRequest(repo, stringToInt(mergeID), l)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(aurora.Yellow("Accepting Merge Request #" + mergeID + "..."))

	if resp.StatusCode == 200 {
		fmt.Println(aurora.Green("Merge Request accepted successfully"))
		displayMergeRequest(mr)
	} else if resp.StatusCode == 405 {
		fmt.Println("Merge request cannot be merged")
	} else if resp.StatusCode == 401 {
		fmt.Println("You don't have enough permission to accept this merge request")
	} else if resp.StatusCode == 406 {
		fmt.Println("Branch cannot be merged. There are merge conflicts.")
	} else {
		fmt.Println(resp)
	}
}

func approveMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")
	l := &gitlab.ApproveMergeRequestOptions{}
	if CommandArgExists(cmdArgs, "sha") {
		l.SHA = gitlab.String(cmdArgs["sha"])
	}

	fmt.Println(aurora.Yellow("Approving Merge Request #" + mergeID + "..."))
	git, repo := InitGitlabClient()
	_, resp, err := git.MergeRequestApprovals.ApproveMergeRequest(repo, stringToInt(mergeID), l)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		fmt.Println(aurora.Green("Merge Request approved successfully"))
	} else if resp.StatusCode == 405 {
		fmt.Println("Merge request cannot be approved")
	} else if resp.StatusCode == 401 {
		fmt.Println("You don't have enough permission to approve this merge request")
	} else {
		fmt.Println(resp)
	}
}

func listMergeRequests(cmdArgs map[string]string, _ map[int]string) {
	var state = "all"
	if CommandArgExists(cmdArgs, "closed") {
		state = "closed"
	} else {
		state = "opened"
	}

	l := &gitlab.ListProjectMergeRequestsOptions{
		State: gitlab.String(state),
	}
	if CommandArgExists(cmdArgs, "label") || CommandArgExists(cmdArgs, "labels") {
		label := gitlab.Labels{
			cmdArgs["label"],
		}
		l.Labels = &label
	}
	if CommandArgExists(cmdArgs, "milestone") {
		l.Milestone = gitlab.String(cmdArgs["milestone"])
	}

	git, repo := InitGitlabClient()
	// Create new label
	mergeRequests, _, err := git.MergeRequests.ListProjectMergeRequests(repo, l)
	if err != nil {
		log.Fatal(err)
	}
	displayAllMergeRequests(mergeRequests)
}

func issuesRelatedMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")
	l := &gitlab.GetIssuesClosedOnMergeOptions{}
	git, repo := InitGitlabClient()
	mr, _, err := git.MergeRequests.GetIssuesClosedOnMerge(repo, stringToInt(mergeID), l)
	if err != nil {
		log.Fatal(err)
	}
	displayAllIssues(mr)
}

func deleteMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")
	git, repo := InitGitlabClient()

	if CommandArgExists(cmdArgs, mergeID) {
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Deleting Merge Request #" + i2)
			issue, _ := git.MergeRequests.DeleteMergeRequest(repo, stringToInt(i2))
			if issue.StatusCode == 204 {
				fmt.Println(aurora.Green("Merge Request Deleted Successfully"))
			} else if issue.StatusCode == 404 {
				fmt.Println(aurora.Red("Merge Request does not exist"))
			} else {
				fmt.Println(aurora.Red("Could not complete request."))
			}
		}
	} else {
		fmt.Println(aurora.Red("Invalid command"))
		fmt.Println("Usage: glab issue delete <issue-id>")
	}
}

func subscribeMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")
	git, repo := InitGitlabClient()
	if CommandArgExists(cmdArgs, mergeID) {
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Subscribing Merge Request #" + i2)
			issue, resp, _ := git.MergeRequests.SubscribeToMergeRequest(repo, stringToInt(i2), nil)

			if resp.StatusCode == 204 {
				bodyString := resp.Body
				fmt.Println(bodyString)
				fmt.Println(aurora.Green("You have successfully subscribe to merge request #" + i2))
				displayMergeRequest(issue)
			} else if resp.StatusCode == 404 {
				fmt.Println(aurora.Red("MergeRequest does not exist"))
			} else {
				fmt.Println(aurora.Red("Could not complete request."), resp.Status)
			}
		}
	} else {
		fmt.Println(aurora.Red("Invalid command"))
		fmt.Println("Usage: glab issue subscribe <issue-id>")
	}
}

func unsubscribeMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")
	git, repo := InitGitlabClient()
	if CommandArgExists(cmdArgs, mergeID) {
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Unsubscribing Merge Request #" + i2)
			issue, resp, _ := git.MergeRequests.UnsubscribeFromMergeRequest(repo, stringToInt(i2))

			if resp.StatusCode == 204 {
				bodyString := resp.Body
				fmt.Println(bodyString)
				fmt.Println(aurora.Green("You have successfully unsubscribe to issue #" + i2))
				displayMergeRequest(issue)
			} else if resp.StatusCode == 404 {
				fmt.Println(aurora.Red("MergeRequest does not exist"))
			} else {
				fmt.Println(aurora.Red("Could not complete request."), resp.Status)
			}
		}
	} else {
		fmt.Println(aurora.Red("Invalid command"))
		fmt.Println("Usage: glab issue unsubscribe <issue-id>")
	}
}

func changeMergeRequestState(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")
	git, repo := InitGitlabClient()
	if CommandArgExists(cmdArgs, mergeID) {
		reqType := arrFlags[0]
		l := &gitlab.UpdateMergeRequestOptions{}
		if reqType == "close" {
			l.StateEvent = gitlab.String("close")
		} else {
			l.StateEvent = gitlab.String("reopen")
		}
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Updating Merge request...")
			mr, resp, err := git.MergeRequests.UpdateMergeRequest(repo, stringToInt(i2), l)
			if err != nil {
				log.Fatal(err)
			}
			if resp.StatusCode == 200 {
				fmt.Println(aurora.Green("You have successfully updated merge request #" + i2))
				displayMergeRequest(mr)
			} else if resp.StatusCode == 404 {
				fmt.Println(aurora.Red("MergeRequest does not exist"))
			} else {
				fmt.Println("Could not complete request: ", resp)
			}
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
		"approve":     approveMergeRequest,
	}
	if _, ok := commandList[arrCmd[0]]; ok {
		commandList[arrCmd[0]](cmdArgs, arrCmd)
	} else {
		fmt.Println(arrCmd[0]+":", "Invalid Command")
	}
}

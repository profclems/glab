package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/logrusorgru/aurora"
	"github.com/xanzy/go-gitlab"
)

func displayAllIssues(m []*gitlab.Issue) {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()
	if len(m) > 0 {
		fmt.Printf("Showing issues %d of %d on %s\n\n", len(m), len(m), getRepo())
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
		fmt.Println("No Issues available on " + getRepo())
	}
}

func displayIssue(hm *gitlab.Issue) {
	duration := TimeAgo(*hm.CreatedAt)
	if hm.State == "opened" {
		fmt.Println(aurora.Green(fmt.Sprint("#", hm.IID)), hm.Title, aurora.Magenta(duration))
	} else {
		fmt.Println(aurora.Red(fmt.Sprint("#", hm.IID)), hm.Title, aurora.Magenta(duration))
	}
}

func createIssue(cmdArgs map[string]string, _ map[int]string) {
	l := &gitlab.CreateIssueOptions{}
	reader := bufio.NewReader(os.Stdin)
	var issueTitle string
	var issueLabel string
	var issueDescription string
	if !CommandArgExists(cmdArgs, "title") {
		fmt.Print(aurora.Cyan("Title" + "\n" + "-> "))
		issueTitle, _ = reader.ReadString('\n')
	} else {
		issueTitle = strings.Trim(cmdArgs["title"], " ")
	}
	if !CommandArgExists(cmdArgs, "label") {
		fmt.Print(aurora.Cyan("Label(s) [Comma Separated]" + "\n" + "-> "))
		issueLabel, _ = reader.ReadString('\n')
	} else {
		issueLabel = strings.Trim(cmdArgs["label"], "[] ")
	}
	issueTitle = strings.ReplaceAll(issueTitle, "\n", "")
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
			issueDescription += "\n" + input
		}
	} else {
		issueDescription = strings.Trim(cmdArgs["description"], " ")
	}
	fmt.Print(aurora.Cyan("Due Date"))
	fmt.Print(aurora.Yellow("(Format: YYYY-MM-DD)" + "\n" + "-> "))
	//issueDue, _ := reader.ReadString('\n')
	l.Title = gitlab.String(issueTitle)
	l.Labels = &gitlab.Labels{issueLabel}
	l.Description = &issueDescription
	l.DueDate = &gitlab.ISOTime{}
	if CommandArgExists(cmdArgs, "confidential") {
		l.Confidential = gitlab.Bool(true)
	}
	if CommandArgExists(cmdArgs, "weight") {
		l.Weight = gitlab.Int(stringToInt(cmdArgs["weight"]))
	}
	if CommandArgExists(cmdArgs, "mr") {
		l.MergeRequestToResolveDiscussionsOf = gitlab.Int(stringToInt(cmdArgs["mr"]))
	}
	if CommandArgExists(cmdArgs, "milestone") {
		l.MilestoneID = gitlab.Int(stringToInt(cmdArgs["milestone"]))
	}
	if CommandArgExists(cmdArgs, "resolved-by-merge-request") {
		l.MergeRequestToResolveDiscussionsOf = gitlab.Int(stringToInt(cmdArgs["resolved-by-merge"]))
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
	issue, _, err := git.Issues.CreateIssue(repo, l)
	if err != nil {
		log.Fatal(err)
	}
	displayIssue(issue)
}

func listIssues(cmdArgs map[string]string, _ map[int]string) {
	var state = "all"
	if CommandArgExists(cmdArgs, "closed") {
		state = "closed"
	} else {
		state = "opened"
	}

	l := &gitlab.ListProjectIssuesOptions{
		State: gitlab.String(state),
	}
	if CommandArgExists(cmdArgs, "label") || CommandArgExists(cmdArgs, "labels") {
		label := gitlab.Labels{
			cmdArgs["label"],
		}
		l.Labels = label
	}
	if CommandArgExists(cmdArgs, "milestone") {
		l.Milestone = gitlab.String(cmdArgs["milestone"])
	}
	if CommandArgExists(cmdArgs, "confidential") {
		var confidential bool
		if cmdArgs["confidential"] == "true" {
			confidential = true
		}
		l.Confidential = gitlab.Bool(confidential)
	}

	git, repo := InitGitlabClient()
	// Create new label
	issues, _, err := git.Issues.ListProjectIssues(repo, l)
	if err != nil {
		log.Fatal(err)
	}
	displayAllIssues(issues)

}

func deleteIssue(cmdArgs map[string]string, arrFlags map[int]string) {
	issueID := strings.Trim(arrFlags[1], " ")
	git, repo := InitGitlabClient()

	if CommandArgExists(cmdArgs, issueID) {
		arrIds := strings.Split(strings.Trim(issueID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Deleting Issue #" + i2)
			issue, _ := git.Issues.DeleteIssue(repo, stringToInt(i2))
			if issue.StatusCode == 204 {
				fmt.Println(aurora.Green("Issue Deleted Successfully"))
			} else if issue.StatusCode == 404 {
				fmt.Println(aurora.Red("Issue does not exist"))
			} else {
				fmt.Println(aurora.Red("Could not complete request."))
			}
		}
	} else {
		fmt.Println(aurora.Red("Invalid command"))
		fmt.Println("Usage: glab issue delete <issue-id>")
	}
}

func subscribeIssue(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")
	git, repo := InitGitlabClient()
	if CommandArgExists(cmdArgs, mergeID) {
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Subscribing Issue #" + i2)
			issue, resp, _ := git.Issues.SubscribeToIssue(repo, stringToInt(i2), nil)

			if resp.StatusCode == 204 {
				bodyString := resp.Body
				fmt.Println(bodyString)
				fmt.Println(aurora.Green("You have successfully subscribe to issue #" + i2))
				displayIssue(issue)
			} else if resp.StatusCode == 404 {
				fmt.Println(aurora.Red("Issue does not exist"))
			} else {
				fmt.Println(aurora.Red("Could not complete request."), resp.Status)
			}
		}
	} else {
		fmt.Println(aurora.Red("Invalid command"))
		fmt.Println("Usage: glab issue subscribe <issue-id>")
	}
}

func unsubscribeIssue(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")
	git, repo := InitGitlabClient()
	if CommandArgExists(cmdArgs, mergeID) {
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Unsubscribing Issue #" + i2)
			issue, resp, _ := git.Issues.UnsubscribeFromIssue(repo, stringToInt(i2))

			if resp.StatusCode == 204 {
				bodyString := resp.Body
				fmt.Println(bodyString)
				fmt.Println(aurora.Green("You have successfully unsubscribe to issue #" + i2))
				displayIssue(issue)
			} else if resp.StatusCode == 404 {
				fmt.Println(aurora.Red("Issue does not exist"))
			} else {
				fmt.Println(aurora.Red("Could not complete request."), resp.Status)
			}
		}
	} else {
		fmt.Println(aurora.Red("Invalid command"))
		fmt.Println("Usage: glab issue unsubscribe <issue-id>")
	}
}

func changeIssueState(cmdArgs map[string]string, arrFlags map[int]string) {
	issueID := strings.Trim(arrFlags[1], " ")
	git, repo := InitGitlabClient()
	if CommandArgExists(cmdArgs, issueID) {
		reqType := arrFlags[0]
		l := &gitlab.UpdateIssueOptions{}
		if reqType == "close" {
			l.StateEvent = gitlab.String("close")
		} else {
			l.StateEvent = gitlab.String("reopen")
		}
		arrIds := strings.Split(strings.Trim(issueID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Updating Issue")
			issue, resp, err := git.Issues.UpdateIssue(repo, stringToInt(i2), l)
			if err != nil {
				log.Fatal(err)
			}
			if resp.StatusCode == 200 {
				fmt.Println(aurora.Green("You have successfully updated issue #" + i2))
				displayIssue(issue)
			} else if resp.StatusCode == 404 {
				fmt.Println(aurora.Red("Issue does not exist"))
			} else {
				fmt.Println("Could not complete request: ", resp)
			}
		}
	} else {
		fmt.Println(aurora.Red("Invalid command"))
		fmt.Println("Usage: glab issue <state> <merge-id>")
	}
}

// ExecIssue is exported
func ExecIssue(cmdArgs map[string]string, arrCmd map[int]string) {
	commandList := map[interface{}]func(map[string]string, map[int]string){
		"create":      createIssue,
		"list":        listIssues,
		"ls":          listIssues,
		"delete":      deleteIssue,
		"subscribe":   subscribeIssue,
		"unsubscribe": unsubscribeIssue,
		"open":        changeIssueState,
		"reopen":      changeIssueState,
		"close":       changeIssueState,
	}
	if _, ok := commandList[arrCmd[0]]; ok {
		commandList[arrCmd[0]](cmdArgs, arrCmd)
	} else {
		fmt.Println(arrCmd[0]+":", "Invalid Command")
	}
}

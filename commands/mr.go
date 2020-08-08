package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"glab/commands"
	"glab/internal/git"
	"glab/internal/manip"
	"log"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/logrusorgru/aurora"
	"github.com/xanzy/go-gitlab"
)

func displayMergeRequest(hm *gitlab.MergeRequest) {
	duration := manip.TimeAgo(*hm.CreatedAt)
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
		fmt.Printf("Showing mergeRequests %d of %d on %s\n\n", len(m), len(m), git.GetRepo())
		for _, issue := range m {
			labels := issue.Labels
			duration := manip.TimeAgo(*issue.CreatedAt)
			if issue.State == "opened" {
				_, _ = fmt.Fprintln(w, aurora.Green(fmt.Sprint("#", issue.IID)), "\t", issue.Title, "\t", aurora.Magenta(labels), "\t", aurora.Magenta(duration))
			} else {
				_, _ = fmt.Fprintln(w, aurora.Red(fmt.Sprint("#", issue.IID)), "\t", issue.Title, "\t", aurora.Magenta(labels), "\t", aurora.Magenta(duration))
			}
		}
	} else {
		fmt.Println("No Merge Requests available on " + git.GetRepo())
	}
}


func acceptMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")
	params := url.Values{}
	l := &gitlab.AcceptMergeRequestOptions{}
	if manip.CommandArgExists(cmdArgs, "message") {
		l.MergeCommitMessage = gitlab.String(cmdArgs["message"])
	}
	if manip.CommandArgExists(cmdArgs, "squash-message") {
		l.SquashCommitMessage = gitlab.String(cmdArgs["squash-message"])
	}
	if manip.CommandArgExists(cmdArgs, "squash") {
		l.Squash = gitlab.Bool(true)
	}
	if manip.CommandArgExists(cmdArgs, "remove-source-branch") {
		l.ShouldRemoveSourceBranch = gitlab.Bool(true)
	}
	if manip.CommandArgExists(cmdArgs, "when-pipeline-succeed") {
		l.MergeWhenPipelineSucceeds = gitlab.Bool(true)
	}
	if manip.CommandArgExists(cmdArgs, "sha") {
		params.Add("sha", cmdArgs["sha"])
		l.SHA = gitlab.String(cmdArgs["sha"])
	}
	gitlabClient, repo := git.InitGitlabClient()
	mr, resp, _ := gitlabClient.MergeRequests.AcceptMergeRequest(repo, manip.StringToInt(mergeID), l)

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
	if manip.CommandArgExists(cmdArgs, "sha") {
		l.SHA = gitlab.String(cmdArgs["sha"])
	}

	fmt.Println(aurora.Yellow("Approving Merge Request #" + mergeID + "..."))
	gitlabClient, repo := git.InitGitlabClient()
	_, resp, _ := gitlabClient.MergeRequestApprovals.ApproveMergeRequest(repo, manip.StringToInt(mergeID), l)
	if resp != nil {
		if resp.StatusCode == 201 {
			fmt.Println(aurora.Green("Merge Request approved successfully"))
		} else if resp.StatusCode == 405 {
			fmt.Println("Merge request cannot be approved")
		} else if resp.StatusCode == 401 {
			fmt.Println("Merge request already approved or you don't have enough permission to approve this merge request")
		} else {
			fmt.Println(resp.Status)
		}
	} else {
		fmt.Println(resp)
	}
}

func revokeMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")

	fmt.Println(aurora.Yellow("Revoking approval for Merge Request #" + mergeID + "..."))
	gitlabClient, repo := git.InitGitlabClient()
	resp, _ := gitlabClient.MergeRequestApprovals.UnapproveMergeRequest(repo, manip.StringToInt(mergeID))
	if resp != nil {
		if resp.StatusCode == 201 {
			fmt.Println(aurora.Green("Merge Request approval revoked successfully"))
		} else if resp.StatusCode == 405 {
			fmt.Println("Merge request cannot be unapproved")
		} else if resp.StatusCode == 401 {
			fmt.Println("Merge request already unapproved or you don't have enough permission to unapprove this merge request")
		} else {
			fmt.Println(resp.Status)
		}
	} else {
		fmt.Println(resp)
	}
}

func listMergeRequests(cmdArgs map[string]string, _ map[int]string) {
	var state = "all"
	if manip.CommandArgExists(cmdArgs, "closed") {
		state = "closed"
	} else {
		state = "opened"
	}

	l := &gitlab.ListProjectMergeRequestsOptions{
		State: gitlab.String(state),
	}
	if manip.CommandArgExists(cmdArgs, "label") || manip.CommandArgExists(cmdArgs, "labels") {
		label := gitlab.Labels{
			cmdArgs["label"],
		}
		l.Labels = &label
	}
	if manip.CommandArgExists(cmdArgs, "milestone") {
		l.Milestone = gitlab.String(cmdArgs["milestone"])
	}

	gitlabClient, repo := git.InitGitlabClient()
	// Create new label
	mergeRequests, _, err := gitlabClient.MergeRequests.ListProjectMergeRequests(repo, l)
	if err != nil {
		log.Fatal(err)
	}
	displayAllMergeRequests(mergeRequests)
}

/*
func issuesRelatedMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")
	l := &gitlab.GetIssuesClosedOnMergeOptions{}
	gitlabClient, repo := git.InitGitlabClient()
	mr, _, err := gitlabClient.MergeRequests.GetIssuesClosedOnMerge(repo, manip.StringToInt(mergeID), l)
	if err != nil {
		log.Fatal(err)
	}
	displayAllIssues(mr)
}
 */

func updateMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")
	l := &gitlab.UpdateMergeRequestOptions{}
	if manip.CommandArgExists(cmdArgs, "title") {
		l.Title = gitlab.String(cmdArgs["title"])
	}
	if manip.CommandArgExists(cmdArgs, "lock-discussion") {
		l.DiscussionLocked = gitlab.Bool(true)
	}
	if manip.CommandArgExists(cmdArgs, "description") {
		l.Description = gitlab.String(cmdArgs["description"])
	}
	gitlabClient, repo := git.InitGitlabClient()
	mr, _, err := gitlabClient.MergeRequests.UpdateMergeRequest(repo, manip.StringToInt(mergeID), l)
	if err != nil {
		log.Fatal(err)
	}
	displayMergeRequest(mr)
}

func deleteMergeRequest(cmdArgs map[string]string, arrFlags map[int]string) {
	mergeID := strings.Trim(arrFlags[1], " ")
	gitlabClient, repo := git.InitGitlabClient()

	if manip.CommandArgExists(cmdArgs, mergeID) {
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Deleting Merge Request #" + i2)
			issue, _ := gitlabClient.MergeRequests.DeleteMergeRequest(repo, manip.StringToInt(i2))
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
	gitlabClient, repo := git.InitGitlabClient()
	if manip.CommandArgExists(cmdArgs, mergeID) {
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Subscribing Merge Request #" + i2)
			issue, resp, _ := gitlabClient.MergeRequests.SubscribeToMergeRequest(repo, manip.StringToInt(i2), nil)

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
	gitlabClient, repo := git.InitGitlabClient()
	if manip.CommandArgExists(cmdArgs, mergeID) {
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Unsubscribing Merge Request #" + i2)
			issue, resp, _ := gitlabClient.MergeRequests.UnsubscribeFromMergeRequest(repo, manip.StringToInt(i2))

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
	gitlabClient, repo := git.InitGitlabClient()
	if manip.CommandArgExists(cmdArgs, mergeID) {
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
			mr, resp, _ := gitlabClient.MergeRequests.UpdateMergeRequest(repo, manip.StringToInt(i2), l)
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

// mrCmd is merge request command
var mrCmd = &cobra.Command{
	Use:   "mr [subcommand] [flags]",
	Short: `Create, view and manage merge requests`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if create, _ := cmd.Flags().GetBool("create"); create {
			mrCreateCmd.Run(cmd, args)
			return
		}

		if len(args) == 0 || len(args) > 2 {
			cmd.Help()
			return
		}
	},
}

func init() {
	mrCmd.Flags().BoolP("list", "l", false, "List merge requests")
	mrCmd.Flags().BoolP("browse", "b", false, "View merge request <id> in a browser")
	mrCmd.Flags().StringP("close", "c", "", "Close merge request <id>")
	mrCmd.Flags().StringP("reopen", "o", "", "reopen a merge request <id>")
	mrCmd.Flags().StringP("delete", "d", "", "delete merge request <id>")
	mrCmd.Flags().StringP("subscribe", "s", "", "subscribe to a merge request <id>")
	mrCmd.Flags().StringP("unsubscribe", "u", "", "Unsubscribe to a merge request <id>")
	mrCmd.Flags().StringP("delete", "d", "", "Close merge request <id>")
	commands.RootCmd.AddCommand(mrCmd)
}

/*
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
		"revoke":      revokeMergeRequest,
		"update":      updateMergeRequest,
	}
	if _, ok := commandList[arrCmd[0]]; ok {
		commandList[arrCmd[0]](cmdArgs, arrCmd)
	} else {
		fmt.Println(arrCmd[0]+":", "Invalid Command")
	}
}
*/
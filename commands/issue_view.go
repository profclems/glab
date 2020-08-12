package commands

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/browser"
	"glab/internal/git"
	"glab/internal/manip"
	"glab/internal/utils"
	"log"
	"strings"
	"time"
)

var issueViewCmd = &cobra.Command{
	Use:     "view <id>",
	Short:   `Display the title, body, and other information about an issue.`,
	Long:    ``,
	Aliases: []string{"show"},
	Args:      cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) >1 {
			cmdErr(cmd, args)
			return
		}
		pid := manip.StringToInt(args[0])

		gitlabClient, repo := git.InitGitlabClient()

		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo = r
		}

		issue, _, err := gitlabClient.Issues.GetIssue(repo, pid)
		if err != nil {
			log.Fatal(err)
		}
		if lb, _ := cmd.Flags().GetBool("web"); lb { //open in browser if --web flag is specified
			a, err := browser.Command(issue.WebURL)
			if err != nil {
				er(err)
			}
			if err:= a.Run(); err != nil {
				er(err)
			}
			return
		}
		var issueState string
		if issue.State == "opened" {
			issueState = color.Green.Sprint(issue.State)
		} else {
			issueState = color.Red.Sprint(issue.State)
		}
		now := time.Now()
		ago := now.Sub(*issue.CreatedAt)
		color.Printf("\n%s <gray>#%d</>\n", issue.Title, issue.IID)
		color.Printf("(%s)<gray> • opened by %s (%s) %s</>\n", issueState,
			issue.Author.Username,
			issue.Author.Name,
			utils.PrettyTimeAgo(ago),
			)
		if issue.Description != "" {
			issue.Description, _ = utils.RenderMarkdown(issue.Description)
			fmt.Println(issue.Description)
		}
		color.Printf("\n<gray>%d upvotes • %d downvotes • %d comments</>\n\n",
			issue.Upvotes, issue.Downvotes, issue.UserNotesCount)
		var labels string
		for _, l := range issue.Labels {
			labels += " " + l + ","
		}
		labels = strings.Trim(labels, ", ")

		var assignees string
		for _, a := range issue.Assignees {
			assignees += " " + a.Username + "(" + a.Name + "),"
		}
		assignees = strings.Trim(assignees, ", ")
		table := uitable.New()
		table.MaxColWidth = 50
		table.Wrap = true
		table.AddRow("Project ID:", issue.ProjectID)
		table.AddRow("Labels:", prettifyNilEmptyValues(labels, "None"))
		table.AddRow("Milestone:", prettifyNilEmptyValues(issue.Milestone, "None"))
		table.AddRow("Assignees:", prettifyNilEmptyValues(assignees, "None"))
		table.AddRow("Due date:", prettifyNilEmptyValues(issue.DueDate, "None"))
		table.AddRow("Weight:", prettifyNilEmptyValues(issue.Weight, "None"))
		table.AddRow("Confidential:", prettifyNilEmptyValues(issue.Confidential, "None"))
		table.AddRow("Discussion Locked:", prettifyNilEmptyValues(issue.DiscussionLocked, "false"))
		table.AddRow("Subscribed:", prettifyNilEmptyValues(issue.Subscribed, "false"))

		if issue.State == "closed" {
			now := time.Now()
			ago := now.Sub(*issue.ClosedAt)
			table.AddRow("Closed By:",
				fmt.Sprintf("%s (%s) %s", issue.ClosedBy.Username, issue.ClosedBy.Name, utils.PrettyTimeAgo(ago)))
		}
		table.AddRow("Reference:", issue.References.Full)
		table.AddRow("Web URL:", issue.WebURL)
		fmt.Println(table)
		fmt.Println() // Empty Space

		if c, _ := cmd.Flags().GetBool("comments"); c { //open in browser if --web flag is specified
			l := &gitlab.ListIssueNotesOptions{}
			notes, _, err := gitlabClient.Notes.ListIssueNotes(repo, pid, l)
			if err != nil {
				er(err)
			}

			table := uitable.New()
			table.MaxColWidth = 100
			table.Wrap = true
			fmt.Println(heredoc.Doc(` 
			--------------------------------------------
			Comments / Notes
			--------------------------------------------
			`))
			if len(notes) > 0 {
				for _, note := range notes {
					if note.System {
						continue
					}
					//body, _ := utils.RenderMarkdown(note.Body)
					table.AddRow(note.Author.Username+":",
						fmt.Sprintf("%s\n%s",
							note.Body,
							color.Gray.Sprint(utils.TimeToPrettyTimeAgo(*note.CreatedAt)),
						),
					)
					table.AddRow("")
				}
				fmt.Println(table)
			} else {
				fmt.Println("There are no comments on this issue")
			}
		}
	},
}

func prettifyNilEmptyValues(value interface{}, defVal string) interface{}  {
	if value == nil || value == "" {
		return defVal
	}
	if value == false {
		return false
	}
	return value
}

func init() {
	issueViewCmd.Flags().StringP("repo", "r", "", "Select another repository using the OWNER/REPO format. Supports group namespaces")
	issueViewCmd.Flags().BoolP("comments", "c", false, "Show issue comments and activities")
	issueViewCmd.Flags().BoolP("web", "w", false, "Open issue in a browser. Uses default browser or browser specified in BROWSER variable")
	issueCmd.AddCommand(issueViewCmd)
}

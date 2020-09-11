package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"
	"github.com/profclems/glab/internal/utils"

	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var mrViewCmd = &cobra.Command{
	Use:     "view <id>",
	Short:   `Display the title, body, and other information about a merge request.`,
	Long:    ``,
	Aliases: []string{"show"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || len(args) > 1 {
			cmdErr(cmd, args)
			return nil
		}
		pid := manip.StringToInt(args[0])

		gitlabClient, repo := git.InitGitlabClient()

		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo, _ = fixRepoNamespace(r)
		}
		opts := &gitlab.GetMergeRequestsOptions{}
		opts.IncludeDivergedCommitsCount = gitlab.Bool(true)
		opts.RenderHTML = gitlab.Bool(true)
		opts.IncludeRebaseInProgress = gitlab.Bool(true)

		mr, _, err := gitlabClient.MergeRequests.GetMergeRequest(repo, pid, opts)
		if err != nil {
			return err
		}
		if lb, _ := cmd.Flags().GetBool("web"); lb { //open in browser if --web flag is specified
			fmt.Fprintf(cmd.ErrOrStderr(), "Opening %s in your browser.\n", utils.DisplayURL(mr.WebURL))
			return utils.OpenInBrowser(mr.WebURL)
		}
		showSystemLog, _ := cmd.Flags().GetBool("system-logs")
		var mrState string
		if mr.State == "opened" {
			mrState = color.Green.Sprint(mr.State)
		} else {
			mrState = color.Red.Sprint(mr.State)
		}
		now := time.Now()
		ago := now.Sub(*mr.CreatedAt)
		color.Printf("\n%s <gray>#%d</>\n", mr.Title, mr.IID)
		color.Printf("(%s)<gray> • opened by %s (%s) %s</>\n", mrState,
			mr.Author.Username,
			mr.Author.Name,
			utils.PrettyTimeAgo(ago),
		)
		if mr.Description != "" {
			mr.Description, _ = utils.RenderMarkdown(mr.Description)
			fmt.Println(mr.Description)
		}
		color.Printf("\n<gray>%d upvotes • %d downvotes • %d comments</>\n\n",
			mr.Upvotes, mr.Downvotes, mr.UserNotesCount)
		var labels string
		for _, l := range mr.Labels {
			labels += " " + l + ","
		}
		labels = strings.Trim(labels, ", ")

		var assignees string
		for _, a := range mr.Assignees {
			assignees += " " + a.Username + "(" + a.Name + "),"
		}
		assignees = strings.Trim(assignees, ", ")
		table := uitable.New()
		table.MaxColWidth = 70
		table.Wrap = true
		table.AddRow("Project ID:", mr.ProjectID)
		table.AddRow("Labels:", prettifyNilEmptyValues(labels, "None"))
		table.AddRow("Milestone:", prettifyNilEmptyValues(mr.Milestone, "None"))
		table.AddRow("Assignees:", prettifyNilEmptyValues(assignees, "None"))
		table.AddRow("Discussion Locked:", prettifyNilEmptyValues(mr.DiscussionLocked, "false"))
		table.AddRow("Subscribed:", prettifyNilEmptyValues(mr.Subscribed, "false"))

		if mr.State == "closed" {
			now := time.Now()
			ago := now.Sub(*mr.ClosedAt)
			table.AddRow("Closed By:",
				fmt.Sprintf("%s (%s) %s", mr.ClosedBy.Username, mr.ClosedBy.Name, utils.PrettyTimeAgo(ago)))
		}
		table.AddRow("Web URL:", mr.WebURL)
		fmt.Println(table)
		fmt.Println() // Empty Space

		if c, _ := cmd.Flags().GetBool("comments"); c {
			l := &gitlab.ListMergeRequestNotesOptions{}
			if p, _ := cmd.Flags().GetInt("page"); p != 0 {
				l.Page = p
			}
			if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
				l.PerPage = p
			}
			notes, _, err := gitlabClient.Notes.ListMergeRequestNotes(repo, pid, l)
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
					if note.System && !showSystemLog {
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
				fmt.Println("There are no comments on this mr")
			}
		}
		return nil
	},
}

func init() {
	mrViewCmd.Flags().BoolP("comments", "c", false, "Show mr comments and activities")
	mrViewCmd.Flags().BoolP("system-logs", "s", false, "Show system activities / logs")
	mrViewCmd.Flags().BoolP("web", "w", false, "Open mr in a browser. Uses default browser or browser specified in BROWSER variable")
	mrViewCmd.Flags().IntP("page", "p", 1, "Page number")
	mrViewCmd.Flags().IntP("per-page", "P", 20, "Number of items to list per page")
	mrCmd.AddCommand(mrViewCmd)
}

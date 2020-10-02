package view

import (
	"fmt"
	"strings"
	"time"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/MakeNowJust/heredoc"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdView(f *cmdutils.Factory) *cobra.Command {
	var issueViewCmd = &cobra.Command{
		Use:     "view <id>",
		Short:   `Display the title, body, and other information about an issue.`,
		Long:    ``,
		Aliases: []string{"show"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pid := utils.StringToInt(args[0])

			var err error
			out := utils.ColorableOut(cmd)
			if r, _ := cmd.Flags().GetString("repo"); r != "" {
				f, err = f.NewClient(r)
				if err != nil {
					return err
				}
			}
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}
			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			issue, err := api.GetIssue(apiClient, repo.FullName(), pid)
			if err != nil {
				return err
			}
			if lb, _ := cmd.Flags().GetBool("web"); lb { //open in browser if --web flag is specified
				fmt.Fprintf(out, "Opening %s in your browser.\n", utils.DisplayURL(issue.WebURL))
				cfg, _ := f.Config()
				browser, _ := cfg.Get(repo.RepoHost(), "browser")
				return utils.OpenInBrowser(issue.WebURL, browser)
			}
			var issueState string
			if issue.State == "opened" {
				issueState = utils.Green("open")
			} else {
				issueState = utils.Red(issue.State)
			}
			now := time.Now()
			ago := now.Sub(*issue.CreatedAt)

			var issuePrintDetails string

			issuePrintDetails += fmt.Sprintf("\n%s #%d\n", utils.Bold(issue.Title), issue.IID)
			issuePrintDetails += fmt.Sprintf("(%s) • ", issueState)
			issuePrintDetails += utils.Gray(fmt.Sprintf("opened by %s (%s) %s\n",
				issue.Author.Username,
				issue.Author.Name,
				utils.PrettyTimeAgo(ago),
			))
			if issue.Description != "" {
				cfg, _ := f.Config()
				glamourStyle, _ := cfg.Get(repo.RepoHost(), "glamour_style")
				issue.Description, _ = utils.RenderMarkdown(issue.Description, glamourStyle)
				issuePrintDetails += issue.Description
			}
			issuePrintDetails += utils.Gray(fmt.Sprintf("\n%d upvotes • %d downvotes • %d comments\n\n",
				issue.Upvotes, issue.Downvotes, issue.UserNotesCount))

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
			table.MaxColWidth = 70
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

			fmt.Fprintln(out, issuePrintDetails)
			fmt.Fprintln(out, table)
			fmt.Fprint(out, "") // Empty Space

			if c, _ := cmd.Flags().GetBool("comments"); c {
				showSystemLog, _ := cmd.Flags().GetBool("system-logs")
				var commentsPrintDetails string

				opts := &gitlab.ListIssueNotesOptions{}
				if p, _ := cmd.Flags().GetInt("page"); p != 0 {
					opts.Page = p
				}
				if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
					opts.PerPage = p
				}
				notes, err := api.ListIssueNotes(apiClient, repo.FullName(), pid, opts)
				if err != nil {
					return err
				}

				table := uitable.New()
				table.MaxColWidth = 100
				table.Wrap = true
				commentsPrintDetails += heredoc.Doc(` 
			--------------------------------------------
			Comments / Notes
			--------------------------------------------
			`)
				if len(notes) > 0 {
					for _, note := range notes {
						if note.System && !showSystemLog {
							continue
						}
						//body, _ := utils.RenderMarkdown(note.Body)
						table.AddRow(note.Author.Username+":",
							fmt.Sprintf("%s\n%s",
								note.Body,
								utils.Gray(utils.TimeToPrettyTimeAgo(*note.CreatedAt)),
							),
						)
						table.AddRow("")
					}
					fmt.Fprintln(out, commentsPrintDetails)
					fmt.Fprintln(out, table)
				} else {
					fmt.Fprintln(out, commentsPrintDetails, "There are no comments on this issue")
				}
			}
			return nil
		},
	}

	issueViewCmd.Flags().BoolP("comments", "c", false, "Show issue comments and activities")
	issueViewCmd.Flags().BoolP("system-logs", "s", false, "Show system activities / logs")
	issueViewCmd.Flags().BoolP("web", "w", false, "Open issue in a browser. Uses default browser or browser specified in BROWSER variable")
	issueViewCmd.Flags().IntP("page", "p", 1, "Page number")
	issueViewCmd.Flags().IntP("per-page", "P", 20, "Number of items to list per page")

	return issueViewCmd
}

func prettifyNilEmptyValues(value interface{}, defVal string) interface{} {
	if value == nil || value == "" {
		return defVal
	}
	if value == false {
		return false
	}
	return value
}

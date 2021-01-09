package view

import (
	"fmt"
	"io"
	"strings"

	"github.com/profclems/glab/commands/issue/issueutils"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var (
	showSystemLogs bool
	showComments   bool
	limit          int
	pageNumber     int
	cfg            config.Config
	glamourStyle   string
	notes          []*gitlab.Note
)

func NewCmdView(f *cmdutils.Factory) *cobra.Command {
	var issueViewCmd = &cobra.Command{
		Use:     "view <id>",
		Short:   `Display the title, body, and other information about an issue.`,
		Long:    ``,
		Aliases: []string{"show"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			var err error
			out := f.IO.StdOut

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}
			cfg, _ = f.Config()

			issue, baseRepo, err := issueutils.IssueFromArg(apiClient, f.BaseRepo, args[0])
			if err != nil {
				return err
			}

			//open in browser if --web flag is specified
			if isWeb, _ := cmd.Flags().GetBool("web"); isWeb {
				if f.IO.IsaTTY && f.IO.IsErrTTY {
					fmt.Fprintf(out, "Opening %s in your browser.\n", utils.DisplayURL(issue.WebURL))
				}

				browser, _ := cfg.Get(baseRepo.RepoHost(), "browser")
				return utils.OpenInBrowser(issue.WebURL, browser)
			}

			glamourStyle, _ = cfg.Get(baseRepo.RepoHost(), "glamour_style")

			if showComments {
				l := &gitlab.ListIssueNotesOptions{
					Sort: gitlab.String("asc"),
				}
				if pageNumber != 0 {
					l.Page = pageNumber
				}
				if limit != 0 {
					l.PerPage = limit
				}
				notes, err = api.ListIssueNotes(apiClient, baseRepo.FullName(), issue.IID, l)
				if err != nil {
					return err
				}
			}

			err = f.IO.StartPager()
			if err != nil {
				return err
			}
			defer f.IO.StopPager()
			if f.IO.IsErrTTY && f.IO.IsaTTY {
				return printTTYIssuePreview(f.IO.StdOut, issue)
			}
			return printRawIssuePreview(f.IO.StdOut, issue)
		},
	}

	issueViewCmd.Flags().BoolVarP(&showComments, "comments", "c", false, "Show mr comments and activities")
	issueViewCmd.Flags().BoolVarP(&showSystemLogs, "system-logs", "s", false, "Show system activities / logs")
	issueViewCmd.Flags().BoolP("web", "w", false, "Open mr in a browser. Uses default browser or browser specified in BROWSER variable")
	issueViewCmd.Flags().IntVarP(&pageNumber, "page", "p", 1, "Page number")
	issueViewCmd.Flags().IntVarP(&limit, "per-page", "P", 20, "Number of items to list per page")

	return issueViewCmd
}

func labelsList(issue *gitlab.Issue) string {
	var labels string
	for _, l := range issue.Labels {
		labels += " " + l + ","
	}
	return strings.Trim(labels, ", ")
}

func assigneesList(issue *gitlab.Issue) string {
	var assignees string
	for _, a := range issue.Assignees {
		assignees += " " + a.Username + ","
	}
	return strings.Trim(assignees, ", ")
}

func issueState(issue *gitlab.Issue) (state string) {
	if issue.State == "opened" {
		state = utils.Green("open")
	} else if issue.State == "locked" {
		state = utils.Blue(issue.State)
	} else {
		state = utils.Red(issue.State)
	}

	return
}

func printTTYIssuePreview(out io.Writer, issue *gitlab.Issue) error {
	issueTimeAgo := utils.TimeToPrettyTimeAgo(*issue.CreatedAt)
	// Header
	fmt.Fprint(out, issueState(issue))
	fmt.Fprintf(out, utils.Gray(" • opened by %s %s\n"), issue.Author.Username, issueTimeAgo)
	fmt.Fprint(out, issue.Title)
	fmt.Fprintf(out, utils.Gray(" #%d"), issue.IID)
	fmt.Fprintln(out)

	// Description
	if issue.Description != "" {
		issue.Description, _ = utils.RenderMarkdown(issue.Description, glamourStyle)
		fmt.Fprintln(out, issue.Description)
	}

	fmt.Fprintf(out, utils.Gray("\n%d upvotes • %d downvotes • %d comments\n"), issue.Upvotes, issue.Downvotes, issue.UserNotesCount)

	// Meta information
	if labels := labelsList(issue); labels != "" {
		fmt.Fprint(out, utils.Bold("Labels: "))
		fmt.Fprintln(out, labels)
	}
	if assignees := assigneesList(issue); assignees != "" {
		fmt.Fprint(out, utils.Bold("Assignees: "))
		fmt.Fprintln(out, assignees)
	}
	if issue.Milestone != nil {
		fmt.Fprint(out, utils.Bold("Milestone: "))
		fmt.Fprintln(out, issue.Milestone.Title)
	}
	if issue.State == "closed" {
		fmt.Fprintf(out, "Closed By: %s %s\n", issue.ClosedBy.Username, issueTimeAgo)
	}

	// Comments
	if showComments {
		fmt.Fprintln(out, heredoc.Doc(` 
			--------------------------------------------
			Comments / Notes
			--------------------------------------------
			`))
		if len(notes) > 0 {
			for _, note := range notes {
				if note.System && !showSystemLogs {
					continue
				}
				createdAt := utils.TimeToPrettyTimeAgo(*note.CreatedAt)
				fmt.Fprint(out, note.Author.Username)
				if note.System {
					fmt.Fprintf(out, " %s ", note.Body)
					fmt.Fprintln(out, utils.Gray(createdAt))
				} else {
					body, _ := utils.RenderMarkdown(note.Body, glamourStyle)
					fmt.Fprint(out, " commented ")
					fmt.Fprintf(out, utils.Gray("%s\n"), createdAt)
					fmt.Fprintln(out, utils.Indent(body, " "))
				}
				fmt.Fprintln(out)
			}
		} else {
			fmt.Fprintln(out, "There are no comments on this issue")
		}
	}

	fmt.Fprintln(out)
	fmt.Fprintf(out, utils.Gray("View this issue on GitLab: %s\n"), issue.WebURL)

	return nil
}

func printRawIssuePreview(out io.Writer, issue *gitlab.Issue) error {
	assignees := assigneesList(issue)
	labels := labelsList(issue)

	fmt.Fprintf(out, "title:\t%s\n", issue.Title)
	fmt.Fprintf(out, "state:\t%s\n", issueState(issue))
	fmt.Fprintf(out, "author:\t%s\n", issue.Author.Username)
	fmt.Fprintf(out, "labels:\t%s\n", labels)
	fmt.Fprintf(out, "comments:\t%d\n", issue.UserNotesCount)
	fmt.Fprintf(out, "assignees:\t%s\n", assignees)
	if issue.Milestone != nil {
		fmt.Fprintf(out, "milestone:\t%s\n", issue.Milestone.Title)
	}

	fmt.Fprintln(out, "--")
	fmt.Fprintln(out, issue.Description)
	return nil
}

package view

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/commands/issue/issueutils"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/utils"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type ViewOpts struct {
	ShowComments   bool
	ShowSystemLogs bool
	OpenInBrowser  bool
	Web            bool

	CommentPageNumber int
	CommentLimit      int

	Notes []*gitlab.Note
	Issue *gitlab.Issue

	IO *iostreams.IOStreams
}

func NewCmdView(f *cmdutils.Factory) *cobra.Command {
	opts := &ViewOpts{
		IO: f.IO,
	}
	var issueViewCmd = &cobra.Command{
		Use:     "view <id>",
		Short:   `Display the title, body, and other information about an issue.`,
		Long:    ``,
		Aliases: []string{"show"},
		Example: heredoc.Doc(`
			$ glab issue view 123
			$ glab issue show 123
			$ glab issue view --web 123
			$ glab issue view --comments 123
			$ glab issue view https://gitlab.com/profclems/glab/-/issues/123
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}
			cfg, _ := f.Config()

			issue, baseRepo, err := issueutils.IssueFromArg(apiClient, f.BaseRepo, args[0])
			if err != nil {
				return err
			}

			opts.Issue = issue

			//open in browser if --web flag is specified
			if opts.Web {
				if f.IO.IsaTTY && f.IO.IsErrTTY {
					fmt.Fprintf(opts.IO.StdErr, "Opening %s in your browser.\n", utils.DisplayURL(opts.Issue.WebURL))
				}

				browser, _ := cfg.Get(baseRepo.RepoHost(), "browser")
				return utils.OpenInBrowser(opts.Issue.WebURL, browser)
			}

			if opts.ShowComments {
				l := &gitlab.ListIssueNotesOptions{
					Sort: gitlab.String("asc"),
				}
				if opts.CommentPageNumber != 0 {
					l.Page = opts.CommentPageNumber
				}
				if opts.CommentLimit != 0 {
					l.PerPage = opts.CommentLimit
				}
				opts.Notes, err = api.ListIssueNotes(apiClient, baseRepo.FullName(), opts.Issue.IID, l)
				if err != nil {
					return err
				}
			}

			glamourStyle, _ := cfg.Get(baseRepo.RepoHost(), "glamour_style")
			f.IO.ResolveBackgroundColor(glamourStyle)
			err = f.IO.StartPager()
			if err != nil {
				return err
			}
			defer f.IO.StopPager()
			if f.IO.IsErrTTY && f.IO.IsaTTY {
				return printTTYIssuePreview(opts)
			}
			return printRawIssuePreview(opts)
		},
	}

	issueViewCmd.Flags().BoolVarP(&opts.ShowComments, "comments", "c", false, "Show mr comments and activities")
	issueViewCmd.Flags().BoolVarP(&opts.ShowSystemLogs, "system-logs", "s", false, "Show system activities / logs")
	issueViewCmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open mr in a browser. Uses default browser or browser specified in BROWSER variable")
	issueViewCmd.Flags().IntVarP(&opts.CommentPageNumber, "page", "p", 1, "Page number")
	issueViewCmd.Flags().IntVarP(&opts.CommentLimit, "per-page", "P", 20, "Number of items to list per page")

	return issueViewCmd
}

func labelsList(opts *ViewOpts) string {
	var labels string
	for _, l := range opts.Issue.Labels {
		labels += " " + l + ","
	}
	return strings.Trim(labels, ", ")
}

func assigneesList(opts *ViewOpts) string {
	var assignees string
	for _, a := range opts.Issue.Assignees {
		assignees += " " + a.Username + ","
	}
	return strings.Trim(assignees, ", ")
}

func issueState(opts *ViewOpts, c *iostreams.ColorPalette) (state string) {
	if opts.Issue.State == "opened" {
		state = c.Green("open")
	} else if opts.Issue.State == "locked" {
		state = c.Blue(opts.Issue.State)
	} else {
		state = c.Red(opts.Issue.State)
	}

	return
}

func printTTYIssuePreview(opts *ViewOpts) error {
	c := opts.IO.Color()
	issueTimeAgo := utils.TimeToPrettyTimeAgo(*opts.Issue.CreatedAt)
	// Header
	fmt.Fprint(opts.IO.StdOut, issueState(opts, c))
	fmt.Fprintf(opts.IO.StdOut, c.Gray(" • opened by %s %s\n"), opts.Issue.Author.Username, issueTimeAgo)
	fmt.Fprint(opts.IO.StdOut, c.Bold(opts.Issue.Title))
	fmt.Fprintf(opts.IO.StdOut, c.Gray(" #%d"), opts.Issue.IID)
	fmt.Fprintln(opts.IO.StdOut)

	// Description
	if opts.Issue.Description != "" {
		opts.Issue.Description, _ = utils.RenderMarkdown(opts.Issue.Description, opts.IO.BackgroundColor())
		fmt.Fprintln(opts.IO.StdOut, opts.Issue.Description)
	}

	fmt.Fprintf(opts.IO.StdOut, c.Gray("\n%d upvotes • %d downvotes • %d comments\n"), opts.Issue.Upvotes, opts.Issue.Downvotes, opts.Issue.UserNotesCount)

	// Meta information
	if labels := labelsList(opts); labels != "" {
		fmt.Fprint(opts.IO.StdOut, c.Bold("Labels: "))
		fmt.Fprintln(opts.IO.StdOut, labels)
	}
	if assignees := assigneesList(opts); assignees != "" {
		fmt.Fprint(opts.IO.StdOut, c.Bold("Assignees: "))
		fmt.Fprintln(opts.IO.StdOut, assignees)
	}
	if opts.Issue.Milestone != nil {
		fmt.Fprint(opts.IO.StdOut, c.Bold("Milestone: "))
		fmt.Fprintln(opts.IO.StdOut, opts.Issue.Milestone.Title)
	}
	if opts.Issue.State == "closed" {
		fmt.Fprintf(opts.IO.StdOut, "Closed By: %s %s\n", opts.Issue.ClosedBy.Username, issueTimeAgo)
	}

	// Comments
	if opts.ShowComments {
		fmt.Fprintln(opts.IO.StdOut, heredoc.Doc(` 
			--------------------------------------------
			Comments / Notes
			--------------------------------------------
			`))
		if len(opts.Notes) > 0 {
			for _, note := range opts.Notes {
				if note.System && !opts.ShowSystemLogs {
					continue
				}
				createdAt := utils.TimeToPrettyTimeAgo(*note.CreatedAt)
				fmt.Fprint(opts.IO.StdOut, note.Author.Username)
				if note.System {
					fmt.Fprintf(opts.IO.StdOut, " %s ", note.Body)
					fmt.Fprintln(opts.IO.StdOut, c.Gray(createdAt))
				} else {
					body, _ := utils.RenderMarkdown(note.Body, opts.IO.BackgroundColor())
					fmt.Fprint(opts.IO.StdOut, " commented ")
					fmt.Fprintf(opts.IO.StdOut, c.Gray("%s\n"), createdAt)
					fmt.Fprintln(opts.IO.StdOut, utils.Indent(body, " "))
				}
				fmt.Fprintln(opts.IO.StdOut)
			}
		} else {
			fmt.Fprintln(opts.IO.StdOut, "There are no comments on this issue")
		}
	}

	fmt.Fprintln(opts.IO.StdOut)
	fmt.Fprintf(opts.IO.StdOut, c.Gray("View this issue on GitLab: %s\n"), opts.Issue.WebURL)

	return nil
}

func printRawIssuePreview(opts *ViewOpts) error {
	assignees := assigneesList(opts)
	labels := labelsList(opts)

	fmt.Fprintf(opts.IO.StdOut, "title:\t%s\n", opts.Issue.Title)
	fmt.Fprintf(opts.IO.StdOut, "state:\t%s\n", issueState(opts, opts.IO.Color()))
	fmt.Fprintf(opts.IO.StdOut, "author:\t%s\n", opts.Issue.Author.Username)
	fmt.Fprintf(opts.IO.StdOut, "labels:\t%s\n", labels)
	fmt.Fprintf(opts.IO.StdOut, "comments:\t%d\n", opts.Issue.UserNotesCount)
	fmt.Fprintf(opts.IO.StdOut, "assignees:\t%s\n", assignees)
	if opts.Issue.Milestone != nil {
		fmt.Fprintf(opts.IO.StdOut, "milestone:\t%s\n", opts.Issue.Milestone.Title)
	}

	fmt.Fprintln(opts.IO.StdOut, "--")
	fmt.Fprintln(opts.IO.StdOut, opts.Issue.Description)
	return nil
}

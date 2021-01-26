package view

import (
	"fmt"
	"io"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type ViewOpts struct {
	ShowComments   bool
	ShowSystemLogs bool
	OpenInBrowser  bool

	CommentPageNumber int
	CommentLimit      int
}

func NewCmdView(f *cmdutils.Factory) *cobra.Command {
	opts := ViewOpts{}
	var mrViewCmd = &cobra.Command{
		Use:     "view {<id> | <branch>}",
		Short:   `Display the title, body, and other information about a merge request.`,
		Long:    ``,
		Aliases: []string{"show"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mr, baseRepo, err := mrutils.MRFromArgsWithOpts(f, args, &gitlab.GetMergeRequestsOptions{
				IncludeDivergedCommitsCount: gitlab.Bool(true),
				RenderHTML:                  gitlab.Bool(true),
				IncludeRebaseInProgress:     gitlab.Bool(true),
			})
			if err != nil {
				return err
			}
			cfg, _ := f.Config()

			if opts.OpenInBrowser { //open in browser if --web flag is specified
				if f.IO.IsOutputTTY() {
					fmt.Fprintf(f.IO.StdErr, "Opening %s in your browser.\n", utils.DisplayURL(mr.WebURL))
				}

				browser, _ := cfg.Get(baseRepo.RepoHost(), "browser")
				return utils.OpenInBrowser(mr.WebURL, browser)
			}

			notes := []*gitlab.Note{}

			if opts.ShowComments {
				l := &gitlab.ListMergeRequestNotesOptions{
					Sort: gitlab.String("asc"),
				}
				l.Page = opts.CommentPageNumber
				l.PerPage = opts.CommentLimit

				notes, err = api.ListMRNotes(apiClient, baseRepo.FullName(), mr.IID, l)
				if err != nil {
					return err
				}
			}

			if err := f.IO.StartPager(); err != nil {
				return err
			}
			defer f.IO.StopPager()

			if f.IO.IsOutputTTY() {
				glamourStyle, _ := cfg.Get(baseRepo.RepoHost(), "glamour_style")
				return printTTYMRPreview(f.IO.StdOut, mr, notes, opts, glamourStyle)
			}
			return printRawMRPreview(f.IO.StdOut, mr)
		},
	}

	mrViewCmd.Flags().BoolVarP(&opts.ShowComments, "comments", "c", false, "Show mr comments and activities")
	mrViewCmd.Flags().BoolVarP(&opts.ShowSystemLogs, "system-logs", "s", false, "Show system activities / logs")
	mrViewCmd.Flags().BoolVarP(&opts.OpenInBrowser, "web", "w", false, "Open mr in a browser. Uses default browser or browser specified in BROWSER variable")
	mrViewCmd.Flags().IntVarP(&opts.CommentPageNumber, "page", "p", 0, "Page number")
	mrViewCmd.Flags().IntVarP(&opts.CommentLimit, "per-page", "P", 20, "Number of items to list per page")

	return mrViewCmd
}

func labelsList(mr *gitlab.MergeRequest) string {
	var labels string
	for _, l := range mr.Labels {
		labels += " " + l + ","
	}
	return strings.Trim(labels, ", ")
}

func assigneesList(mr *gitlab.MergeRequest) string {
	var assignees string
	for _, a := range mr.Assignees {
		assignees += " " + a.Username + ","
	}
	return strings.Trim(assignees, ", ")
}

func mrState(mr *gitlab.MergeRequest) (mrState string) {
	if mr.State == "opened" {
		mrState = utils.Green("open")
	} else if mr.State == "merged" {
		mrState = utils.Blue(mr.State)
	} else {
		mrState = utils.Red(mr.State)
	}

	return mrState
}

func printTTYMRPreview(out io.Writer, mr *gitlab.MergeRequest, notes []*gitlab.Note, opts ViewOpts, glamourStyle string) error {
	mrTimeAgo := utils.TimeToPrettyTimeAgo(*mr.CreatedAt)
	// Header
	fmt.Fprint(out, mrState(mr))
	fmt.Fprintf(out, utils.Gray(" • opened by %s %s\n"), mr.Author.Username, mrTimeAgo)
	fmt.Fprint(out, mr.Title)
	fmt.Fprintf(out, utils.Gray(" !%d"), mr.IID)
	fmt.Fprintln(out)

	// Description
	if mr.Description != "" {
		mr.Description, _ = utils.RenderMarkdown(mr.Description, glamourStyle)
		fmt.Fprintln(out, mr.Description)
	}

	fmt.Fprintf(out, utils.Gray("\n%d upvotes • %d downvotes • %d comments\n"), mr.Upvotes, mr.Downvotes, mr.UserNotesCount)

	// Meta information
	if labels := labelsList(mr); labels != "" {
		fmt.Fprint(out, utils.Bold("Labels: "))
		fmt.Fprintln(out, labels)
	}
	if assignees := assigneesList(mr); assignees != "" {
		fmt.Fprint(out, utils.Bold("Assignees: "))
		fmt.Fprintln(out, assignees)
	}
	if mr.Milestone != nil {
		fmt.Fprint(out, utils.Bold("Milestone: "))
		fmt.Fprintln(out, mr.Milestone.Title)
	}
	if mr.State == "closed" {
		fmt.Fprintf(out, "Closed By: %s %s\n", mr.ClosedBy.Username, mrTimeAgo)
	}
	if mr.Pipeline != nil {
		fmt.Fprint(out, utils.Bold("Pipeline Status: "))
		var status string
		switch s := mr.Pipeline.Status; s {
		case "failed":
			status = utils.Red(s)
		case "success":
			status = utils.Green(s)
		default:
			status = utils.Gray(s)
		}
		fmt.Fprintf(out, "%s (View pipeline with `%s`)\n", status, utils.Bold("glab ci view "+mr.SourceBranch))

		if mr.MergeWhenPipelineSucceeds && mr.Pipeline.Status != "success" {
			fmt.Fprintf(out, "%s Requires pipeline to succeed before merging\n", utils.WarnIcon())
		}
	}
	fmt.Fprintf(out, "%s This merge request has %s changes\n", utils.GreenCheck(), utils.Yellow(mr.ChangesCount))
	if mr.State == "merged" && mr.MergedBy != nil {
		fmt.Fprintf(out, "%s The changes were merged into %s by %s %s\n", utils.GreenCheck(), mr.TargetBranch, mr.MergedBy.Name, utils.TimeToPrettyTimeAgo(*mr.MergedAt))
	}

	if mr.HasConflicts {
		fmt.Fprintf(out, utils.Red("%s This branch has conflicts that must be resolved\n"), utils.FailedIcon())
	}

	// Comments
	if opts.ShowComments {
		fmt.Fprintln(out, heredoc.Doc(` 
			--------------------------------------------
			Comments / Notes
			--------------------------------------------
			`))
		if len(notes) > 0 {
			for _, note := range notes {
				if note.System && !opts.ShowSystemLogs {
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
			fmt.Fprintln(out, "There are no comments on this merge request")
		}
	}

	fmt.Fprintln(out)
	fmt.Fprintf(out, utils.Gray("View this merge request on GitLab: %s\n"), mr.WebURL)

	return nil
}

func printRawMRPreview(out io.Writer, mr *gitlab.MergeRequest) error {
	assignees := assigneesList(mr)
	labels := labelsList(mr)

	fmt.Fprintf(out, "title:\t%s\n", mr.Title)
	fmt.Fprintf(out, "state:\t%s\n", mrState(mr))
	fmt.Fprintf(out, "author:\t%s\n", mr.Author.Username)
	fmt.Fprintf(out, "labels:\t%s\n", labels)
	fmt.Fprintf(out, "assignees:\t%s\n", assignees)
	fmt.Fprintf(out, "comments:\t%d\n", mr.UserNotesCount)
	if mr.Milestone != nil {
		fmt.Fprintf(out, "milestone:\t%s\n", mr.Milestone.Title)
	}
	fmt.Fprintf(out, "number:\t%d\n", mr.IID)
	fmt.Fprintf(out, "url:\t%s\n", mr.WebURL)

	fmt.Fprintln(out, "--")
	fmt.Fprintln(out, mr.Description)

	return nil
}

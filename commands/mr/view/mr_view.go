package view

import (
	"fmt"
	"strings"
	"time"

	"github.com/profclems/glab/commands/mr/mrutils"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/MakeNowJust/heredoc"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdView(f *cmdutils.Factory) *cobra.Command {
	var mrViewCmd = &cobra.Command{
		Use:     "view {<id> | <branch>}",
		Short:   `Display the title, body, and other information about a merge request.`,
		Long:    ``,
		Aliases: []string{"show"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := utils.ColorableOut(cmd)

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mr, repo, err := mrutils.MRFromArgs(f, args)
			if err != nil {
				return err
			}

			opts := &gitlab.GetMergeRequestsOptions{}
			opts.IncludeDivergedCommitsCount = gitlab.Bool(true)
			opts.RenderHTML = gitlab.Bool(true)
			opts.IncludeRebaseInProgress = gitlab.Bool(true)

			mr, err = api.GetMR(apiClient, repo.FullName(), mr.IID, opts)
			if err != nil {
				return err
			}

			if lb, _ := cmd.Flags().GetBool("web"); lb { //open in browser if --web flag is specified
				fmt.Fprintf(cmd.ErrOrStderr(), "Opening %s in your browser.\n", utils.DisplayURL(mr.WebURL))
				cfg, _ := f.Config()
				browser, _ := cfg.Get(repo.RepoHost(), "browser")
				return utils.OpenInBrowser(mr.WebURL, browser)
			}

			showSystemLog, _ := cmd.Flags().GetBool("system-logs")

			var mrState string
			if mr.State == "opened" {
				mrState = utils.Green("open")
			} else if mr.State == "merged" {
				mrState = utils.Blue(mr.State)
			} else {
				mrState = utils.Red(mr.State)
			}
			now := time.Now()
			ago := now.Sub(*mr.CreatedAt)

			mrPrintDetails := mrState
			mrPrintDetails += utils.Gray(fmt.Sprintf(" • opened by %s (%s) %s\n", mr.Author.Username, mr.Author.Name, utils.PrettyTimeAgo(ago)))
			mrPrintDetails += mr.Title
			mrPrintDetails += fmt.Sprintf("!%d", mr.IID)
			mrPrintDetails += "\n"

			if mr.Description != "" {
				cfg, _ := f.Config()
				glamourStyle, _ := cfg.Get(repo.RepoHost(), "glamour_style")
				mr.Description, _ = utils.RenderMarkdown(mr.Description, glamourStyle)
				mrPrintDetails += mr.Description
			}

			mrPrintDetails += utils.Gray(fmt.Sprintf("\n%d upvotes • %d downvotes • %d comments\n\n",
				mr.Upvotes, mr.Downvotes, mr.UserNotesCount))

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
			mrPrintDetails += heredoc.Docf(`
			Labels: %v
			Assignees: %v
			Milestone: %v
			`, prettifyNilEmptyValues(labels, "None"),
				prettifyNilEmptyValues(assignees, "None"),
				prettifyNilEmptyValues(mr.Milestone, "None"))

			if mr.State == "closed" {
				now := time.Now()
				ago := now.Sub(*mr.ClosedAt)
				mrPrintDetails += fmt.Sprintf("Closed By: %s (%s) %s", mr.ClosedBy.Username, mr.ClosedBy.Name, utils.PrettyTimeAgo(ago))
			}

			mrPrintDetails += utils.Gray(fmt.Sprintf("\nView this pull request on GitLab: %s", mr.WebURL))
			mrPrintDetails += "\n"

			fmt.Fprintln(out, mrPrintDetails)

			if c, _ := cmd.Flags().GetBool("comments"); c {
				l := &gitlab.ListMergeRequestNotesOptions{}
				if p, _ := cmd.Flags().GetInt("page"); p != 0 {
					l.Page = p
				}
				if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
					l.PerPage = p
				}
				notes, err := api.ListMRNotes(apiClient, repo.FullName(), mr.IID, l)
				if err != nil {
					return err
				}

				table := uitable.New()
				table.MaxColWidth = 100
				table.Wrap = true
				fmt.Fprintln(out, heredoc.Doc(` 
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
								utils.Gray(utils.TimeToPrettyTimeAgo(*note.CreatedAt)),
							),
						)
						table.AddRow("")
					}
					fmt.Fprintln(out, table)
				} else {
					fmt.Fprintln(out, "There are no comments on this mr")
				}
			}

			return nil
		},
	}

	mrViewCmd.Flags().BoolP("comments", "c", false, "Show mr comments and activities")
	mrViewCmd.Flags().BoolP("system-logs", "s", false, "Show system activities / logs")
	mrViewCmd.Flags().BoolP("web", "w", false, "Open mr in a browser. Uses default browser or browser specified in BROWSER variable")
	mrViewCmd.Flags().IntP("page", "p", 1, "Page number")
	mrViewCmd.Flags().IntP("per-page", "P", 20, "Number of items to list per page")

	return mrViewCmd
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

package view

import (
	"fmt"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/api"
	"strings"
	"time"

	"github.com/profclems/glab/internal/utils"

	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdView(f *cmdutils.Factory) *cobra.Command {
	var mrViewCmd = &cobra.Command{
		Use:     "view <id>",
		Short:   `Display the title, body, and other information about a merge request.`,
		Long:    ``,
		Aliases: []string{"show"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			pid := utils.StringToInt(args[0])

			opts := &gitlab.GetMergeRequestsOptions{}
			opts.IncludeDivergedCommitsCount = gitlab.Bool(true)
			opts.RenderHTML = gitlab.Bool(true)
			opts.IncludeRebaseInProgress = gitlab.Bool(true)

			mr, err := api.GetMR(apiClient, repo.FullName(), pid, opts)
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
				mrState = utils.Green(mr.State)
			} else if mr.State == "merged" {
				mrState = utils.Blue(mr.State)
			} else {
				mrState = utils.Red(mr.State)
			}
			now := time.Now()
			ago := now.Sub(*mr.CreatedAt)

			mrPrintDetails := "\n" + mr.Title
			mrPrintDetails += fmt.Sprintf("#%d", mr.IID)
			mrPrintDetails += fmt.Sprintf("(%s)", mrState)
			mrPrintDetails += utils.Gray(fmt.Sprintf(" • opened by %s (%s) %s\n", mr.Author.Username,
				mr.Author.Name,
				utils.PrettyTimeAgo(ago)))

			if mr.Description != "" {
				cfg, _ := f.Config()
				glamourStyle, _ := cfg.Get(repo.RepoHost(), "glamour_style")
				mr.Description, _ = utils.RenderMarkdown(mr.Description, glamourStyle)
				mrPrintDetails += mr.Description
			}

			mrPrintDetails += utils.Gray(fmt.Sprintf("\n%d upvotes • %d downvotes • %d comments\n\n",
				mr.Upvotes, mr.Downvotes, mr.UserNotesCount))

			fmt.Fprintln(out, mrPrintDetails)

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
			fmt.Fprintln(out, table)
			fmt.Fprint(out, "\n") // Empty Line

			if c, _ := cmd.Flags().GetBool("comments"); c {
				l := &gitlab.ListMergeRequestNotesOptions{}
				if p, _ := cmd.Flags().GetInt("page"); p != 0 {
					l.Page = p
				}
				if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
					l.PerPage = p
				}
				notes, err := api.ListMRNotes(apiClient, repo.FullName(), pid, l)
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
								color.Gray.Sprint(utils.TimeToPrettyTimeAgo(*note.CreatedAt)),
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

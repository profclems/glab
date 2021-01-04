package list

import (
	"fmt"

	"github.com/profclems/glab/internal/glrepo"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type ListOptions struct {
	// metadata
	Assignee  string
	Labels    string
	Milestone string
	Mine      bool

	// issue states
	State        string
	Closed       bool
	Opened       bool
	All          bool
	Confidential bool

	// Pagination
	Page    int
	PerPage int

	// display opts
	ListType       string
	TitleQualifier string

	IO         *utils.IOStreams
	BaseRepo   func() (glrepo.Interface, error)
	HTTPClient func() (*gitlab.Client, error)
}

func NewCmdList(f *cmdutils.Factory, runE func(opts *ListOptions) error) *cobra.Command {
	var opts = &ListOptions{
		IO: f.IO,
	}

	var issueListCmd = &cobra.Command{
		Use:     "list [flags]",
		Short:   `List project issues`,
		Long:    ``,
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.BaseRepo = f.BaseRepo
			opts.HTTPClient = f.HttpClient

			if opts.All {
				opts.State = "all"
			} else if opts.Closed {
				opts.State = "closed"
				opts.TitleQualifier = "closed"
			} else {
				opts.State = "opened"
				opts.TitleQualifier = "open"
			}

			if runE != nil {
				return runE(opts)
			}

			return listRun(opts)
		},
	}
	issueListCmd.Flags().StringVarP(&opts.Assignee, "assignee", "", "", "Filter issue by assignee <username>")
	issueListCmd.Flags().StringVarP(&opts.Labels, "label", "l", "", "Filter issue by label <name>")
	issueListCmd.Flags().StringVarP(&opts.Milestone, "milestone", "", "", "Filter issue by milestone <id>")
	issueListCmd.Flags().BoolVarP(&opts.Mine, "mine", "", false, "Filter only issues issues assigned to me")
	issueListCmd.Flags().BoolVarP(&opts.All, "all", "a", false, "Get all issues")
	issueListCmd.Flags().BoolVarP(&opts.Opened, "closed", "c", false, "Get only closed issues")
	issueListCmd.Flags().BoolVarP(&opts.Opened, "opened", "o", false, "Get only opened issues")
	issueListCmd.Flags().BoolVarP(&opts.Confidential, "confidential", "", false, "Filter by confidential issues")
	issueListCmd.Flags().IntVarP(&opts.Page, "page", "p", 1, "Page number")
	issueListCmd.Flags().IntVarP(&opts.PerPage, "per-page", "P", 30, "Number of items to list per page. (default 30)")

	return issueListCmd
}

func listRun(opts *ListOptions) error {
	apiClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}

	repo, err := opts.BaseRepo()
	if err != nil {
		return err
	}

	listOpts := &gitlab.ListProjectIssuesOptions{
		State: gitlab.String(opts.State),
	}
	listOpts.Page = 1
	listOpts.PerPage = 30

	if opts.Assignee != "" {
		listOpts.AssigneeUsername = gitlab.String(opts.Assignee)
	}
	if opts.Labels != "" {
		label := gitlab.Labels{
			opts.Labels,
		}
		listOpts.Labels = label
		opts.ListType = "search"
	}
	if opts.Milestone != "" {
		listOpts.Milestone = gitlab.String(opts.Milestone)
		opts.ListType = "search"
	}
	if opts.Confidential {
		listOpts.Confidential = gitlab.Bool(opts.Confidential)
		opts.ListType = "search"
	}
	if opts.Page != 0 {
		listOpts.Page = opts.Page
		opts.ListType = "search"
	}
	if opts.PerPage != 0 {
		listOpts.PerPage = opts.PerPage
		opts.ListType = "search"
	}

	if opts.Mine {
		u, err := api.CurrentUser(nil)
		if err != nil {
			return err
		}
		listOpts.AssigneeUsername = gitlab.String(u.Username)
		opts.ListType = "search"
	}
	issues, err := api.ListIssues(apiClient, repo.FullName(), listOpts)
	if err != nil {
		return err
	}

	title := utils.NewListTitle(opts.TitleQualifier + " issue")
	title.RepoName = repo.FullName()
	title.Page = listOpts.Page
	title.ListActionType = opts.ListType
	title.CurrentPageTotal = len(issues)

	if opts.IO.StartPager() != nil {
		return fmt.Errorf("failed to start pager: %q", err)
	}
	defer opts.IO.StopPager()

	fmt.Fprintf(opts.IO.StdOut, "%s\n%s\n", title.Describe(), issueutils.DisplayIssueList(issues, repo.FullName()))
	return nil
}

package list

import (
	"errors"
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
	Author    string
	NotAuthor []string
	Labels    []string
	NotLabels []string
	Milestone string
	Mine      bool
	Search    string

	// issue states
	State        string
	Closed       bool
	Opened       bool
	All          bool
	Confidential bool

	// Pagination
	Page    int
	PerPage int

	// Other
	In string

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
			// support repo override
			opts.BaseRepo = f.BaseRepo
			opts.HTTPClient = f.HttpClient

			if len(opts.Labels) != 0 && len(opts.NotLabels) != 0 {
				return cmdutils.FlagError{
					Err: errors.New("flags --label and --not-label are mutually exclusive"),
				}
			}

			if opts.Author != "" && len(opts.NotAuthor) != 0 {
				return cmdutils.FlagError{
					Err: errors.New("flags --author and --not-author are mutually exclusive"),
				}
			}

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
	issueListCmd.Flags().StringVarP(&opts.Assignee, "assignee", "a", "", "Filter issue by assignee <username>")
	issueListCmd.Flags().StringVar(&opts.Author, "author", "", "Filter issue by author <username>")
	issueListCmd.Flags().StringSliceVar(&opts.NotAuthor, "not-author", []string{}, "Filter by not being by author(s) <username>")
	issueListCmd.Flags().StringVar(&opts.Search, "search", "", "Search <string> in the fields defined by --in")
	issueListCmd.Flags().StringVar(&opts.In, "in", "title,description", "search in {title|description}")
	issueListCmd.Flags().StringSliceVarP(&opts.Labels, "label", "l", []string{}, "Filter issue by label <name>")
	issueListCmd.Flags().StringSliceVar(&opts.NotLabels, "not-label", []string{}, "Filter issue by lack of label <name>")
	issueListCmd.Flags().StringVarP(&opts.Milestone, "milestone", "m", "", "Filter issue by milestone <id>")
	issueListCmd.Flags().BoolVarP(&opts.All, "all", "A", false, "Get all issues")
	issueListCmd.Flags().BoolVarP(&opts.Closed, "closed", "c", false, "Get only closed issues")
	issueListCmd.Flags().BoolVarP(&opts.Confidential, "confidential", "C", false, "Filter by confidential issues")
	issueListCmd.Flags().IntVarP(&opts.Page, "page", "p", 1, "Page number")
	issueListCmd.Flags().IntVarP(&opts.PerPage, "per-page", "P", 30, "Number of items to list per page. (default 30)")

	issueListCmd.Flags().BoolP("opened", "o", false, "Get only opened issues")
	issueListCmd.Flags().MarkHidden("opened")
	issueListCmd.Flags().MarkDeprecated("opened", "default if --closed is not used")

	issueListCmd.Flags().BoolVarP(&opts.Mine, "mine", "M", false, "Filter only issues issues assigned to me")
	issueListCmd.Flags().MarkHidden("mine")
	issueListCmd.Flags().MarkDeprecated("mine", "use --assignee=@me")

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
		In:    gitlab.String(opts.In),
	}
	listOpts.Page = 1
	listOpts.PerPage = 30

	if opts.Assignee != "" || opts.Mine {
		if opts.Assignee == "@me" || opts.Mine {
			u, err := api.CurrentUser(nil)
			if err != nil {
				return err
			}
			opts.Assignee = u.Username
		}
		listOpts.AssigneeUsername = gitlab.String(opts.Assignee)
	}
	if opts.Author != "" {
		var u *gitlab.User
		if opts.Author == "@me" {
			u, err = api.CurrentUser(nil)
			if err != nil {
				return err
			}
		} else {
			// go-gitlab still doesn't have an AuthorUsername field so convert to ID
			u, err = api.UserByName(apiClient, opts.Author)
			if err != nil {
				return err
			}
		}
		listOpts.AuthorID = gitlab.Int(u.ID)
	}
	if len(opts.NotAuthor) != 0 {
		us, err := api.UsersByNames(apiClient, opts.NotAuthor)
		if err != nil {
			return err
		}
		var IDs []int
		for i := range us {
			IDs = append(IDs, us[i].ID)
		}
		listOpts.NotAuthorID = IDs
	}
	if opts.Search != "" {
		listOpts.Search = gitlab.String(opts.Search)
		opts.ListType = "search"
	}
	if len(opts.Labels) != 0 {
		listOpts.Labels = opts.Labels
		opts.ListType = "search"
	}
	if len(opts.NotLabels) != 0 {
		listOpts.NotLabels = opts.NotLabels
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

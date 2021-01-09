package list

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/internal/glrepo"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type ListOptions struct {
	// metadata
	Assignee  []string
	Author    string
	Labels    []string
	Milestone string
	Mine      bool

	// issue states
	State  string
	Closed bool
	Merged bool
	All    bool

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

	var mrListCmd = &cobra.Command{
		Use:     "list [flags]",
		Short:   `List merge requests`,
		Long:    ``,
		Aliases: []string{"ls"},
		Example: heredoc.Doc(`
			$ glab mr list --all
			$ glab mr ls -a
			$ glab mr list --mine
			$ glab mr list --label needs-review
			$ glab mr list -o --per-page 10
		`),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			// supports repo override
			opts.BaseRepo = f.BaseRepo
			opts.HTTPClient = f.HttpClient

			// check if any of the two or all of states flag are specified
			if opts.Closed && opts.Merged {
				return cmdutils.FlagError{
					Err: errors.New("specify either --closed or --merged. Use --all issues in all states"),
				}
			}
			if opts.All {
				opts.State = "all"
			} else if opts.Closed {
				opts.State = "closed"
				opts.TitleQualifier = opts.State
			} else if opts.Merged {
				opts.State = "merged"
				opts.TitleQualifier = opts.State
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

	mrListCmd.Flags().StringSliceVarP(&opts.Labels, "label", "l", []string{}, "Filter merge request by label <name>")
	mrListCmd.Flags().StringVar(&opts.Author, "author", "", "Fitler merge request by Author <username>")
	mrListCmd.Flags().StringVarP(&opts.Milestone, "milestone", "m", "", "Filter merge request by milestone <id>")
	mrListCmd.Flags().BoolVarP(&opts.All, "all", "A", false, "Get all merge requests")
	mrListCmd.Flags().BoolVarP(&opts.Closed, "closed", "c", false, "Get only closed merge requests")
	mrListCmd.Flags().BoolVarP(&opts.Merged, "merged", "M", false, "Get only merged merge requests")
	mrListCmd.Flags().IntVarP(&opts.Page, "page", "p", 1, "Page number")
	mrListCmd.Flags().IntVarP(&opts.PerPage, "per-page", "P", 30, "Number of items to list per page")
	mrListCmd.Flags().BoolVarP(&opts.Mine, "mine", "", false, "Get only merge requests assigned to me")
	mrListCmd.Flags().StringSliceVarP(&opts.Assignee, "assignee", "a", []string{}, "Get only merge requests assigned to users")

	mrListCmd.Flags().BoolP("opened", "o", false, "Get only open merge requests")
	mrListCmd.Flags().MarkHidden("opened")
	mrListCmd.Flags().MarkDeprecated("opened", "default value if neither --closed, --locked or --merged is used")

	return mrListCmd
}

func listRun(opts *ListOptions) error {
	var mergeRequests []*gitlab.MergeRequest

	apiClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}

	repo, err := opts.BaseRepo()
	if err != nil {
		return err
	}

	l := &gitlab.ListProjectMergeRequestsOptions{
		State: gitlab.String(opts.State),
	}
	l.Page = 1
	l.PerPage = 30

	if opts.Author != "" {
		u, err := api.UserByName(apiClient, opts.Author)
		if err != nil {
			return err
		}
		l.AuthorID = gitlab.Int(u.ID)
		opts.ListType = "search"
	}
	if len(opts.Labels) > 0 {
		l.Labels = opts.Labels
		opts.ListType = "search"
	}
	if opts.Milestone != "" {
		l.Milestone = gitlab.String(opts.Milestone)
		opts.ListType = "search"
	}
	if opts.Page != 0 {
		l.Page = opts.Page
	}
	if opts.PerPage != 0 {
		l.PerPage = opts.PerPage
	}

	if opts.Mine {
		l.Scope = gitlab.String("assigned_to_me")
		opts.ListType = "search"
	}

	assigneeIds := make([]int, 0)
	if len(opts.Assignee) > 0 {
		users, err := api.UsersByNames(apiClient, opts.Assignee)
		if err != nil {
			return err
		}
		for _, user := range users {
			assigneeIds = append(assigneeIds, user.ID)
		}
	}

	if len(assigneeIds) > 0 {
		mergeRequests, err = api.ListMRsWithAssignees(apiClient, repo.FullName(), l, assigneeIds)

	} else {
		mergeRequests, err = api.ListMRs(apiClient, repo.FullName(), l)
	}
	if err != nil {
		return err
	}

	title := utils.NewListTitle(opts.TitleQualifier + " merge request")
	title.RepoName = repo.FullName()
	title.Page = l.Page
	title.ListActionType = opts.ListType
	title.CurrentPageTotal = len(mergeRequests)

	if err = opts.IO.StartPager(); err != nil {
		return err
	}
	defer opts.IO.StopPager()
	fmt.Fprintf(opts.IO.StdOut, "%s\n%s\n", title.Describe(), mrutils.DisplayAllMRs(mergeRequests, repo.FullName()))

	return nil
}

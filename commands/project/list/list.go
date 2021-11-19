package list

import (
	"fmt"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/tableprinter"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type Options struct {
	OrderBy       string
	Sort          string
	PerPage       int
	Page          int
	FilterAll     bool
	FilterOwned   bool
	FilterMember  bool
	FilterStarred bool

	BaseRepo   func() (glrepo.Interface, error)
	HTTPClient func() (*gitlab.Client, error)
	IO         *iostreams.IOStreams
}

func NewCmdList(f *cmdutils.Factory) *cobra.Command {
	opts := &Options{
		IO: f.IO,
	}
	var repoListCmd = &cobra.Command{
		Use:   "list",
		Short: `Get list of repositories.`,
		Example: heredoc.Doc(`
	$ glab repo list
	`),
		Args:    cobra.ExactArgs(0),
		Aliases: []string{"users"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Support `-R, --repo` override
			opts.BaseRepo = f.BaseRepo
			opts.HTTPClient = f.HttpClient

			return runE(opts)
		},
	}

	cmdutils.EnableRepoOverride(repoListCmd, f)

	repoListCmd.Flags().StringVarP(&opts.OrderBy, "order", "o", "last_activity_at", "Return repositories ordered by id, created_at, or other fields")
	repoListCmd.Flags().StringVarP(&opts.Sort, "sort", "s", "", "Return repositories sorted in asc or desc order")
	repoListCmd.Flags().IntVarP(&opts.Page, "page", "p", 1, "Page number")
	repoListCmd.Flags().IntVarP(&opts.PerPage, "per-page", "P", 30, "Number of items to list per page.")
	repoListCmd.Flags().BoolVarP(&opts.FilterAll, "all", "a", false, "List all projects on the instance")
	repoListCmd.Flags().BoolVarP(&opts.FilterOwned, "mine", "m", true, "Only list projects you own")
	repoListCmd.Flags().BoolVar(&opts.FilterMember, "member", false, "Only list projects which you are a member")
	repoListCmd.Flags().BoolVar(&opts.FilterStarred, "starred", false, "Only list starred projects")
	return repoListCmd
}

func runE(opts *Options) error {
	var err error
	c := opts.IO.Color()

	apiClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}

	l := &gitlab.ListProjectsOptions{
		OrderBy: gitlab.String(opts.OrderBy),
	}

	// Other filters only valid if FilterAll not true
	if !opts.FilterAll {
		if opts.FilterOwned {
			l.Owned = gitlab.Bool(opts.FilterOwned)
		} else if opts.FilterStarred {
			l.Starred = gitlab.Bool(opts.FilterStarred)
		} else if opts.FilterMember {
			l.Membership = gitlab.Bool(opts.FilterMember)
		}
	}

	if opts.Sort != "" {
		l.Sort = gitlab.String(opts.Sort)
	}

	l.PerPage = opts.PerPage
	l.Page = opts.Page

	projects, resp, err := apiClient.Projects.ListProjects(l)
	if err != nil {
		return err
	}

	// Title
	title := fmt.Sprintf("Showing %d of %d projects (Page %d of %d)\n", len(projects), resp.TotalItems, resp.CurrentPage, resp.TotalPages)

	// List
	table := tableprinter.NewTablePrinter()
	for _, prj := range projects {
		table.AddCell(c.Blue(prj.PathWithNamespace))
		table.AddCell(prj.SSHURLToRepo)
		table.AddCell(prj.Description)
		table.EndRow()
	}

	fmt.Fprintf(opts.IO.StdOut, "%s\n%s\n", title, table.String())
	return err
}

package issueutils

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/api"
	"golang.org/x/sync/errgroup"

	"github.com/profclems/glab/pkg/tableprinter"
	"github.com/profclems/glab/pkg/utils"

	"github.com/xanzy/go-gitlab"
)

func DisplayIssueList(c *iostreams.ColorPalette, issues []*gitlab.Issue, projectID string) string {
	table := tableprinter.NewTablePrinter()
	for _, issue := range issues {
		table.AddCell(IssueState(c, issue))
		table.AddCell(issue.Title)

		if len(issue.Labels) > 0 {
			table.AddCellf("(%s)", c.Cyan(strings.Trim(strings.Join(issue.Labels, ", "), ",")))
		} else {
			table.AddCell("")
		}

		table.AddCell(c.Gray(utils.TimeToPrettyTimeAgo(*issue.CreatedAt)))
		table.EndRow()
	}

	return table.Render()
}

func DisplayIssue(c *iostreams.ColorPalette, i *gitlab.Issue) string {
	duration := utils.TimeToPrettyTimeAgo(*i.CreatedAt)
	issueID := IssueState(c, i)

	return fmt.Sprintf("%s %s (%s)\n %s\n",
		issueID, i.Title, duration, i.WebURL)
}

func IssueState(c *iostreams.ColorPalette, i *gitlab.Issue) (issueID string) {
	if i.State == "opened" {
		issueID = c.Green(fmt.Sprintf("#%d", i.IID))
	} else {
		issueID = c.Red(fmt.Sprintf("#%d", i.IID))
	}
	return
}

func IssuesFromArgs(apiClient *gitlab.Client, baseRepoFn func() (glrepo.Interface, error), args []string) ([]*gitlab.Issue, glrepo.Interface, error) {
	baseRepo, err := baseRepoFn()
	if err != nil {
		return nil, nil, err
	}
	if len(args) <= 1 {
		if len(args) == 1 {
			args = strings.Split(args[0], ",")
		}
		if len(args) <= 1 {
			issue, repo, err := IssueFromArg(apiClient, baseRepoFn, args[0])
			if err != nil {
				return nil, nil, err
			}
			baseRepo = repo
			return []*gitlab.Issue{issue}, baseRepo, err
		}
	}

	errGroup, _ := errgroup.WithContext(context.Background())
	issues := make([]*gitlab.Issue, len(args))
	for i, arg := range args {
		i, arg := i, arg
		errGroup.Go(func() error {
			issue, repo, err := IssueFromArg(apiClient, baseRepoFn, arg)
			if err != nil {
				return err
			}
			baseRepo = repo
			issues[i] = issue
			return nil
		})
	}
	if err := errGroup.Wait(); err != nil {
		return nil, nil, err
	}
	return issues, baseRepo, nil

}

func IssueFromArg(apiClient *gitlab.Client, baseRepoFn func() (glrepo.Interface, error), arg string) (*gitlab.Issue, glrepo.Interface, error) {
	issueIID, baseRepo := issueMetadataFromURL(arg)
	if issueIID == 0 {
		var err error
		issueIID, err = strconv.Atoi(strings.TrimPrefix(arg, "#"))
		if err != nil {
			return nil, nil, fmt.Errorf("invalid issue format: %q", arg)
		}
	}

	if baseRepo == nil {
		var err error
		baseRepo, err = baseRepoFn()
		if err != nil {
			return nil, nil, fmt.Errorf("could not determine base repo: %w", err)
		}
	} else {
		// initialize a new HTTP Client with the new host
		// TODO: avoid reinitializing the config, get the config as a parameter

		cfg, _ := config.Init()
		a, err := api.NewClientWithCfg(baseRepo.RepoHost(), cfg, false)
		if err != nil {
			return nil, nil, err
		}
		apiClient = a.Lab()
	}

	issue, err := issueFromIID(apiClient, baseRepo, issueIID)
	return issue, baseRepo, err
}

// FIXME: have a single regex to match either of the following
//  OWNER/REPO/issues/id
//  GROUP/NAMESPACE/REPO/issues/id
var issueURLPersonalRE = regexp.MustCompile(`^/([^/]+)/([^/]+)/issues/(\d+)`)
var issueURLGroupRE = regexp.MustCompile(`^/([^/]+)/([^/]+)/([^/]+)/issues/(\d+)`)

func issueMetadataFromURL(s string) (int, glrepo.Interface) {
	u, err := url.Parse(s)
	if err != nil {
		return 0, nil
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return 0, nil
	}

	u.Path = strings.Replace(u.Path, "/-/issues", "/issues", 1)

	m := issueURLPersonalRE.FindStringSubmatch(u.Path)
	if m == nil {
		m = issueURLGroupRE.FindStringSubmatch(u.Path)
		if m == nil {
			return 0, nil
		}
	}
	var issueIID int
	if len(m) > 0 {
		issueIID, _ = strconv.Atoi(m[len(m)-1])
	}

	u.Path = strings.Replace(u.Path, fmt.Sprintf("/issues/%d", issueIID), "", 1)

	repo, err := glrepo.FromURL(u)
	if err != nil {
		return 0, nil
	}
	return issueIID, repo
}

func issueFromIID(apiClient *gitlab.Client, repo glrepo.Interface, issueIID int) (*gitlab.Issue, error) {
	return api.GetIssue(apiClient, repo.FullName(), issueIID)
}

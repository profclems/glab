package mrutils

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/prompt"
	"github.com/profclems/glab/pkg/tableprinter"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
)

type MRCheckErrOptions struct {
	// WorkInProgress: check and return err if merge request is a DRAFT
	WorkInProgress bool
	// Closed : check and return err if merge request is closed
	Closed bool
	// Merged : check and return err if merge request is already merged
	Merged bool
	// Opened : check and return err if merge request is already opened
	Opened bool
	// Conflict : check and return err if there are merge conflicts
	Conflict bool
	// PipelineStatus : check and return err pipeline did not succeed and it is required before merging
	PipelineStatus bool
	// MergePermitted : check and return err if user is not authorized to merge
	MergePermitted bool
	// Subscribed : check and return err if user is already subscribed to MR
	Subscribed bool
	// Unsubscribed : check and return err if user is already unsubscribed to MR
	Unsubscribed bool
	// MergePrivilege : check and return err if user is not authorized to merge
	MergePrivilege bool
}

// MRCheckErrors checks and return merge request errors specified in MRCheckErrOptions{}
func MRCheckErrors(mr *gitlab.MergeRequest, err MRCheckErrOptions) error {
	if mr.WorkInProgress && err.WorkInProgress {
		return fmt.Errorf("this merge request is still a work in progress. Run `glab mr update %d --ready` to mark it as ready for review", mr.IID)
	}

	if mr.MergeWhenPipelineSucceeds && err.PipelineStatus && mr.Pipeline != nil {
		if mr.Pipeline.Status != "success" {
			return fmt.Errorf("pipeline for this merge request has failed. Pipeline is required to succeed before merging")
		}
	}

	if mr.State == "merged" && err.Merged {
		return fmt.Errorf("this merge request has already been merged")
	}

	if mr.State == "closed" && err.Closed {
		return fmt.Errorf("this merge request has been closed")
	}

	if mr.State == "opened" && err.Opened {
		return fmt.Errorf("this merge request is already open")
	}

	if mr.Subscribed && err.Subscribed {
		return fmt.Errorf("you are already subscribed to this merge request")
	}

	if !mr.Subscribed && err.Unsubscribed {
		return fmt.Errorf("you are not subscribed to this merge request")
	}

	if err.MergePrivilege && !mr.User.CanMerge {
		return fmt.Errorf("you do not have enough priviledges to merge this merge request")
	}

	if err.Conflict && mr.HasConflicts {
		return fmt.Errorf("there are merge conflicts. Resolve conflicts and try again or merge locally")
	}

	return nil
}

func DisplayMR(c *iostreams.ColorPalette, mr *gitlab.MergeRequest, isTTY bool) string {
	mrID := MRState(c, mr)
	if isTTY {
		return fmt.Sprintf("%s %s (%s)\n %s\n", mrID, mr.Title, mr.SourceBranch, mr.WebURL)
	} else {
		return mr.WebURL
	}
}

func MRState(c *iostreams.ColorPalette, m *gitlab.MergeRequest) string {
	if m.State == "opened" {
		return c.Green(fmt.Sprintf("!%d", m.IID))
	} else if m.State == "merged" {
		return c.Magenta(fmt.Sprintf("!%d", m.IID))
	} else {
		return c.Red(fmt.Sprintf("!%d", m.IID))
	}
}

func DisplayAllMRs(streams *iostreams.IOStreams, mrs []*gitlab.MergeRequest, projectID string) string {
	c := streams.Color()
	table := tableprinter.NewTablePrinter()
	table.SetIsTTY(streams.IsOutputTTY())
	for _, m := range mrs {
		table.AddCell(streams.Hyperlink(MRState(c, m), m.WebURL))
		table.AddCell(m.Title)
		table.AddCell(c.Cyan(fmt.Sprintf("(%s) ← (%s)", m.TargetBranch, m.SourceBranch)))
		table.EndRow()
	}

	return table.Render()
}

//MRFromArgs is wrapper around MRFromArgsWithOpts without any custom options
func MRFromArgs(f *cmdutils.Factory, args []string, state string) (*gitlab.MergeRequest, glrepo.Interface, error) {
	return MRFromArgsWithOpts(f, args, &gitlab.GetMergeRequestsOptions{}, state)
}

//MRFromArgsWithOpts gets MR with custom request options passed down to it
func MRFromArgsWithOpts(
	f *cmdutils.Factory,
	args []string,
	opts *gitlab.GetMergeRequestsOptions,
	state string,
) (*gitlab.MergeRequest, glrepo.Interface, error) {
	var mrID int
	var mr *gitlab.MergeRequest

	apiClient, err := f.HttpClient()
	if err != nil {
		return nil, nil, err
	}

	baseRepo, err := f.BaseRepo()
	if err != nil {
		return nil, nil, err
	}

	var branch string

	if len(args) > 0 {
		mrID, err = strconv.Atoi(strings.TrimPrefix(args[0], "!"))
		if err != nil {
			branch = args[0]
		} else if mrID == 0 { // to check for cases where the user explicitly specified mrID to be zero
			return nil, nil, fmt.Errorf("invalid merge request ID provided")
		}
	}

	if branch == "" && mrID == 0 {
		branch, err = f.Branch()
		if err != nil {
			return nil, nil, err
		}
	}

	if mrID == 0 {
		mr, err = getMRForBranch(apiClient, baseRepo, branch, state)
		if err != nil {
			return nil, nil, err
		}
		mrID = mr.IID
	}
	mr, err = api.GetMR(apiClient, baseRepo.FullName(), mrID, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get merge request %d: %w", mrID, err)
	}

	return mr, baseRepo, nil
}

func MRsFromArgs(f *cmdutils.Factory, args []string, state string) ([]*gitlab.MergeRequest, glrepo.Interface, error) {
	if len(args) <= 1 {
		var arrIDs []string
		if len(args) == 1 {
			arrIDs = strings.Split(args[0], ",")
		}
		if len(arrIDs) <= 1 {
			// If there are no args then try to auto-detect from the branch name
			mr, baseRepo, err := MRFromArgs(f, args, state)
			if err != nil {
				return nil, nil, err
			}
			return []*gitlab.MergeRequest{mr}, baseRepo, nil
		}
		args = arrIDs
	}

	baseRepo, err := f.BaseRepo()
	if err != nil {
		return nil, nil, err
	}

	errGroup, _ := errgroup.WithContext(context.Background())
	mrs := make([]*gitlab.MergeRequest, len(args))
	for i, arg := range args {
		i, arg := i, arg
		errGroup.Go(func() error {
			// fetching multiple MRs does not return many major params in the payload
			// so we fetch again using the single mr endpoint
			mr, _, err := MRFromArgs(f, []string{arg}, state)
			if err != nil {
				return err
			}
			mrs[i] = mr
			return nil
		})
	}
	if err := errGroup.Wait(); err != nil {
		return nil, nil, err
	}
	return mrs, baseRepo, nil

}

var getMRForBranch = func(apiClient *gitlab.Client, baseRepo glrepo.Interface, arg string, state string) (*gitlab.MergeRequest, error) {
	currentBranch := arg // Assume the user is using only 'branch', not 'OWNER:branch'
	var owner string

	// If the string contains a ':' then it is using the OWNER:branch format, split it and
	// assign them to the appropriate values, do note that we do not expect multiple ':' as
	// git does not allow ':' to be used on branch names
	if strings.Contains(arg, ":") {
		t := strings.Split(arg, ":")
		owner = t[0]
		currentBranch = t[1]
	}

	opts := gitlab.ListProjectMergeRequestsOptions{
		SourceBranch: gitlab.String(currentBranch),
	}

	// Set the state value if it is not empty, if it is empty then it will look at all possible
	// values, 'any' is also a descriptive keyword used in the source code that is equivalent to
	// passing nothing on
	if state != "" && state != "any" {
		opts.State = gitlab.String(state)
	}

	mrs, err := api.ListMRs(apiClient, baseRepo.FullName(), &opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get open merge request for %q: %w", currentBranch, err)
	}

	if len(mrs) == 0 {
		return nil, fmt.Errorf("no open merge request available for %q", currentBranch)
	}

	// The user gave us an 'OWNER:' so try to match the merge request with it
	if owner != "" {
		for i := range mrs {
			// We found a match!
			if mrs[i].Author.Username == owner {
				return mrs[i], nil
			}
		}
		// No match, error out, tell the user which branch and which username we looked for
		return nil, fmt.Errorf("no open merge request available for %q owned by @%s", currentBranch, owner)
	}

	// This is done after the 'OWNER:' check because we don't want to give the wrong MR
	// to someone that **explicitly** asked for a OWNER.
	if len(mrs) == 1 {
		return mrs[0], nil
	}

	// No 'OWNER:' prompt the user to pick a merge request
	mrMap := map[string]*gitlab.MergeRequest{}
	var mrNames []string
	for i := range mrs {
		t := fmt.Sprintf("!%d (%s) by @%s", mrs[i].IID, currentBranch, mrs[i].Author.Username)
		mrMap[t] = mrs[i]
		mrNames = append(mrNames, t)
	}
	pickedMR := mrNames[0]
	err = prompt.Select(&pickedMR,
		"mr",
		"There are multiple merge requests matching the requested branch, pick one",
		mrNames,
	)
	if err != nil {
		return nil, fmt.Errorf("a merge request must be picked: %w", err)
	}
	return mrMap[pickedMR], nil
}

func RebaseMR(ios *iostreams.IOStreams, apiClient *gitlab.Client, repo glrepo.Interface, mr *gitlab.MergeRequest) error {
	ios.StartSpinner("Sending rebase request...")
	err := api.RebaseMR(apiClient, repo.FullName(), mr.IID)
	if err != nil {
		return err
	}
	ios.StopSpinner("")

	opts := &gitlab.GetMergeRequestsOptions{}
	opts.IncludeRebaseInProgress = gitlab.Bool(true)
	ios.StartSpinner("Checking rebase status...")
	errorMSG := ""
	i := 0
	for {
		mr, err := api.GetMR(apiClient, repo.FullName(), mr.IID, opts)
		if err != nil {
			errorMSG = err.Error()
			break
		}
		if i == 0 {
			ios.StopSpinner("")
			ios.StartSpinner("Rebase in progress...")
		}
		if !mr.RebaseInProgress {
			if mr.MergeError != "" && mr.MergeError != "null" {
				errorMSG = mr.MergeError
			}
			break
		}
		i++
	}
	ios.StopSpinner("")
	if errorMSG != "" {
		return errors.New(errorMSG)
	}
	fmt.Fprintln(ios.StdOut, ios.Color().GreenCheck(), "Rebase successful")
	return nil
}

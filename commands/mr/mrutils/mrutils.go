package mrutils

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
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
		return fmt.Errorf("you are already unsubscribed to this merge request")
	}

	if err.MergePrivilege && !mr.User.CanMerge {
		return fmt.Errorf("you do not have enough priviledges to merge this merge request")
	}

	if err.Conflict && mr.HasConflicts {
		return fmt.Errorf("there are merge conflicts. Resolve conflicts and try again or merge locally")
	}

	return nil
}

func DisplayMR(mr *gitlab.MergeRequest) string {
	mrID := MRState(mr)
	return fmt.Sprintf("%s %s (%s)\n %s\n",
		mrID, mr.Title, mr.SourceBranch, mr.WebURL)
}

func MRState(m *gitlab.MergeRequest) string {
	if m.State == "opened" {
		return utils.Green(fmt.Sprintf("!%d", m.IID))
	} else if m.State == "merged" {
		return utils.Blue(fmt.Sprintf("!%d", m.IID))
	} else {
		return utils.Red(fmt.Sprintf("!%d", m.IID))
	}
}

func DisplayAllMRs(mrs []*gitlab.MergeRequest, projectID string) string {
	table := tableprinter.NewTablePrinter()
	for _, m := range mrs {
		table.AddCell(MRState(m))
		table.AddCell(m.Title)
		table.AddCell(utils.Cyan(fmt.Sprintf("(%s) ← (%s)", m.TargetBranch, m.SourceBranch)))
		table.EndRow()
	}

	return table.Render()
}

func MRFromArgs(f *cmdutils.Factory, args []string) (*gitlab.MergeRequest, glrepo.Interface, error) {
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

	branch, err := f.Branch()
	if err != nil {
		return nil, nil, err
	}

	if len(args) > 0 {
		mrID, err = strconv.Atoi(args[0])
		if err != nil {
			branch = args[0]
		} else if mrID == 0 { // to check for cases where the user explicitly specified mrID to be zero
			return nil, nil, fmt.Errorf("invalid merge request ID provided")
		}
	}

	if mrID == 0 {
		mr, err = GetOpenMRForBranch(apiClient, baseRepo, branch)
		if err != nil {
			return nil, nil, err
		}
		mrID = mr.IID
	}
	// fetching multiple MRs does not return many major params in the payload
	// so we fetch again using the single mr endpoint
	mr, err = api.GetMR(apiClient, baseRepo.FullName(), mrID, &gitlab.GetMergeRequestsOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get merge request %d: %w", mrID, err)
	}

	return mr, baseRepo, nil
}

func MRsFromArgs(f *cmdutils.Factory, args []string) ([]*gitlab.MergeRequest, glrepo.Interface, error) {
	if len(args) <= 1 {
		var arrIDs []string
		if len(args) == 1 {
			arrIDs = strings.Split(args[0], ",")
		}
		if len(arrIDs) <= 1 {
			mr, baseRepo, err := MRFromArgs(f, args)
			if err != nil {
				return nil, nil, err
			}
			return []*gitlab.MergeRequest{mr}, baseRepo, err
		}
		args = arrIDs
	}

	apiClient, err := f.HttpClient()
	if err != nil {
		return nil, nil, err
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
			mrID, err := strconv.Atoi(arg)
			if err != nil {
				return err
			}
			if mrID == 0 {
				return fmt.Errorf("invalid merge request ID provided")
			}
			// fetching multiple MRs does not return many major params in the payload
			// so we fetch again using the single mr endpoint
			mr, err := api.GetMR(apiClient, baseRepo.FullName(), mrID, &gitlab.GetMergeRequestsOptions{})
			if err != nil {
				return fmt.Errorf("failed to get merge request %d: %w", mrID, err)
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

func GetOpenMRForBranch(apiClient *gitlab.Client, baseRepo glrepo.Interface, currentBranch string) (*gitlab.MergeRequest, error) {
	mrs, err := api.ListMRs(apiClient, baseRepo.FullName(), &gitlab.ListProjectMergeRequestsOptions{
		SourceBranch: gitlab.String(currentBranch),
		State:        gitlab.String("opened"),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get open merge request for %q: %w", currentBranch, err)
	}
	if len(mrs) == 0 {
		return nil, fmt.Errorf("no open merge request available for %q", currentBranch)
	}
	// A single result is expected since gitlab does not allow multiple merge requests for a single source branch
	return mrs[0], nil
}

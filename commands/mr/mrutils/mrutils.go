package mrutils

import (
	"fmt"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/profclems/glab/pkg/tableprinter"
	"github.com/xanzy/go-gitlab"
	"strconv"
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
}

// MRCheckErrors checks and return merge request errors specified in MRCheckErrOptions{}
func MRCheckErrors(mr *gitlab.MergeRequest, opts MRCheckErrOptions) error {
	if mr.WorkInProgress && opts.WorkInProgress {
		return fmt.Errorf("this merge request is still a work in progress. Run `glab mr update %d --ready` to mark it as ready for review", mr.IID)
	}

	if mr.MergeWhenPipelineSucceeds && opts.PipelineStatus && mr.Pipeline != nil {
		if mr.Pipeline.Status != "success" {
			return fmt.Errorf("pipeline for this merge request has failed. Pipeline is required to succeed before merging")
		}
	}

	if mr.State == "merged" && opts.Merged {
		return fmt.Errorf("this merge request has already been merged")
	}

	if mr.State == "closed" && opts.Closed {
		return fmt.Errorf("this merge request has been closed")
	}

	if mr.State == "opened" && opts.Opened {
		return fmt.Errorf("this merge request is already open")
	}

	if mr.Subscribed && opts.Subscribed {
		return fmt.Errorf("you are already subscribed to this merge request")
	}

	if !mr.Subscribed && opts.Unsubscribed {
		return fmt.Errorf("you are already unsubscribed to this merge request")
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
	} else {
		mr, err = api.GetMR(apiClient, baseRepo.FullName(), mrID, &gitlab.GetMergeRequestsOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get merge request %d: %w", mrID, err)
		}

	}

	return mr, baseRepo, nil
}

func GetOpenMRForBranch(apiClient *gitlab.Client, baseRepo glrepo.Interface, currentBranch string) (*gitlab.MergeRequest, error)  {
	mrs, err := api.ListMRs(apiClient, baseRepo.FullName(), &gitlab.ListProjectMergeRequestsOptions{
		SourceBranch: gitlab.String(currentBranch),
		State: gitlab.String("opened"),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get open merge request for %q: %w", currentBranch, err)
	}
	if len(mrs) == 0 {
		return nil, fmt.Errorf("no open merge request availabe for %q", currentBranch)
	}
	// A single result is expected since gitlab does not allow multiple merge requests for a single source branch
	return mrs[0], nil
}

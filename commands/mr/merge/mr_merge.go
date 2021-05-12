package merge

import (
	"errors"
	"fmt"
	"time"

	"github.com/profclems/glab/pkg/surveyext"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/avast/retry-go"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/pkg/prompt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type MRMergeMethod int

const (
	MRMergeMethodMerge MRMergeMethod = iota
	MRMergeMethodSquash
	MRMergeMethodRebase
)

type MergeOpts struct {
	MergeWhenPipelineSucceeds bool
	SquashBeforeMerge         bool
	RebaseBeforeMerge         bool
	RemoveSourceBranch        bool

	SquashMessage      string
	MergeCommitMessage string
	SHA                string

	MergeMethod MRMergeMethod
}

func NewCmdMerge(f *cmdutils.Factory) *cobra.Command {
	var opts = &MergeOpts{
		MergeMethod: MRMergeMethodMerge,
	}

	var mrMergeCmd = &cobra.Command{
		Use:     "merge {<id> | <branch>}",
		Short:   `Merge/Accept merge requests`,
		Long:    ``,
		Aliases: []string{"accept"},
		Example: heredoc.Doc(`
			$ glab mr merge 235
			$ glab mr accept 235
			$ glab mr merge    # Finds open merge request from current branch
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			c := f.IO.Color()

			if opts.SquashBeforeMerge && opts.RebaseBeforeMerge {
				return &cmdutils.FlagError{Err: errors.New("only one of --rebase, or --squash can be enabled")}
			}

			if !opts.SquashBeforeMerge && opts.SquashMessage != "" {
				return &cmdutils.FlagError{Err: errors.New("--squash-message can only be used with --squash")}
			}

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mr, repo, err := mrutils.MRFromArgs(f, args, "opened")
			if err != nil {
				return err
			}

			if err = mrutils.MRCheckErrors(mr, mrutils.MRCheckErrOptions{
				WorkInProgress: true,
				Closed:         true,
				Merged:         true,
				Conflict:       true,
				PipelineStatus: true,
				MergePrivilege: true,
			}); err != nil {
				return err
			}

			if !cmd.Flags().Changed("when-pipeline-succeeds") && f.IO.IsOutputTTY() && mr.Pipeline != nil {
				_ = prompt.Confirm(&opts.MergeWhenPipelineSucceeds, "Merge when pipeline succeeds?", true)
			}

			if f.IO.IsOutputTTY() {
				if !opts.SquashBeforeMerge && !opts.RebaseBeforeMerge && opts.MergeCommitMessage == "" {
					opts.MergeMethod, err = mergeMethodSurvey()
					if err != nil {
						return err
					}
					if opts.MergeMethod == MRMergeMethodSquash {
						opts.SquashBeforeMerge = true
					} else if opts.MergeMethod == MRMergeMethodRebase {
						opts.RebaseBeforeMerge = true
					}
				}

				if opts.MergeCommitMessage == "" && opts.SquashMessage == "" {
					action, err := confirmSurvey(opts.MergeMethod != MRMergeMethodRebase)
					if err != nil {
						return fmt.Errorf("unable to prompt: %w", err)
					}

					if action == cmdutils.EditCommitMessageAction {
						var mergeMessage string

						editor, err := cmdutils.GetEditor(f.Config)
						if err != nil {
							return err
						}
						mergeMessage, err = surveyext.Edit(editor, "*.md", mr.Title, f.IO.In, f.IO.StdOut, f.IO.StdErr, nil)
						if err != nil {
							return err
						}

						if opts.SquashBeforeMerge {
							opts.SquashMessage = mergeMessage
						} else {
							opts.MergeCommitMessage = mergeMessage
						}

						action, err = confirmSurvey(false)
						if err != nil {
							return fmt.Errorf("unable to confirm: %w", err)
						}
					}
					if action == cmdutils.CancelAction {
						fmt.Fprintln(f.IO.StdErr, "Cancelled.")
						return cmdutils.SilentError
					}
				}
			}

			mergeOpts := &gitlab.AcceptMergeRequestOptions{}
			if opts.MergeCommitMessage != "" {
				mergeOpts.MergeCommitMessage = gitlab.String(opts.MergeCommitMessage)
			}
			if opts.SquashMessage != "" {
				mergeOpts.SquashCommitMessage = gitlab.String(opts.SquashMessage)
			}
			if opts.SquashBeforeMerge {
				mergeOpts.Squash = gitlab.Bool(true)
			}
			if opts.RemoveSourceBranch {
				mergeOpts.ShouldRemoveSourceBranch = gitlab.Bool(true)
			}
			if opts.MergeWhenPipelineSucceeds && mr.Pipeline != nil {
				if mr.Pipeline.Status == "canceled" || mr.Pipeline.Status == "failed" {
					fmt.Fprintln(f.IO.StdOut, c.FailedIcon(), "Pipeline Status:", mr.Pipeline.Status)
					fmt.Fprintln(f.IO.StdOut, c.FailedIcon(), "Cannot perform merge action")
					return cmdutils.SilentError
				}
				mergeOpts.MergeWhenPipelineSucceeds = gitlab.Bool(true)
			}
			if opts.SHA != "" {
				mergeOpts.SHA = gitlab.String(opts.SHA)
			}

			if opts.RebaseBeforeMerge {
				err := mrutils.RebaseMR(f.IO, apiClient, repo, mr)
				if err != nil {
					return err
				}
			}

			f.IO.StartSpinner("Merging merge request !%d", mr.IID)

			err = retry.Do(func() error {
				retry.Attempts(3)
				retry.Delay(time.Second * 6)
				mr, err = api.MergeMR(apiClient, repo.FullName(), mr.IID, mergeOpts)
				if err != nil {
					return err
				}
				return nil
			})

			if err != nil {
				return err
			}
			f.IO.StopSpinner("")
			isMerged := true
			if opts.MergeWhenPipelineSucceeds {
				if mr.Pipeline == nil {
					fmt.Fprintln(f.IO.StdOut, c.WarnIcon(), "No pipeline running on", mr.SourceBranch)
				} else {
					switch mr.Pipeline.Status {
					case "success":
						fmt.Fprintln(f.IO.StdOut, c.GreenCheck(), "Pipeline Succeeded")
					default:
						fmt.Fprintln(f.IO.StdOut, c.WarnIcon(), "Pipeline Status:", mr.Pipeline.Status)
						if mr.State != "merged" {
							fmt.Fprintln(f.IO.StdOut, c.GreenCheck(), "Will merge when pipeline succeeds")
							isMerged = false
						}
					}
				}
			}
			if isMerged {
				action := "Merged"
				switch opts.MergeMethod {
				case MRMergeMethodRebase:
					action = "Rebased and merged"
				case MRMergeMethodSquash:
					action = "Squashed and merged"
				}
				fmt.Fprintln(f.IO.StdOut, c.GreenCheck(), action)
			}
			fmt.Fprintln(f.IO.StdOut, mrutils.DisplayMR(c, mr))
			return nil
		},
	}

	mrMergeCmd.Flags().StringVarP(&opts.SHA, "sha", "", "", "Merge Commit sha")
	mrMergeCmd.Flags().BoolVarP(&opts.RemoveSourceBranch, "remove-source-branch", "d", false, "Remove source branch on merge")
	mrMergeCmd.Flags().BoolVarP(&opts.MergeWhenPipelineSucceeds, "when-pipeline-succeeds", "", true, "Merge only when pipeline succeeds")
	mrMergeCmd.Flags().StringVarP(&opts.MergeCommitMessage, "message", "m", "", "Custom merge commit message")
	mrMergeCmd.Flags().StringVarP(&opts.SquashMessage, "squash-message", "", "", "Custom Squash commit message")
	mrMergeCmd.Flags().BoolVarP(&opts.SquashBeforeMerge, "squash", "s", false, "Squash commits on merge")
	mrMergeCmd.Flags().BoolVarP(&opts.RebaseBeforeMerge, "rebase", "r", false, "Rebase the commits onto the base branch\n")

	return mrMergeCmd
}

func mergeMethodSurvey() (MRMergeMethod, error) {
	type mergeOption struct {
		title  string
		method MRMergeMethod
	}

	var mergeOpts = []mergeOption{
		{title: "Create a merge commit", method: MRMergeMethodMerge},
		{title: "Rebase and merge", method: MRMergeMethodRebase},
		{title: "Squash and merge", method: MRMergeMethodSquash},
	}

	var surveyOpts []string
	for _, v := range mergeOpts {
		surveyOpts = append(surveyOpts, v.title)
	}

	mergeQuestion := &survey.Select{
		Message: "What merge method would you like to use?",
		Options: surveyOpts,
	}

	var result int
	err := prompt.AskOne(mergeQuestion, &result)
	return mergeOpts[result].method, err
}

func confirmSurvey(allowEditMsg bool) (cmdutils.Action, error) {
	const (
		submitLabel        = "Submit"
		editCommitMsgLabel = "Edit commit message"
		cancelLabel        = "Cancel"
	)

	options := []string{submitLabel}
	if allowEditMsg {
		options = append(options, editCommitMsgLabel)
	}
	options = append(options, cancelLabel)

	var result string
	submit := &survey.Select{
		Message: "What's next?",
		Options: options,
	}
	err := prompt.AskOne(submit, &result)
	if err != nil {
		return cmdutils.CancelAction, fmt.Errorf("could not prompt: %w", err)
	}

	switch result {
	case submitLabel:
		return cmdutils.SubmitAction, nil
	case editCommitMsgLabel:
		return cmdutils.EditCommitMessageAction, nil
	default:
		return cmdutils.CancelAction, nil
	}
}

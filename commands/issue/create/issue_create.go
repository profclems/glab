package create

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/pkg/api"
	"github.com/profclems/glab/pkg/prompt"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type CreateOpts struct {
	Title       string
	Description string
	Labels      string
	Assignees   []string

	Weight    int
	MileStone int
	LinkedMR  int

	NoEditor       bool
	IsConfidential bool
	IsInteractive  bool
}

func NewCmdCreate(f *cmdutils.Factory) *cobra.Command {
	opts := &CreateOpts{}
	var issueCreateCmd = &cobra.Command{
		Use:     "create [flags]",
		Short:   `Create an issue`,
		Long:    ``,
		Aliases: []string{"new"},
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			issueCreateOpts := &gitlab.CreateIssueOptions{}

			hasTitle := cmd.Flags().Changed("title")
			hasDescription := cmd.Flags().Changed("description")

			// disable interactive mode if title and description are explicitly defined
			opts.IsInteractive = !(hasTitle && hasDescription)

			if opts.IsInteractive && !f.IO.PromptEnabled() {
				return &cmdutils.FlagError{Err: errors.New("--title and --description required for non-interactive mode")}
			}

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			var templateName string
			var templateContents string

			if opts.IsInteractive {
				if opts.Description == "" {
					if editor, _ := cmd.Flags().GetBool("no-editor"); editor {
						err = prompt.AskMultiline(&opts.Description, "Description:", "")
						if err != nil {
							return err
						}
					} else {

						templateResponse := struct {
							Index int
						}{}
						templateNames, err := cmdutils.ListGitLabTemplates(cmdutils.IssueTemplate)
						if err != nil {
							return fmt.Errorf("error getting templates: %w", err)
						}

						templateNames = append(templateNames, "Open a blank Issue")

						selectQs := []*survey.Question{
							{
								Name: "index",
								Prompt: &survey.Select{
									Message: "Choose a template",
									Options: templateNames,
								},
							},
						}

						if err := prompt.Ask(selectQs, &templateResponse); err != nil {
							return fmt.Errorf("could not prompt: %w", err)
						}
						if templateResponse.Index != len(templateNames) {
							templateName = templateNames[templateResponse.Index]
							templateContents, err = cmdutils.LoadGitLabTemplate(cmdutils.IssueTemplate, templateName)
							if err != nil {
								return fmt.Errorf("failed to get template contents: %w", err)
							}
						}
					}
				}
				if opts.Title == "" {
					err = prompt.AskQuestionWithInput(&opts.Title, "Title", "", true)
					if err != nil {
						return err
					}
				}
				if opts.Description == "" {
					if opts.NoEditor {
						err = prompt.AskMultiline(&opts.Description, "Description:", "")
						if err != nil {
							return err
						}
					} else {
						editor, err := cmdutils.GetEditor(f.Config)
						if err != nil {
							return err
						}
						err = cmdutils.DescriptionPrompt(&opts.Description, templateContents, editor)
						if err != nil {
							return err
						}
					}
				}
				if opts.Labels == "" {
					remotes, err := f.Remotes()
					if err != nil {
						return err
					}
					repoRemote, err := remotes.FindByRepo(repo.RepoOwner(), repo.RepoName())
					if err != nil {
						return err
					}
					err = cmdutils.LabelsPrompt(&opts.Labels, apiClient, repoRemote)
					if err != nil {
						return err
					}
				}
			} else {
				if opts.Title == "" {
					return fmt.Errorf("title can't be blank")
				}
			}

			issueCreateOpts.Title = gitlab.String(opts.Title)
			issueCreateOpts.Labels = gitlab.Labels{opts.Labels}
			issueCreateOpts.Description = &opts.Description
			if opts.IsConfidential {
				issueCreateOpts.Confidential = gitlab.Bool(opts.IsConfidential)
			}
			if opts.Weight != -1 {
				issueCreateOpts.Weight = gitlab.Int(opts.Weight)
			}
			if opts.LinkedMR != -1 {
				issueCreateOpts.MergeRequestToResolveDiscussionsOf = gitlab.Int(opts.LinkedMR)
			}
			if opts.MileStone != -1 {
				issueCreateOpts.MilestoneID = gitlab.Int(opts.MileStone)
			}
			if len(opts.Assignees) > 0 {
				users, err := api.UsersByNames(apiClient, opts.Assignees)
				if err != nil {
					return err
				}
				issueCreateOpts.AssigneeIDs = cmdutils.IDsFromUsers(users)
			}
			fmt.Fprintln(f.IO.StdErr, "\n- Creating issue in", repo.FullName())
			issue, err := api.CreateIssue(apiClient, repo.FullName(), issueCreateOpts)
			if err != nil {
				return err
			}
			fmt.Fprintln(f.IO.StdOut, issueutils.DisplayIssue(issue))
			return nil
		},
	}
	issueCreateCmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Supply a title for issue")
	issueCreateCmd.Flags().StringVarP(&opts.Description, "description", "d", "", "Supply a description for issue")
	issueCreateCmd.Flags().StringVarP(&opts.Labels, "label", "l", "", "Add label by name. Multiple labels should be comma separated")
	issueCreateCmd.Flags().StringSliceVarP(&opts.Assignees, "assignee", "a", []string{}, "Assign issue to people by their `usernames`")
	issueCreateCmd.Flags().IntVarP(&opts.MileStone, "milestone", "m", -1, "The global ID of a milestone to assign issue")
	issueCreateCmd.Flags().BoolVarP(&opts.IsConfidential, "confidential", "c", false, "Set an issue to be confidential. Default is false")
	issueCreateCmd.Flags().IntVarP(&opts.LinkedMR, "linked-mr", "", -1, "The IID of a merge request in which to resolve all issues")
	issueCreateCmd.Flags().IntVarP(&opts.Weight, "weight", "w", -1, "The weight of the issue. Valid values are greater than or equal to 0.")
	issueCreateCmd.Flags().BoolVarP(&opts.NoEditor, "no-editor", "", false, "Don't open editor to enter description. If set to true, uses prompt. Default is false")

	return issueCreateCmd
}

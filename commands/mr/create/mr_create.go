package create

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/profclems/glab/pkg/prompt"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type CreateOpts struct {
	Title        string
	Description  string
	SourceBranch string
	TargetBranch string
	Labels       string
	Assignees    string

	MileStone     int
	TargetProject int

	CreateSourceBranch bool
	RemoveSourceBranch bool
	AllowCollaboration bool

	Autofill      bool
	IsDraft       bool
	IsWIP         bool
	ShouldPush    bool
	NoEditor      bool
	IsInteractive bool
}

func NewCmdCreate(f *cmdutils.Factory) *cobra.Command {
	opts := &CreateOpts{}
	var mrCreateCmd = &cobra.Command{
		Use:     "create",
		Short:   `Create new merge request`,
		Long:    ``,
		Aliases: []string{"new"},
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			out := f.IO.StdOut
			mrCreateOpts := &gitlab.CreateMergeRequestOptions{}

			hasTitle := cmd.Flags().Changed("title")
			hasDescription := cmd.Flags().Changed("description")

			// disable interactive mode if title and description are explicitly defined
			opts.IsInteractive = !(hasTitle && hasDescription)

			if opts.IsInteractive && !f.IO.PromptEnabled() && !opts.Autofill {
				return &cmdutils.FlagError{Err: errors.New("--title or --fill required for non-interactive mode")}
			}

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			remotes, err := f.Remotes()
			if err != nil {
				return err
			}
			repoRemote, err := remotes.FindByRepo(repo.RepoOwner(), repo.RepoName())
			if err != nil {
				return err
			}

			if opts.TargetBranch == "" {
				opts.TargetBranch, _ = git.GetDefaultBranch(repoRemote.PushURL.String())
			}

			if opts.CreateSourceBranch && opts.SourceBranch == "" {
				opts.SourceBranch = utils.ReplaceNonAlphaNumericChars(opts.Title, "-")
			} else if opts.SourceBranch == "" {
				b, err := git.CurrentBranch()
				if err != nil {
					return err
				}
				opts.SourceBranch = b
			}

			if opts.Autofill {
				branch, err := f.Branch()
				if err != nil {
					return err
				}
				commit, _ := git.LatestCommit(branch)
				if commit != nil {
					opts.Description, err = git.CommitBody(strings.Trim(commit.Sha, `'`))
					if err != nil {
						return err
					}
					opts.Title = utils.Humanize(commit.Title)
				} else {
					opts.Title = utils.Humanize(branch)
				}
				_, err = api.GetCommit(apiClient, repo.FullName(), opts.TargetBranch)
				if err != nil {
					return fmt.Errorf("target branch %s does not exist on remote. Specify target branch with --target-branch flag",
						opts.TargetBranch)
				}
				if c, err := git.UncommittedChangeCount(); c != 0 {
					if err != nil {
						return err
					}
					fmt.Fprintf(f.IO.StdErr, "\nwarning: you have %s\n", utils.Pluralize(c, "uncommitted change"))
				}

				err = git.Push(repoRemote.PushURL.String(), opts.SourceBranch)
				if err != nil {
					return err
				}
			} else if opts.IsInteractive {
				var templateName string
				var templateContents string
				if opts.Description == "" {
					if opts.NoEditor {
						err = prompt.AskMultiline(&opts.Description, "Description:", "")
						if err != nil {
							return err
						}
					} else {
						templateResponse := struct {
							Index int
						}{}
						templateNames, err := cmdutils.ListGitLabTemplates(cmdutils.MergeRequestTemplate)
						if err != nil {
							return fmt.Errorf("error getting templates: %w", err)
						}

						templateNames = append(templateNames, "Open a blank merge request")

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
							templateContents, err = cmdutils.LoadGitLabTemplate(cmdutils.MergeRequestTemplate, templateName)
							if err != nil {
								return fmt.Errorf("failed to get template contents: %w", err)
							}
						}
					}
				}

				if opts.Title == "" {
					err = prompt.AskQuestionWithInput(&opts.Title, "Title:", "", true)
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

			if opts.IsDraft || opts.IsWIP {
				if opts.IsDraft {
					opts.Title = "Draft: " + opts.Title
				} else {
					opts.Title = "WIP: " + opts.Title
				}
			}
			mrCreateOpts.Title = gitlab.String(opts.Title)
			mrCreateOpts.Description = gitlab.String(opts.Description)
			mrCreateOpts.Labels = gitlab.Labels{opts.Labels}
			mrCreateOpts.SourceBranch = gitlab.String(opts.SourceBranch)
			mrCreateOpts.TargetBranch = gitlab.String(opts.TargetBranch)
			if opts.MileStone != -1 {
				mrCreateOpts.MilestoneID = gitlab.Int(opts.MileStone)
			}
			if opts.AllowCollaboration {
				mrCreateOpts.AllowCollaboration = gitlab.Bool(true)
			}
			if opts.RemoveSourceBranch {
				mrCreateOpts.RemoveSourceBranch = gitlab.Bool(true)
			}
			if opts.TargetProject != -1 {
				mrCreateOpts.TargetProjectID = gitlab.Int(opts.TargetProject)
			}
			if opts.Assignees != "" {
				arrIds := strings.Split(strings.Trim(opts.Assignees, "[] "), ",")
				var assigneeIDs []int

				for _, id := range arrIds {
					assigneeIDs = append(assigneeIDs, utils.StringToInt(id))
				}
				mrCreateOpts.AssigneeIDs = assigneeIDs
			}

			if opts.CreateSourceBranch {
				lb := &gitlab.CreateBranchOptions{
					Branch: gitlab.String(opts.SourceBranch),
					Ref:    gitlab.String(opts.TargetBranch),
				}
				fmt.Fprintln(f.IO.StdErr, "\nCreating related branch...")
				branch, err := api.CreateBranch(apiClient, repo.FullName(), lb)
				if err == nil {
					fmt.Fprintln(f.IO.StdErr, "Branch created: ", branch.WebURL)
				} else {
					fmt.Fprintln(f.IO.StdErr, "Error creating branch: ", err.Error())
				}
			}

			if opts.ShouldPush {
				err = git.Push(repoRemote.PushURL.String(), opts.SourceBranch)
				if err != nil {
					return err
				}
			}

			message := "\nCreating merge request for %s into %s in %s\n\n"
			if opts.IsDraft || opts.IsWIP {
				message = "\nCreating draft merge request for %s into %s in %s\n\n"
			}

			fmt.Fprintf(f.IO.StdErr, message, utils.Cyan(opts.SourceBranch), utils.Cyan(opts.TargetBranch), repo.FullName())

			mr, err := api.CreateMR(apiClient, repo.FullName(), mrCreateOpts)
			if err != nil {
				return err
			}

			fmt.Fprintln(out, mrutils.DisplayMR(mr))
			return nil
		},
	}
	mrCreateCmd.Flags().BoolVarP(&opts.Autofill, "fill", "f", false, "Do not prompt for title/description and just use commit info")
	mrCreateCmd.Flags().BoolVarP(&opts.IsDraft, "draft", "", false, "Mark merge request as a draft")
	mrCreateCmd.Flags().BoolVarP(&opts.IsWIP, "wip", "", false, "Mark merge request as a work in progress. Alternative to --draft")
	mrCreateCmd.Flags().BoolVarP(&opts.ShouldPush, "push", "", false, "Push committed changes after creating merge request. Make sure you have committed changes")
	mrCreateCmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Supply a title for merge request")
	mrCreateCmd.Flags().StringVarP(&opts.Description, "description", "d", "", "Supply a description for merge request")
	mrCreateCmd.Flags().StringVarP(&opts.Labels, "label", "l", "", "Add label by name. Multiple labels should be comma separated")
	mrCreateCmd.Flags().StringVarP(&opts.Assignees, "assignee", "a", "", "Assign merge request to people by their IDs. Multiple values should be comma separated ")
	mrCreateCmd.Flags().StringVarP(&opts.SourceBranch, "source-branch", "s", "", "The Branch you are creating the merge request. Default is the current branch.")
	mrCreateCmd.Flags().StringVarP(&opts.TargetBranch, "target-branch", "b", "", "The target or base branch into which you want your code merged")
	mrCreateCmd.Flags().IntVarP(&opts.TargetProject, "target-project", "", -1, "Add target project by id")
	mrCreateCmd.Flags().BoolVarP(&opts.CreateSourceBranch, "create-source-branch", "", false, "Create source branch if it does not exist")
	mrCreateCmd.Flags().IntVarP(&opts.MileStone, "milestone", "m", -1, "add milestone by <id> for merge request")
	mrCreateCmd.Flags().BoolVarP(&opts.AllowCollaboration, "allow-collaboration", "", false, "Allow commits from other members")
	mrCreateCmd.Flags().BoolVarP(&opts.RemoveSourceBranch, "remove-source-branch", "", false, "Remove Source Branch on merge")
	mrCreateCmd.Flags().BoolVarP(&opts.NoEditor, "no-editor", "", false, "Don't open editor to enter description. If set to true, uses prompt. Default is false")

	return mrCreateCmd
}

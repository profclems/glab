package create

import (
	"errors"
	"fmt"
	"strings"

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glrepo"

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
	Title                string
	Description          string
	SourceBranch         string
	TargetBranch         string
	TargetTrackingBranch string
	Labels               string
	Assignees            []string

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

	IO       *utils.IOStreams
	Branch   func() (string, error)
	Remotes  func() (glrepo.Remotes, error)
	Lab      func() (*gitlab.Client, error)
	Config   func() (config.Config, error)
	BaseRepo func() (glrepo.Interface, error)
}

func NewCmdCreate(f *cmdutils.Factory) *cobra.Command {
	opts := &CreateOpts{
		IO:       f.IO,
		Branch:   f.Branch,
		Remotes:  f.Remotes,
		Lab:      f.HttpClient,
		Config:   f.Config,
		BaseRepo: f.BaseRepo,
	}

	var mrCreateCmd = &cobra.Command{
		Use:     "create",
		Short:   `Create new merge request`,
		Long:    ``,
		Aliases: []string{"new"},
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			out := opts.IO.StdOut
			mrCreateOpts := &gitlab.CreateMergeRequestOptions{}

			hasTitle := cmd.Flags().Changed("title")
			hasDescription := cmd.Flags().Changed("description")

			// disable interactive mode if title and description are explicitly defined
			opts.IsInteractive = !(hasTitle && hasDescription)

			if opts.IsInteractive && !opts.IO.PromptEnabled() && !opts.Autofill {
				return &cmdutils.FlagError{Err: errors.New("--title or --fill required for non-interactive mode")}
			}

			labClient, err := opts.Lab()
			if err != nil {
				return err
			}

			repo, err := opts.BaseRepo()
			if err != nil {
				return err
			}

			remotes, err := opts.Remotes()
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
			opts.TargetTrackingBranch = fmt.Sprintf("%s/%s", repoRemote.Name, opts.TargetBranch)

			if opts.CreateSourceBranch && opts.SourceBranch == "" {
				opts.SourceBranch = utils.ReplaceNonAlphaNumericChars(opts.Title, "-")
			} else if opts.SourceBranch == "" {
				opts.SourceBranch, err = opts.Branch()
				if err != nil {
					return err
				}
			}

			if opts.Autofill {
				if err = mrBodyAndTitle(opts); err != nil {
					return err
				}
				_, err = api.GetCommit(labClient, repo.FullName(), opts.TargetBranch)
				if err != nil {
					return fmt.Errorf("target branch %s does not exist on remote. Specify target branch with --target-branch flag",
						opts.TargetBranch)
				}
				if c, err := git.UncommittedChangeCount(); c != 0 {
					if err != nil {
						return err
					}
					fmt.Fprintf(opts.IO.StdErr, "\nwarning: you have %s\n", utils.Pluralize(c, "uncommitted change"))
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
						editor, err := cmdutils.GetEditor(opts.Config)
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
					err = cmdutils.LabelsPrompt(&opts.Labels, labClient, repoRemote)
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
			if len(opts.Assignees) > 0 {
				users, err := api.UsersByNames(labClient, opts.Assignees)
				if err != nil {
					return err
				}
				mrCreateOpts.AssigneeIDs = cmdutils.IDsFromUsers(users)
			}

			if opts.CreateSourceBranch {
				lb := &gitlab.CreateBranchOptions{
					Branch: gitlab.String(opts.SourceBranch),
					Ref:    gitlab.String(opts.TargetBranch),
				}
				fmt.Fprintln(opts.IO.StdErr, "\nCreating related branch...")
				branch, err := api.CreateBranch(labClient, repo.FullName(), lb)
				if err == nil {
					fmt.Fprintln(opts.IO.StdErr, "Branch created: ", branch.WebURL)
				} else {
					fmt.Fprintln(opts.IO.StdErr, "Error creating branch: ", err.Error())
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

			fmt.Fprintf(opts.IO.StdErr, message, utils.Cyan(opts.SourceBranch), utils.Cyan(opts.TargetBranch), repo.FullName())

			mr, err := api.CreateMR(labClient, repo.FullName(), mrCreateOpts)
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
	mrCreateCmd.Flags().StringSliceVarP(&opts.Assignees, "assignee", "a", []string{}, "Assign merge request to people by their `usernames`")
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

func mrBodyAndTitle(opts *CreateOpts) error {
	// TODO: detect forks
	commits, err := git.Commits(opts.TargetTrackingBranch, opts.SourceBranch)
	if err != nil {
		return err
	}
	if len(commits) == 1 {
		opts.Title = commits[0].Title
		body, err := git.CommitBody(commits[0].Sha)
		if err != nil {
			return err
		}
		opts.Description = body
	} else {
		opts.Title = utils.Humanize(opts.SourceBranch)

		var body strings.Builder
		for i := len(commits) - 1; i >= 0; i-- {
			fmt.Fprintf(&body, "- %s\n", commits[i].Title)
		}
		opts.Description = body.String()
	}
	return nil
}

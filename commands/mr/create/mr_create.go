package create

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/profclems/glab/pkg/surveyext"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/profclems/glab/pkg/prompt"

	"strings"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type CreateOpts struct {
	Title              string
	Description        string
	MileStone          int
	SourceBranch       string
	TargetBranch       string
	Labels             string
	Assignees          string
	TargetProject      int
	CreateSourceBranch bool
	RemoveSourceBranch bool
	AllowCollaboration bool

	Autofill   bool
	IsDraft    bool
	IsWIP      bool
	ShouldPush bool
	NoEditor   bool
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
			} else {
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
					fmt.Fprintf(out, "warning: you have %s\n", utils.Pluralize(c, "uncommitted changes"))
				}

				err = git.Push(repoRemote.PushURL.String(), opts.SourceBranch)
				if err != nil {
					return err
				}
			} else {
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
						err = DescriptionPrompt(opts, templateContents, editor)
					}
				}
			}

			if opts.IsDraft || opts.IsWIP {
				if opts.IsDraft {
					opts.Title = "Draft: " + opts.Title
				} else {
					opts.Title = "WIP: " + opts.Title
				}
			}
			mergeLabel, _ := cmd.Flags().GetString("label")
			mrCreateOpts.Title = gitlab.String(opts.Title)
			mrCreateOpts.Description = gitlab.String(opts.Description)
			mrCreateOpts.Labels = gitlab.Labels{mergeLabel}
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
				var t2 []int

				for _, i := range arrIds {
					j := utils.StringToInt(i)
					t2 = append(t2, j)
				}
				mrCreateOpts.AssigneeIDs = t2
			}

			if opts.CreateSourceBranch {
				lb := &gitlab.CreateBranchOptions{
					Branch: gitlab.String(opts.SourceBranch),
					Ref:    gitlab.String(opts.TargetBranch),
				}
				fmt.Fprintln(out, "Creating related branch...")
				branch, err := api.CreateBranch(apiClient, repo.FullName(), lb)
				if err == nil {
					fmt.Fprintln(out, "Branch created: ", branch.WebURL)
				} else {
					fmt.Fprintln(out, "Error creating branch: ", err)
				}
			}

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

func DescriptionPrompt(mrOpts *CreateOpts, templateContent, editorCommand string) error {
	if templateContent != "" {
		if mrOpts.Description != "" {
			// prevent excessive newlines between default body and template
			mrOpts.Description = strings.TrimRight(mrOpts.Description, "\n")
			mrOpts.Description += "\n\n"
		}
		mrOpts.Description += templateContent
	}

	qs := []*survey.Question{
		{
			Name: "Body",
			Prompt: &surveyext.GLabEditor{
				BlankAllowed:  true,
				EditorCommand: editorCommand,
				Editor: &survey.Editor{
					Message:       "Body",
					FileName:      "*.md",
					Default:       mrOpts.Description,
					HideDefault:   true,
					AppendDefault: true,
				},
			},
		},
	}

	err := prompt.Ask(qs, mrOpts)
	if err != nil {
		return err
	}

	return nil
}

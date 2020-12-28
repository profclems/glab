package create

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/internal/utils"

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
	Labels      []string
	Assignees   []string

	Weight    int
	MileStone int
	LinkedMR  int

	NoEditor       bool
	IsConfidential bool
	IsInteractive  bool
	OpenInWeb      bool

	IO         *utils.IOStreams
	BaseRepo   func() (glrepo.Interface, error)
	HTTPClient func() (*gitlab.Client, error)
	Remotes    func() (glrepo.Remotes, error)
	Config     func() (config.Config, error)

	BaseProject *gitlab.Project
}

func NewCmdCreate(f *cmdutils.Factory) *cobra.Command {
	opts := &CreateOpts{
		IO:         f.IO,
		BaseRepo:   f.BaseRepo,
		HTTPClient: f.HttpClient,
		Remotes:    f.Remotes,
		Config:     f.Config,
	}
	var issueCreateCmd = &cobra.Command{
		Use:     "create [flags]",
		Short:   `Create an issue`,
		Long:    ``,
		Aliases: []string{"new"},
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			apiClient, err := opts.HTTPClient()
			if err != nil {
				return err
			}

			repo, err := opts.BaseRepo()
			if err != nil {
				return err
			}

			hasTitle := cmd.Flags().Changed("title")
			hasDescription := cmd.Flags().Changed("description")

			// disable interactive mode if title and description are explicitly defined
			opts.IsInteractive = !(hasTitle && hasDescription)

			if opts.IsInteractive && !opts.IO.PromptEnabled() {
				return &cmdutils.FlagError{Err: errors.New("--title and --description required for non-interactive mode")}
			}

			opts.BaseProject, err = api.GetProject(apiClient, repo.FullName())
			if err != nil {
				return err
			}

			if !opts.BaseProject.IssuesEnabled {
				fmt.Fprintf(opts.IO.StdErr, "Issues are disabled for %q\n", opts.BaseProject.PathWithNamespace)
				return cmdutils.SilentError
			}

			return createRun(opts)
		},
	}
	issueCreateCmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Supply a title for issue")
	issueCreateCmd.Flags().StringVarP(&opts.Description, "description", "d", "", "Supply a description for issue")
	issueCreateCmd.Flags().StringSliceVarP(&opts.Labels, "label", "l", []string{}, "Add label by name. Multiple labels should be comma separated")
	issueCreateCmd.Flags().StringSliceVarP(&opts.Assignees, "assignee", "a", []string{}, "Assign issue to people by their `usernames`")
	issueCreateCmd.Flags().IntVarP(&opts.MileStone, "milestone", "m", 0, "The global ID of a milestone to assign issue")
	issueCreateCmd.Flags().BoolVarP(&opts.IsConfidential, "confidential", "c", false, "Set an issue to be confidential. Default is false")
	issueCreateCmd.Flags().IntVarP(&opts.LinkedMR, "linked-mr", "", 0, "The IID of a merge request in which to resolve all issues")
	issueCreateCmd.Flags().IntVarP(&opts.Weight, "weight", "w", 0, "The weight of the issue. Valid values are greater than or equal to 0.")
	issueCreateCmd.Flags().BoolVarP(&opts.NoEditor, "no-editor", "", false, "Don't open editor to enter description. If set to true, uses prompt. Default is false")

	return issueCreateCmd
}

func createRun(opts *CreateOpts) error {
	apiClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}

	repo, err := opts.BaseRepo()
	if err != nil {
		return err
	}

	var templateName string
	var templateContents string

	issueCreateOpts := &gitlab.CreateIssueOptions{}

	if opts.IsInteractive {
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
		if len(opts.Labels) == 0 {
			remotes, err := opts.Remotes()
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
	} else if opts.Title == "" {
		return fmt.Errorf("title can't be blank")
	}

	var action cmdutils.Action

	// submit without prompting for non interactive mode
	if !opts.IsInteractive {
		action = cmdutils.SubmitAction
	}

	if action == 0 {
		action, err = cmdutils.ConfirmSubmission(true)
		if err != nil {
			return fmt.Errorf("unable to confirm: %w", err)
		}
	}

	if action == cmdutils.CancelAction {
		fmt.Fprintln(opts.IO.StdErr, "Discarded.")
		return nil
	}

	if action == cmdutils.PreviewAction {
		return previewIssue(opts)
	}

	if action == cmdutils.SubmitAction {
		issueCreateOpts.Title = gitlab.String(opts.Title)
		issueCreateOpts.Labels = opts.Labels
		issueCreateOpts.Description = &opts.Description
		if opts.IsConfidential {
			issueCreateOpts.Confidential = gitlab.Bool(opts.IsConfidential)
		}
		if opts.Weight != 0 {
			issueCreateOpts.Weight = gitlab.Int(opts.Weight)
		}
		if opts.LinkedMR != 0 {
			issueCreateOpts.MergeRequestToResolveDiscussionsOf = gitlab.Int(opts.LinkedMR)
		}
		if opts.MileStone != 0 {
			issueCreateOpts.MilestoneID = gitlab.Int(opts.MileStone)
		}
		if len(opts.Assignees) > 0 {
			users, err := api.UsersByNames(apiClient, opts.Assignees)
			if err != nil {
				return err
			}
			issueCreateOpts.AssigneeIDs = cmdutils.IDsFromUsers(users)
		}
		fmt.Fprintln(opts.IO.StdErr, "\n- Creating issue in", repo.FullName())
		issue, err := api.CreateIssue(apiClient, repo.FullName(), issueCreateOpts)
		if err != nil {
			return err
		}
		fmt.Fprintln(opts.IO.StdOut, issueutils.DisplayIssue(issue))
		return nil
	}

	return errors.New("expected to cancel, preview in browser, or submit")
}

func previewIssue(opts *CreateOpts) error {
	repo, err := opts.BaseRepo()
	if err != nil {
		return err
	}

	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	openURL, err := generateIssueWebURL(opts, repo)
	if err != nil {
		return err
	}

	if opts.IO.IsOutputTTY() {
		fmt.Fprintf(opts.IO.StdErr, "Opening %s in your browser.\n", utils.DisplayURL(openURL))
	}
	browser, _ := cfg.Get(repo.RepoHost(), "browser")
	return utils.OpenInBrowser(openURL, browser)
}

func generateIssueWebURL(opts *CreateOpts, repo glrepo.Interface) (string, error) {
	description := opts.Description

	if len(opts.Labels) > 0 {
		// this uses the slash commands to add labels to the description
		// See https://docs.gitlab.com/ee/user/project/quick_actions.html
		// See also https://gitlab.com/gitlab-org/gitlab-foss/-/issues/19731#note_32550046
		description += "\n/label "
		for _, label := range opts.Labels {
			description += fmt.Sprintf("~%q", label)
		}
	}
	if len(opts.Assignees) > 0 {
		// this uses the slash commands to add assignees to the description
		description += fmt.Sprintf("\n/assign %s", strings.Join(opts.Assignees, ", "))
	}
	if opts.MileStone != 0 {
		// this uses the slash commands to add milestone to the description
		description += fmt.Sprintf("\n/milestone %%%d", opts.MileStone)
	}
	if opts.Weight != 0 {
		// this uses the slash commands to add weight to the description
		description += fmt.Sprintf("\n/weight %%%d", opts.Weight)
	}
	if opts.IsConfidential {
		// this uses the slash commands to add confidential to the description
		description += "\n/confidential"
	}

	u, err := url.Parse(opts.BaseProject.WebURL)
	if err != nil {
		return "", err
	}
	u.Path += "/-/issues/new"
	u.RawQuery = fmt.Sprintf(
		"utf8=✓&issue[title]=%s&issue[description]=%s",
		opts.Title,
		description)
	return u.String(), nil
}

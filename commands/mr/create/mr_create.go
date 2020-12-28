package create

import (
	"errors"
	"fmt"
	"net/url"
	"os"
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
	Labels               []string
	Assignees            []string
	MileStone            int

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
	HeadRepo func() (glrepo.Interface, error)

	// SourceProject is the Project we create the merge request in and where we push our branch
	// it is the project we have permission to push so most likely one's fork
	SourceProject *gitlab.Project
	// TargetProject is the one we query for changes between our branch and the target branch
	// it is the one we merge request will appear in
	TargetProject *gitlab.Project
}

func NewCmdCreate(f *cmdutils.Factory) *cobra.Command {
	var mrCreateTargetProject string

	opts := &CreateOpts{
		IO:       f.IO,
		Branch:   f.Branch,
		Remotes:  f.Remotes,
		Lab:      f.HttpClient,
		Config:   f.Config,
		BaseRepo: f.BaseRepo,
		HeadRepo: resolvedHeadRepo(f),
	}

	var mrCreateCmd = &cobra.Command{
		Use:     "create",
		Short:   `Create new merge request`,
		Long:    ``,
		Aliases: []string{"new"},
		Args:    cobra.ExactArgs(0),
		PreRun: func(cmd *cobra.Command, args []string) {
			repoOverride, _ := cmd.Flags().GetString("head")
			if repoFromEnv := os.Getenv("GITLAB_HEAD_REPO"); repoOverride == "" && repoFromEnv != "" {
				repoOverride = repoFromEnv
			}
			if repoOverride != "" {
				_ = headRepoOverride(opts, repoOverride)
			}
		},
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
			if cmd.Flags().Changed("wip") && cmd.Flags().Changed("draft") {
				return &cmdutils.FlagError{Err: errors.New("specify either of --draft or --wip")}
			}

			labClient, err := opts.Lab()
			if err != nil {
				return err
			}

			baseRepo, err := opts.BaseRepo()
			if err != nil {
				return err
			}

			headRepo, err := opts.HeadRepo()
			if err != nil {
				return err
			}

			opts.SourceProject, err = api.GetProject(labClient, headRepo.FullName())
			if err != nil {
				return err
			}

			// if the user set the target_project, get details of the target
			if mrCreateTargetProject != "" {
				opts.TargetProject, err = api.GetProject(labClient, mrCreateTargetProject)
				if err != nil {
					return err
				}
			} else {
				// If both the baseRepo and headRepo are the same then re-use the SourceProject
				if baseRepo.FullName() == headRepo.FullName() {
					opts.TargetProject = opts.SourceProject
				} else {
					// Otherwise assume the user wants to create the merge request against the
					// baseRepo
					opts.TargetProject, err = api.GetProject(labClient, baseRepo.FullName())
					if err != nil {
						return err
					}
				}
			}

			if !opts.TargetProject.MergeRequestsEnabled {
				fmt.Fprintf(opts.IO.StdErr, "Merge requests are disabled for %q\n", opts.TargetProject.PathWithNamespace)
				return cmdutils.SilentError
			}

			headRepoRemote, err := repoRemote(labClient, opts, headRepo, opts.SourceProject, "glab-head")
			if err != nil {
				return nil
			}

			var baseRepoRemote *glrepo.Remote

			// check if baseRepo is the same as the headRepo and set the remote
			if glrepo.IsSame(baseRepo, headRepo) {
				baseRepoRemote = headRepoRemote
			} else {
				baseRepoRemote, err = repoRemote(labClient, opts, baseRepo, opts.TargetProject, "glab-base")
				if err != nil {
					return nil
				}
			}

			if opts.TargetBranch == "" {
				opts.TargetBranch, _ = git.GetDefaultBranch(baseRepoRemote.PushURL.String())
			}
			opts.TargetTrackingBranch = fmt.Sprintf("%s/%s", baseRepoRemote.Name, opts.TargetBranch)

			if opts.CreateSourceBranch && opts.SourceBranch == "" {
				opts.SourceBranch = utils.ReplaceNonAlphaNumericChars(opts.Title, "-")
			} else if opts.SourceBranch == "" {
				opts.SourceBranch, err = opts.Branch()
				if err != nil {
					return err
				}
			}

			if opts.SourceBranch == opts.TargetBranch && glrepo.IsSame(baseRepo, headRepo) {
				fmt.Fprintf(opts.IO.StdErr, "You must be on a different branch other than %q\n", opts.TargetBranch)
				return cmdutils.SilentError
			}

			if opts.Autofill {
				if err = mrBodyAndTitle(opts); err != nil {
					return err
				}
				_, err = api.GetCommit(labClient, baseRepo.FullName(), opts.TargetBranch)
				if err != nil {
					return fmt.Errorf("target branch %s does not exist on remote. Specify target branch with --target-branch flag",
						opts.TargetBranch)
				}

				opts.ShouldPush = true
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
				if len(opts.Labels) == 0 {
					err = cmdutils.LabelsPrompt(&opts.Labels, labClient, baseRepoRemote)
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
			mrCreateOpts.Labels = opts.Labels
			mrCreateOpts.SourceBranch = gitlab.String(opts.SourceBranch)
			mrCreateOpts.TargetBranch = gitlab.String(opts.TargetBranch)
			if opts.MileStone != 0 {
				mrCreateOpts.MilestoneID = gitlab.Int(opts.MileStone)
			}
			if opts.AllowCollaboration {
				mrCreateOpts.AllowCollaboration = gitlab.Bool(true)
			}
			if opts.RemoveSourceBranch {
				mrCreateOpts.RemoveSourceBranch = gitlab.Bool(true)
			}
			if opts.TargetProject != nil {
				mrCreateOpts.TargetProjectID = gitlab.Int(opts.TargetProject.ID)
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
				branch, err := api.CreateBranch(labClient, headRepo.FullName(), lb)
				if err == nil {
					fmt.Fprintln(opts.IO.StdErr, "Branch created: ", branch.WebURL)
				} else {
					fmt.Fprintln(opts.IO.StdErr, "Error creating branch: ", err.Error())
				}
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

			if err := handlePush(opts, headRepoRemote); err != nil {
				return err
			}

			if action == cmdutils.PreviewAction {
				return previewMR(opts)
			}

			if action == cmdutils.SubmitAction {

				message := "\nCreating merge request for %s into %s in %s\n\n"
				if opts.IsDraft || opts.IsWIP {
					message = "\nCreating draft merge request for %s into %s in %s\n\n"
				}

				fmt.Fprintf(opts.IO.StdErr, message, utils.Cyan(opts.SourceBranch), utils.Cyan(opts.TargetBranch), baseRepo.FullName())

				// It is intentional that we create against the head repo, it is necessary
				// for cross-repository merge requests
				mr, err := api.CreateMR(labClient, headRepo.FullName(), mrCreateOpts)
				if err != nil {
					return err
				}

				fmt.Fprintln(out, mrutils.DisplayMR(mr))
				return nil
			}

			return errors.New("expected to cancel, preview in browser, or submit")
		},
	}
	mrCreateCmd.Flags().BoolVarP(&opts.Autofill, "fill", "f", false, "Do not prompt for title/description and just use commit info")
	mrCreateCmd.Flags().BoolVarP(&opts.IsDraft, "draft", "", false, "Mark merge request as a draft")
	mrCreateCmd.Flags().BoolVarP(&opts.IsWIP, "wip", "", false, "Mark merge request as a work in progress. Alternative to --draft")
	mrCreateCmd.Flags().BoolVarP(&opts.ShouldPush, "push", "", false, "Push committed changes after creating merge request. Make sure you have committed changes")
	mrCreateCmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Supply a title for merge request")
	mrCreateCmd.Flags().StringVarP(&opts.Description, "description", "d", "", "Supply a description for merge request")
	mrCreateCmd.Flags().StringSliceVarP(&opts.Labels, "label", "l", []string{}, "Add label by name. Multiple labels should be comma separated")
	mrCreateCmd.Flags().StringSliceVarP(&opts.Assignees, "assignee", "a", []string{}, "Assign merge request to people by their `usernames`")
	mrCreateCmd.Flags().StringVarP(&opts.SourceBranch, "source-branch", "s", "", "The Branch you are creating the merge request. Default is the current branch.")
	mrCreateCmd.Flags().StringVarP(&opts.TargetBranch, "target-branch", "b", "", "The target or base branch into which you want your code merged")
	mrCreateCmd.Flags().BoolVarP(&opts.CreateSourceBranch, "create-source-branch", "", false, "Create source branch if it does not exist")
	mrCreateCmd.Flags().IntVarP(&opts.MileStone, "milestone", "m", 0, "add milestone by <id> for merge request")
	mrCreateCmd.Flags().BoolVarP(&opts.AllowCollaboration, "allow-collaboration", "", false, "Allow commits from other members")
	mrCreateCmd.Flags().BoolVarP(&opts.RemoveSourceBranch, "remove-source-branch", "", false, "Remove Source Branch on merge")
	mrCreateCmd.Flags().BoolVarP(&opts.NoEditor, "no-editor", "", false, "Don't open editor to enter description. If set to true, uses prompt. Default is false")
	mrCreateCmd.Flags().StringP("head", "H", "", "Select another head repository using the `OWNER/REPO` or `GROUP/NAMESPACE/REPO` format or the project ID or full URL")

	mrCreateCmd.Flags().StringVarP(&mrCreateTargetProject, "target-project", "", "", "Add target project by id or OWNER/REPO or GROUP/NAMESPACE/REPO")
	mrCreateCmd.Flags().MarkHidden("target-project")
	mrCreateCmd.Flags().MarkDeprecated("target-project", "Use --repo instead")

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

func handlePush(opts *CreateOpts, remote *glrepo.Remote) error {
	if opts.ShouldPush {
		var sourceRemote = remote

		sourceBranch := opts.SourceBranch

		if sourceBranch != "" {
			if idx := strings.IndexRune(sourceBranch, ':'); idx >= 0 {
				sourceBranch = sourceBranch[idx+1:]
			}
		}

		if c, err := git.UncommittedChangeCount(); c != 0 {
			if err != nil {
				return err
			}
			fmt.Fprintf(opts.IO.StdErr, "\nwarning: you have %s\n", utils.Pluralize(c, "uncommitted change"))
		}
		err := git.Push(sourceRemote.Name, fmt.Sprintf("HEAD:%s", sourceBranch), opts.IO.StdOut, opts.IO.StdErr)
		if err == nil {
			branchConfig := git.ReadBranchConfig(sourceBranch)
			if branchConfig.RemoteName == "" && (branchConfig.MergeRef == "" || branchConfig.RemoteURL == nil) {
				// No remote is set so set it
				_ = git.SetUpstream(sourceRemote.Name, sourceBranch, opts.IO.StdOut, opts.IO.StdErr)
			}
		}
		return err
	}

	return nil
}

func previewMR(opts *CreateOpts) error {
	repo, err := opts.BaseRepo()
	if err != nil {
		return err
	}

	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	openURL, err := generateMRCompareURL(opts)
	if err != nil {
		return err
	}

	if opts.IO.IsOutputTTY() {
		fmt.Fprintf(opts.IO.StdErr, "Opening %s in your browser.\n", utils.DisplayURL(openURL))
	}
	browser, _ := cfg.Get(repo.RepoHost(), "browser")
	return utils.OpenInBrowser(openURL, browser)
}

func generateMRCompareURL(opts *CreateOpts) (string, error) {
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

	// The merge request **must** be opened against the head repo
	u, err := url.Parse(opts.SourceProject.WebURL)
	if err != nil {
		return "", err
	}
	u.Path += "/-/merge_requests/new"
	u.RawQuery = fmt.Sprintf(
		"utf8=✓&merge_request[title]=%s&merge_request[description]=%s&merge_request[source_branch]=%s&merge_request[target_branch]=%s&merge_request[source_project_id]=%d&merge_request[target_project_id]=%d",
		opts.Title,
		description,
		opts.SourceBranch,
		opts.TargetBranch,
		opts.SourceProject.ID,
		opts.TargetProject.ID)
	return u.String(), nil
}

func resolvedHeadRepo(f *cmdutils.Factory) func() (glrepo.Interface, error) {
	return func() (glrepo.Interface, error) {
		httpClient, err := f.HttpClient()
		if err != nil {
			return nil, err
		}
		remotes, err := f.Remotes()
		if err != nil {
			return nil, err
		}
		repoContext, err := glrepo.ResolveRemotesToRepos(remotes, httpClient, "")
		if err != nil {
			return nil, err
		}
		headRepo, err := repoContext.HeadRepo(true)
		if err != nil {
			return nil, err
		}

		return headRepo, nil
	}
}

func headRepoOverride(opts *CreateOpts, repo string) error {
	opts.HeadRepo = func() (glrepo.Interface, error) {
		return glrepo.FromFullName(repo)
	}
	return nil
}

func repoRemote(labClient *gitlab.Client, opts *CreateOpts, repo glrepo.Interface, project *gitlab.Project, remoteName string) (*glrepo.Remote, error) {
	remotes, err := opts.Remotes()
	if err != nil {
		return nil, err
	}
	repoRemote, _ := remotes.FindByRepo(repo.RepoOwner(), repo.RepoName())
	if repoRemote == nil {
		cfg, err := opts.Config()
		if err != nil {
			return nil, err
		}
		gitProtocol, _ := cfg.Get(repo.RepoHost(), "git_protocol")
		token, _ := cfg.Get(repo.RepoHost(), "token")
		repoURL := project.SSHURLToRepo

		if gitProtocol != "ssh" {
			u, err := api.CurrentUser(labClient)
			if err != nil {
				return nil, fmt.Errorf("failed to get current user info: %q", err)
			}
			remoteArgs := &glrepo.RemoteArgs{
				Protocol: gitProtocol,
				Token:    token,
				Url:      repo.RepoHost(),
				Username: u.Username,
			}

			repoURL, _ = glrepo.RemoteURL(project, remoteArgs)

		}

		gitRemote, err := git.AddRemote(remoteName, repoURL)
		if err != nil {
			return nil, fmt.Errorf("error adding remote: %w", err)
		}
		repoRemote = &glrepo.Remote{
			Remote: gitRemote,
			Repo:   repo,
		}
	}

	return repoRemote, nil
}

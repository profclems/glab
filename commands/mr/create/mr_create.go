package create

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/pkg/prompt"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glrepo"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/pkg/git"
	"github.com/profclems/glab/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type CreateOpts struct {
	Title                 string
	Description           string
	SourceBranch          string
	TargetBranch          string
	TargetTrackingBranch  string
	Labels                []string
	Assignees             []string
	Reviewers             []string
	MileStone             int
	MilestoneFlag         string
	MRCreateTargetProject string

	RelatedIssue    string
	CopyIssueLabels bool

	CreateSourceBranch bool
	RemoveSourceBranch bool
	AllowCollaboration bool
	SquashBeforeMerge  bool

	Autofill       bool
	FillCommitBody bool
	IsDraft        bool
	IsWIP          bool
	ShouldPush     bool
	NoEditor       bool
	IsInteractive  bool
	Yes            bool
	Web            bool

	IO       *iostreams.IOStreams
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

func NewCmdCreate(f *cmdutils.Factory, runE func(opts *CreateOpts) error) *cobra.Command {
	opts := &CreateOpts{
		IO:       f.IO,
		Branch:   f.Branch,
		Remotes:  f.Remotes,
		Config:   f.Config,
		HeadRepo: resolvedHeadRepo(f),
	}

	var mrCreateCmd = &cobra.Command{
		Use:     "create",
		Short:   `Create new merge request`,
		Long:    ``,
		Aliases: []string{"new"},
		Example: heredoc.Doc(`
			$ glab mr new
			$ glab mr create -a username -t "fix annoying bug"
			$ glab mr create -f --draft --label RFC
			$ glab mr create --fill --yes --web
			$ glab mr create --fill --fill-commit-body --yes
		`),
		Args: cobra.ExactArgs(0),
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
			// support `-R, --repo` override
			//
			// NOTE: it is important to assign the BaseRepo and HTTPClient in RunE because
			// they are overridden in a PersistentRun hook (when `-R, --repo` is specified)
			// which runs before RunE is executed
			opts.BaseRepo = f.BaseRepo
			opts.Lab = f.HttpClient

			hasTitle := cmd.Flags().Changed("title")
			hasDescription := cmd.Flags().Changed("description")

			// disable interactive mode if title and description are explicitly defined
			opts.IsInteractive = !(hasTitle && hasDescription)

			if hasTitle && hasDescription && opts.Autofill {
				return &cmdutils.FlagError{
					Err: errors.New("usage of --title and --description completely override --fill"),
				}
			}
			if opts.IsInteractive && !opts.IO.PromptEnabled() && !opts.Autofill {
				return &cmdutils.FlagError{Err: errors.New("--title or --fill required for non-interactive mode")}
			}
			if cmd.Flags().Changed("wip") && cmd.Flags().Changed("draft") {
				return &cmdutils.FlagError{Err: errors.New("specify either of --draft or --wip")}
			}
			if !opts.Autofill && opts.FillCommitBody {
				return &cmdutils.FlagError{Err: errors.New("--fill-commit-body should be used with --fill")}
			}
			// Remove this once --yes does more than just skip the prompts that --web happen to skip
			// by design
			if opts.Yes && opts.Web {
				return &cmdutils.FlagError{Err: errors.New("--web already skips all prompts currently skipped by --yes")}
			}

			if opts.CopyIssueLabels && opts.RelatedIssue == "" {
				return &cmdutils.FlagError{Err: errors.New("--copy-issue-labels can only be used with --related-issue")}
			}

			if runE != nil {
				return runE(opts)
			}

			return createRun(opts)
		},
	}
	mrCreateCmd.Flags().BoolVarP(&opts.Autofill, "fill", "f", false, "Do not prompt for title/description and just use commit info")
	mrCreateCmd.Flags().BoolVarP(&opts.FillCommitBody, "fill-commit-body", "", false, "Fill description with each commit body when multiple commits. Can only be used with --fill")
	mrCreateCmd.Flags().BoolVarP(&opts.IsDraft, "draft", "", false, "Mark merge request as a draft")
	mrCreateCmd.Flags().BoolVarP(&opts.IsWIP, "wip", "", false, "Mark merge request as a work in progress. Alternative to --draft")
	mrCreateCmd.Flags().BoolVarP(&opts.ShouldPush, "push", "", false, "Push committed changes after creating merge request. Make sure you have committed changes")
	mrCreateCmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Supply a title for merge request")
	mrCreateCmd.Flags().StringVarP(&opts.Description, "description", "d", "", "Supply a description for merge request")
	mrCreateCmd.Flags().StringSliceVarP(&opts.Labels, "label", "l", []string{}, "Add label by name. Multiple labels should be comma separated")
	mrCreateCmd.Flags().StringSliceVarP(&opts.Assignees, "assignee", "a", []string{}, "Assign merge request to people by their `usernames`")
	mrCreateCmd.Flags().StringSliceVarP(&opts.Reviewers, "reviewer", "", []string{}, "Request review from users by their `usernames`")
	mrCreateCmd.Flags().StringVarP(&opts.SourceBranch, "source-branch", "s", "", "The Branch you are creating the merge request. Default is the current branch.")
	mrCreateCmd.Flags().StringVarP(&opts.TargetBranch, "target-branch", "b", "", "The target or base branch into which you want your code merged")
	mrCreateCmd.Flags().BoolVarP(&opts.CreateSourceBranch, "create-source-branch", "", false, "Create source branch if it does not exist")
	mrCreateCmd.Flags().StringVarP(&opts.MilestoneFlag, "milestone", "m", "", "The global ID or title of a milestone to assign")
	mrCreateCmd.Flags().BoolVarP(&opts.AllowCollaboration, "allow-collaboration", "", false, "Allow commits from other members")
	mrCreateCmd.Flags().BoolVarP(&opts.RemoveSourceBranch, "remove-source-branch", "", false, "Remove Source Branch on merge")
	mrCreateCmd.Flags().BoolVarP(&opts.SquashBeforeMerge, "squash-before-merge", "", false, "Squash commits into a single commit when merging")
	mrCreateCmd.Flags().BoolVarP(&opts.NoEditor, "no-editor", "", false, "Don't open editor to enter description. If set to true, uses prompt. Default is false")
	mrCreateCmd.Flags().StringP("head", "H", "", "Select another head repository using the `OWNER/REPO` or `GROUP/NAMESPACE/REPO` format or the project ID or full URL")
	mrCreateCmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Skip submission confirmation prompt, with --fill it skips all optional prompts")
	mrCreateCmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "continue merge request creation on web browser")
	mrCreateCmd.Flags().BoolVarP(&opts.CopyIssueLabels, "copy-issue-labels", "", false, "Copy labels from issue to the merge request. Used with --related-issue")
	mrCreateCmd.Flags().StringVarP(&opts.RelatedIssue, "related-issue", "i", "", "Create merge request for an issue. The merge request title will be created from the issue if --title is not provided.")

	mrCreateCmd.Flags().StringVarP(&opts.MRCreateTargetProject, "target-project", "", "", "Add target project by id or OWNER/REPO or GROUP/NAMESPACE/REPO")
	_ = mrCreateCmd.Flags().MarkHidden("target-project")
	_ = mrCreateCmd.Flags().MarkDeprecated("target-project", "Use --repo instead")

	return mrCreateCmd
}

func parseIssue(apiClient *gitlab.Client, opts *CreateOpts) (*gitlab.Issue, error) {
	issue, _, err := issueutils.IssueFromArg(apiClient, opts.BaseRepo, opts.RelatedIssue)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func createRun(opts *CreateOpts) error {
	out := opts.IO.StdOut
	c := opts.IO.Color()
	mrCreateOpts := &gitlab.CreateMergeRequestOptions{}

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
	if opts.MRCreateTargetProject != "" {
		opts.TargetProject, err = api.GetProject(labClient, opts.MRCreateTargetProject)
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

	headRepoRemote, err := repoRemote(opts, headRepo, opts.SourceProject, "glab-head")
	if err != nil {
		return nil
	}

	var baseRepoRemote *glrepo.Remote

	// check if baseRepo is the same as the headRepo and set the remote
	if glrepo.IsSame(baseRepo, headRepo) {
		baseRepoRemote = headRepoRemote
	} else {
		baseRepoRemote, err = repoRemote(opts, baseRepo, opts.TargetProject, "glab-base")
		if err != nil {
			return nil
		}
	}

	if opts.MilestoneFlag != "" {
		opts.MileStone, err = cmdutils.ParseMilestone(labClient, baseRepo, opts.MilestoneFlag)
		if err != nil {
			return err
		}
	}

	if opts.CreateSourceBranch && opts.SourceBranch == "" {
		opts.SourceBranch = utils.ReplaceNonAlphaNumericChars(opts.Title, "-")
	} else if opts.SourceBranch == "" && opts.RelatedIssue == "" {
		opts.SourceBranch, err = opts.Branch()
		if err != nil {
			return err
		}
	}

	if opts.TargetBranch == "" {
		opts.TargetBranch = getTargetBranch(baseRepoRemote)
	}

	if opts.RelatedIssue != "" {
		issue, err := parseIssue(labClient, opts)
		if err != nil {
			return err
		}

		if opts.CopyIssueLabels {
			*mrCreateOpts.Labels = issue.Labels
		}
		opts.Description = fmt.Sprintf("Closes #%d", issue.IID)
		opts.Title = fmt.Sprintf("Resolve \"%s\"", issue.Title)
		if !opts.IsDraft && !opts.IsWIP {
			opts.IsDraft = true
		}

		if opts.SourceBranch == "" {
			sourceBranch := fmt.Sprintf("%d-%s", issue.IID, utils.ReplaceNonAlphaNumericChars(strings.ToLower(issue.Title), "-"))
			branchOpts := &gitlab.CreateBranchOptions{
				Branch: &sourceBranch,
				Ref:    &opts.TargetBranch,
			}

			_, err = api.CreateBranch(labClient, baseRepo.FullName(), branchOpts)
			if err != nil {
				for branchErr, branchCount := err, 1; branchErr != nil; branchCount++ {
					sourceBranch = fmt.Sprintf("%d-%s-%d", issue.IID, strings.ReplaceAll(strings.ToLower(issue.Title), " ", "-"), branchCount)
					_, branchErr = api.CreateBranch(labClient, baseRepo.FullName(), branchOpts)
				}
			}
			opts.SourceBranch = sourceBranch
		}
	} else {
		opts.TargetTrackingBranch = fmt.Sprintf("%s/%s", baseRepoRemote.Name, opts.TargetBranch)
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
					err = prompt.AskMultiline(&opts.Description, "description", "Description:", "")
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
				err = prompt.AskQuestionWithInput(&opts.Title, "title", "Title:", "", true)
				if err != nil {
					return err
				}
			}
			if opts.Description == "" {
				if opts.NoEditor {
					err = prompt.AskMultiline(&opts.Description, "description", "Description:", "")
					if err != nil {
						return err
					}
				} else {
					editor, err := cmdutils.GetEditor(opts.Config)
					if err != nil {
						return err
					}
					err = cmdutils.EditorPrompt(&opts.Description, "Description", templateContents, editor)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	if opts.Title == "" {
		return fmt.Errorf("title can't be blank")
	}

	if opts.IsDraft || opts.IsWIP {
		if opts.IsDraft {
			opts.Title = "Draft: " + opts.Title
		} else {
			opts.Title = "WIP: " + opts.Title
		}
	}
	mrCreateOpts.Title = &opts.Title
	mrCreateOpts.Description = &opts.Description
	mrCreateOpts.SourceBranch = &opts.SourceBranch
	mrCreateOpts.TargetBranch = &opts.TargetBranch

	if opts.AllowCollaboration {
		mrCreateOpts.AllowCollaboration = gitlab.Bool(true)
	}

	if opts.RemoveSourceBranch {
		mrCreateOpts.RemoveSourceBranch = gitlab.Bool(true)
	}

	if opts.SquashBeforeMerge {
		mrCreateOpts.Squash = gitlab.Bool(true)
	}

	if opts.TargetProject != nil {
		mrCreateOpts.TargetProjectID = &opts.TargetProject.ID
	}

	if opts.CreateSourceBranch {
		lb := &gitlab.CreateBranchOptions{
			Branch: &opts.SourceBranch,
			Ref:    &opts.TargetBranch,
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
	if !opts.IsInteractive || opts.Yes {
		action = cmdutils.SubmitAction
	}

	if opts.Web {
		action = cmdutils.PreviewAction
	}

	if action == cmdutils.NoAction {
		action, err = cmdutils.ConfirmSubmission(true, true)
		if err != nil {
			return fmt.Errorf("unable to confirm: %w", err)
		}
	}

	if action == cmdutils.AddMetadataAction {
		var metadataActions []cmdutils.Action

		metadataActions, err = cmdutils.PickMetadata()
		if err != nil {
			return fmt.Errorf("failed to pick metadata to add: %w", err)
		}

		for _, x := range metadataActions {
			if x == cmdutils.AddLabelAction {
				err = cmdutils.LabelsPrompt(&opts.Labels, labClient, baseRepoRemote)
				if err != nil {
					return err
				}
			}
			if x == cmdutils.AddAssigneeAction {
				// Use minimum permission level 30 (Maintainer) as it is the minimum level
				// to accept a merge request
				err = cmdutils.AssigneesPrompt(&opts.Assignees, labClient, baseRepoRemote, opts.IO, 30)
				if err != nil {
					return err
				}
			}
			if x == cmdutils.AddMilestoneAction {
				err = cmdutils.MilestonesPrompt(&opts.MileStone, labClient, baseRepoRemote, opts.IO)
				if err != nil {
					return err
				}
			}
		}

		// Ask the user again but don't permit AddMetadata a second time
		action, err = cmdutils.ConfirmSubmission(true, false)
		if err != nil {
			return err
		}
	}

	// These actions need to be done here, after the `Add metadata` prompt because
	// they are metadata that can be modified by the prompt
	*mrCreateOpts.Labels = append(*mrCreateOpts.Labels, opts.Labels...)

	if len(opts.Assignees) > 0 {
		users, err := api.UsersByNames(labClient, opts.Assignees)
		if err != nil {
			return err
		}
		mrCreateOpts.AssigneeIDs = cmdutils.IDsFromUsers(users)
	}

	if len(opts.Reviewers) > 0 {
		users, err := api.UsersByNames(labClient, opts.Reviewers)
		if err != nil {
			return err
		}
		mrCreateOpts.ReviewerIDs = cmdutils.IDsFromUsers(users)
	}

	if opts.MileStone != 0 {
		mrCreateOpts.MilestoneID = gitlab.Int(opts.MileStone)
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

		fmt.Fprintf(opts.IO.StdErr, message, c.Cyan(opts.SourceBranch), c.Cyan(opts.TargetBranch), baseRepo.FullName())

		// It is intentional that we create against the head repo, it is necessary
		// for cross-repository merge requests
		mr, err := api.CreateMR(labClient, headRepo.FullName(), mrCreateOpts)
		if err != nil {
			return err
		}

		fmt.Fprintln(out, mrutils.DisplayMR(c, mr, opts.IO.IsaTTY))
		return nil
	}

	return errors.New("expected to cancel, preview in browser, or submit")
}

func mrBodyAndTitle(opts *CreateOpts) error {
	// TODO: detect forks
	commits, err := git.Commits(opts.TargetTrackingBranch, opts.SourceBranch)
	if err != nil {
		return err
	}
	if len(commits) == 1 {
		if opts.Title == "" {
			opts.Title = commits[0].Title
		}
		if opts.Description == "" {
			body, err := git.CommitBody(commits[0].Sha)
			if err != nil {
				return err
			}
			opts.Description = body
		}
	} else {
		if opts.Title == "" {
			opts.Title = utils.Humanize(opts.SourceBranch)
		}

		if opts.Description == "" {
			var body strings.Builder
			for i := len(commits) - 1; i >= 0; i-- {
				// adds 2 spaces for markdown line wrapping
				fmt.Fprintf(&body, "- %s  \n", commits[i].Title)

				if opts.FillCommitBody {
					commitBody, err := git.CommitBody(commits[i].Sha)
					if err != nil {
						return err
					}
					re := regexp.MustCompile(`\r?\n\n`)
					commitBody = re.ReplaceAllString(commitBody, "  \n")
					fmt.Fprintf(&body, "%s\n", commitBody)
				}
			}
			opts.Description = body.String()
		}
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
	if len(opts.Reviewers) > 0 {
		// this uses the slash commands to add reviewers to the description
		description += fmt.Sprintf("\n/reviewer %s", strings.Join(opts.Reviewers, ", "))
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
		"merge_request[title]=%s&merge_request[description]=%s&merge_request[source_branch]=%s&merge_request[target_branch]=%s&merge_request[source_project_id]=%d&merge_request[target_project_id]=%d",
		strings.ReplaceAll(url.PathEscape(opts.Title), "+", "%2B"),
		strings.ReplaceAll(url.PathEscape(description), "+", "%2B"),
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

func repoRemote(opts *CreateOpts, repo glrepo.Interface, project *gitlab.Project, remoteName string) (*glrepo.Remote, error) {
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
		repoURL := glrepo.RemoteURL(project, gitProtocol)

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

func getTargetBranch(baseRepoRemote *glrepo.Remote) string {
	br, _ := git.GetDefaultBranch(baseRepoRemote.PushURL.String())
	// we ignore the error since git.GetDefaultBranch returns master and an error
	// if the default branch cannot be determined
	return br
}

package create

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/profclems/glab/pkg/prompt"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/internal/glrepo"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/run"
	"github.com/profclems/glab/pkg/git"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdCreate(f *cmdutils.Factory) *cobra.Command {
	var projectCreateCmd = &cobra.Command{
		Use:   "create [path] [flags]",
		Short: `Create a new Gitlab project/repository`,
		Long:  `Create a new GitHub repository.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateProject(cmd, args, f)
		},
		Example: heredoc.Doc(`
			# create a repository under your account using the current directory name
			$ glab repo create

			# create a repository under a group using the current directory name
			$ glab repo create --group glab-cli

			# create a repository with a specific name
			$ glab repo create my-project

			# create a repository for a group
			$ glab repo create glab-cli/my-project
	  `),
	}

	projectCreateCmd.Flags().StringP("name", "n", "", "Name of the new project")
	projectCreateCmd.Flags().StringP("group", "g", "", "Namespace/group for the new project (defaults to the current user’s namespace)")
	projectCreateCmd.Flags().StringP("description", "d", "", "Description of the new project")
	projectCreateCmd.Flags().String("defaultBranch", "", "Default branch of the project. If not provided, `master` by default.")
	projectCreateCmd.Flags().StringArrayP("tag", "t", []string{}, "The list of tags for the project.")
	projectCreateCmd.Flags().Bool("internal", false, "Make project internal: visible to any authenticated user (default)")
	projectCreateCmd.Flags().BoolP("private", "p", false, "Make project private: visible only to project members")
	projectCreateCmd.Flags().BoolP("public", "P", false, "Make project public: visible without any authentication")
	projectCreateCmd.Flags().Bool("readme", false, "Initialize project with README.md")

	return projectCreateCmd
}

func runCreateProject(cmd *cobra.Command, args []string, f *cmdutils.Factory) error {
	var (
		projectPath string
		visiblity   gitlab.VisibilityValue
		err         error
		isPath      bool
		namespaceID int
		namespace   string
	)
	c := f.IO.Color()

	defaultBranch, err := cmd.Flags().GetString("defaultBranch")
	if err != nil {
		return err
	}
	err = initGit(defaultBranch)
	if err != nil {
		return err
	}

	if len(args) == 1 {
		projectPath = args[0]
		if strings.Contains(projectPath, "/") {
			pp := strings.Split(projectPath, "/")
			projectPath = pp[1]
			namespace = pp[0]
		}
	} else {
		projectPath, err = git.ToplevelDir()
		projectPath = path.Base(projectPath)
		if err != nil {
			return err
		}
		isPath = true
	}

	apiClient, err := f.HttpClient()
	if err != nil {
		return err
	}

	group, err := cmd.Flags().GetString("group")
	if err != nil {
		return fmt.Errorf("could not parse group flag: %v", err)
	}
	if group != "" {
		namespace = group
	}

	if namespace != "" {
		group, err := api.GetGroup(apiClient, namespace)
		if err != nil {
			return fmt.Errorf("could not find group or namespace %s: %v", namespace, err)
		}
		namespaceID = group.ID
	}

	name, _ := cmd.Flags().GetString("name")

	if projectPath == "" && name == "" {
		fmt.Println("ERROR: Path or Name required to create project.")
		return cmd.Usage()
	} else if name == "" {
		name = projectPath
	}

	description, _ := cmd.Flags().GetString("description")

	if internal, _ := cmd.Flags().GetBool("internal"); internal {
		visiblity = gitlab.InternalVisibility
	} else if private, _ := cmd.Flags().GetBool("private"); private {
		visiblity = gitlab.PrivateVisibility
	} else if public, _ := cmd.Flags().GetBool("public"); public {
		visiblity = gitlab.PublicVisibility
	}

	tags, _ := cmd.Flags().GetStringArray("tag")
	readme, _ := cmd.Flags().GetBool("readme")

	opts := &gitlab.CreateProjectOptions{
		Name:                 gitlab.String(name),
		Path:                 gitlab.String(projectPath),
		Description:          gitlab.String(description),
		DefaultBranch:        gitlab.String(defaultBranch),
		TagList:              &tags,
		InitializeWithReadme: gitlab.Bool(readme),
	}

	if visiblity != "" {
		opts.Visibility = &visiblity
	}

	if namespaceID != 0 {
		opts.NamespaceID = &namespaceID
	}

	project, err := api.CreateProject(apiClient, opts)

	greenCheck := c.Green("✓")

	if err == nil {
		fmt.Fprintf(f.IO.StdOut, "%s Created repository %s on GitLab: %s\n", greenCheck, project.NameWithNamespace, project.WebURL)
		if isPath {
			cfg, _ := f.Config()
			webURL, _ := url.Parse(project.WebURL)
			protocol, _ := cfg.Get(webURL.Host, "git_protocol")

			remoteArgs := &glrepo.RemoteArgs{
				Protocol: protocol,
				Url:      webURL.Host,
			}

			if protocol != "ssh" {
				remoteArgs.Token, _ = cfg.Get(webURL.Host, "token")
				currentUser, err := api.CurrentUser(apiClient)
				if err == nil {
					remoteArgs.Username = currentUser.Username
				}
			}
			remote, err := glrepo.RemoteURL(project, remoteArgs)
			if err != nil {
				return err
			}
			_, err = git.AddRemote("origin", remote)
			if err != nil {
				return err
			}
			fmt.Fprintf(f.IO.StdOut, "%s Added remote %s\n", greenCheck, remote)

		} else if f.IO.IsaTTY {
			var doSetup bool
			err := prompt.Confirm(&doSetup, fmt.Sprintf("Create a local project directory for %s?", project.NameWithNamespace), true)
			if err != nil {
				return err
			}

			if doSetup {
				projectPath := project.Path
				err = initialiseRepo(projectPath, project.SSHURLToRepo)
				if err != nil {
					return err
				}
				fmt.Fprintf(f.IO.StdOut, "%s Initialized repository in './%s/'\n", greenCheck, projectPath)
			}
		}
	} else {
		return fmt.Errorf("error creating project: %v", err)
	}
	return err
}

func initGit(defaultBranch string) error {
	if stat, err := os.Stat(".git"); err == nil && stat.IsDir() {
		return nil
	}
	var doInit bool
	err := prompt.Confirm(&doInit, "Directory not git initialized. Run `git init`?", true)
	if err != nil || !doInit {
		return err
	}

	gitInit := git.GitCommand("init")
	gitInit.Stdout = os.Stdout
	gitInit.Stderr = os.Stderr
	err = run.PrepareCmd(gitInit).Run()
	if err != nil {
		return err
	}

	if defaultBranch == "" {
		return nil
	}

	gitBranch := git.GitCommand("checkout", "-b", defaultBranch)
	gitBranch.Stdout = os.Stdout
	gitBranch.Stdin = os.Stdin
	return run.PrepareCmd(gitBranch).Run()
}

func initialiseRepo(projectPath, remoteURL string) error {

	gitInit := git.GitCommand("init", projectPath)
	gitInit.Stdout = os.Stdout
	gitInit.Stderr = os.Stderr
	err := run.PrepareCmd(gitInit).Run()
	if err != nil {
		return err
	}
	gitRemoteAdd := git.GitCommand("-C", projectPath, "remote", "add", "origin", remoteURL)
	gitRemoteAdd.Stdout = os.Stdout
	gitRemoteAdd.Stderr = os.Stderr
	err = run.PrepareCmd(gitRemoteAdd).Run()
	if err != nil {
		return err
	}
	return nil
}

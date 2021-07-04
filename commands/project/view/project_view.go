package view

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/iostreams"
	"github.com/profclems/glab/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type ViewOptions struct {
	ProjectID    string
	APIClient    *gitlab.Client
	Web          bool
	Branch       string
	Browser      string
	GlamourStyle string

	IO   *iostreams.IOStreams
	Repo glrepo.Interface
}

func NewCmdView(f *cmdutils.Factory) *cobra.Command {
	opts := ViewOptions{
		IO: f.IO,
	}

	var projectViewCmd = &cobra.Command{
		Use:   "view [repository] [flags]",
		Short: "View a project/repository",
		Long:  `Display the description and README of a project or open it in the browser.`,
		Args:  cobra.MaximumNArgs(1),
		Example: heredoc.Doc(`
			# view project information for the current directory
			$ glab repo view

			# view project information of specified name
			$ glab repo view my-project
			$ glab repo view user/repo
			$ glab repo view group/namespace/repo

			# specify repo by full [git] URL
			$ glab repo view git@gitlab.com:user/repo.git
			$ glab repo view https://gitlab.company.org/user/repo
			$ glab repo view https://gitlab.company.org/user/repo.git
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			cfg, err := f.Config()
			if err != nil {
				return err
			}

			if len(args) == 1 {
				opts.ProjectID = args[0]
			}

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}
			opts.APIClient = apiClient

			if opts.ProjectID == "" {
				opts.Repo, err = f.BaseRepo()
				if err != nil {
					return cmdutils.WrapError(err, "`repository` is required when not running in a git repository")
				}
				opts.ProjectID = opts.Repo.FullName()
			}

			if opts.ProjectID != "" {
				if !strings.Contains(opts.ProjectID, "/") {
					currentUser, err := api.CurrentUser(opts.APIClient)
					if err != nil {
						return cmdutils.WrapError(err, "Failed to retrieve your current user")
					}

					opts.ProjectID = currentUser.Username + "/" + opts.ProjectID
				}

				repo, err := glrepo.FromFullName(opts.ProjectID)
				if err != nil {
					return err
				}

				if !glrepo.IsSame(repo, opts.Repo) {
					client, err := api.NewClientWithCfg(repo.RepoHost(), cfg, false)
					if err != nil {
						return err
					}
					opts.APIClient = client.Lab()
				}
				opts.Repo = repo
				opts.ProjectID = repo.FullName()
			}

			if opts.Branch == "" {
				opts.Branch, _ = f.Branch()
			}

			browser, _ := cfg.Get(opts.Repo.RepoHost(), "browser")
			opts.Browser = browser

			opts.GlamourStyle, _ = cfg.Get(opts.Repo.RepoHost(), "glamour_style")

			return runViewProject(&opts)
		},
	}

	projectViewCmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open a project in the browser")
	projectViewCmd.Flags().StringVarP(&opts.Branch, "branch", "b", "", "View a specific branch of the repository")

	return projectViewCmd
}

func runViewProject(opts *ViewOptions) error {
	project, err := opts.Repo.Project(opts.APIClient)
	if err != nil {
		return cmdutils.WrapError(err, "Failed to retrieve project information")
	}

	readmeFile, err := getReadmeFile(opts, project)
	if err != nil {
		return err
	}

	if opts.Web {
		openURL := generateProjectURL(project, opts.Branch)

		if opts.IO.IsaTTY {
			fmt.Fprintf(opts.IO.StdOut, "Opening %s in your browser.\n", utils.DisplayURL(openURL))
		}

		return utils.OpenInBrowser(openURL, opts.Browser)
	} else {
		if opts.IO.IsaTTY {
			if err := opts.IO.StartPager(); err != nil {
				return err
			}
			defer opts.IO.StopPager()

			printProjectContentTTY(opts, project, readmeFile)
		} else {
			printProjectContentRaw(opts, project, readmeFile)
		}
	}

	return nil
}

func getReadmeFile(opts *ViewOptions, project *gitlab.Project) (*gitlab.File, error) {
	if project.ReadmeURL == "" {
		return nil, nil
	}

	readmePath := strings.Replace(project.ReadmeURL, project.WebURL+"/-/blob/", "", 1)
	readmePathComponents := strings.Split(readmePath, "/")
	readmeRef := readmePathComponents[0]
	readmeFileName := readmePathComponents[1]
	readmeFile, err := api.GetFile(opts.APIClient, project.PathWithNamespace, readmeFileName, readmeRef)

	if err != nil {
		return nil, cmdutils.WrapError(err, "Failed to retrieve README file")
	}

	decoded, err := base64.StdEncoding.DecodeString(readmeFile.Content)
	if err != nil {
		return nil, cmdutils.WrapError(err, "Failed to decode README file")
	}

	readmeFile.Content = string(decoded)

	return readmeFile, nil
}

func generateProjectURL(project *gitlab.Project, branch string) string {
	if project.DefaultBranch != branch {
		return project.WebURL + "/-/tree/" + branch
	}

	return project.WebURL
}

func printProjectContentTTY(opts *ViewOptions, project *gitlab.Project, readme *gitlab.File) {
	var description string
	var readmeContent string
	var err error

	fullName := project.NameWithNamespace
	if project.Description != "" {
		description, err = utils.RenderMarkdownWithoutIndentations(project.Description, opts.GlamourStyle)

		if err != nil {
			description = project.Description
		}
	} else {
		description = "\n(No description provided)\n\n"
	}

	if readme != nil {
		readmeContent, err = utils.RenderMarkdown(readme.Content, opts.GlamourStyle)

		if err != nil {
			readmeContent = readme.Content
		}
	}

	c := opts.IO.Color()
	// Header
	fmt.Fprint(opts.IO.StdOut, c.Bold(fullName))
	fmt.Fprint(opts.IO.StdOut, c.Gray(description))

	// Readme
	if readme != nil {
		fmt.Fprint(opts.IO.StdOut, readmeContent)
	} else {
		fmt.Fprintln(opts.IO.StdOut, c.Gray("(This repository does not have a README file)"))
	}

	fmt.Fprintln(opts.IO.StdOut)
	fmt.Fprintf(opts.IO.StdOut, c.Gray("View this project on GitLab: %s\n"), project.WebURL)
}

func printProjectContentRaw(opts *ViewOptions, project *gitlab.Project, readme *gitlab.File) {
	fullName := project.NameWithNamespace
	description := project.Description

	fmt.Fprintf(opts.IO.StdOut, "name:\t%s\n", fullName)
	fmt.Fprintf(opts.IO.StdOut, "description:\t%s\n", description)

	if readme != nil {
		fmt.Fprintln(opts.IO.StdOut, "---")
		fmt.Fprintf(opts.IO.StdOut, readme.Content)
		fmt.Fprintln(opts.IO.StdOut)
	}
}

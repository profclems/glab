package run

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/pkg/api"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

const keyValuePair = ".+:.+"

var re = regexp.MustCompile(keyValuePair)

func getDefaultBranch(f *cmdutils.Factory) string {
	repo, err := f.BaseRepo()
	if err != nil {
		return "master"
	}

	remotes, err := f.Remotes()
	if err != nil {
		return "master"
	}

	repoRemote, err := remotes.FindByRepo(repo.RepoOwner(), repo.RepoName())
	if err != nil {
		return "master"
	}

	branch, _ := git.GetDefaultBranch(repoRemote.Name)

	return branch
}

func NewCmdRun(f *cmdutils.Factory) *cobra.Command {
	var pipelineRunCmd = &cobra.Command{
		Use:     "run [flags]",
		Short:   `Create or run a new CI pipeline`,
		Aliases: []string{"create"},
		Example: heredoc.Doc(`
	$ glab ci run
	$ glab ci run -b trunk
	`),
		Long: ``,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			pipelineVars := []*gitlab.PipelineVariable{}

			if customPipelineVars, _ := cmd.Flags().GetStringSlice("variables"); len(customPipelineVars) > 0 {
				for _, v := range customPipelineVars {
					if !re.MatchString(v) {
						return fmt.Errorf("Bad pipeline variable : \"%s\" should be of format KEY:VALUE", v)
					}
					s := strings.Split(v, ":")
					pipelineVars = append(pipelineVars, &gitlab.PipelineVariable{
						Key:          s[0],
						Value:        s[1],
						VariableType: "env_var",
					})
				}
			}

			c := &gitlab.CreatePipelineOptions{
				Variables: pipelineVars,
			}

			if m, _ := cmd.Flags().GetString("branch"); m != "" {
				c.Ref = gitlab.String(m)
			} else {
				c.Ref = gitlab.String(getDefaultBranch(f))
			}

			pipe, err := api.CreatePipeline(apiClient, repo.FullName(), c)
			if err != nil {
				return err
			}

			fmt.Fprintln(f.IO.StdOut, "Created pipeline (id:", pipe.ID, "), status:", pipe.Status, ", ref:", pipe.Ref, ", weburl: ", pipe.WebURL, ")")
			return nil
		},
	}
	pipelineRunCmd.Flags().StringP("branch", "b", "", "Create pipeline on branch/ref <string>")
	pipelineRunCmd.Flags().StringSliceP("variables", "", []string{}, "Pass variables to pipeline")

	return pipelineRunCmd
}

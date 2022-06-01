package trigger

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/git"

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
		Use:     "trigger [flags]",
		Short:   `Trigger a pipeline in the CI`,
		Aliases: []string{"t"},
		Example: heredoc.Doc(`
	$ glab ci trigger
	$ glab ci trigger -b trunk
  	$ glab ci run -b trunk --variables MYKEY:some_value --variables KEY2:another_value
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

			pipelineVars := map[string]string{
				// Do this, so the
				"CI_PIPELINE_SOURCE" : "trigger",
			}

			if customRunPipelineVars, _ := cmd.Flags().GetStringSlice("variables"); len(customRunPipelineVars) > 0 {
				for _, v := range customRunPipelineVars {
					if !re.MatchString(v) {
						return fmt.Errorf("Bad pipeline variable : \"%s\" should be of format KEY:VALUE", v)
					}
					s := strings.SplitN(v, ":", 2)
					pipelineVars[s[0]] = s[1]
				}
			}

			c := &gitlab.RunPipelineTriggerOptions{
				Variables: pipelineVars,
			}

			if m, _ := cmd.Flags().GetString("branch"); m != "" {
				c.Ref = gitlab.String(m)
			} else {
				c.Ref = gitlab.String(getDefaultBranch(f))
			}

			pipe, err := api.RunPipelineTrigger(apiClient, repo.FullName(), c)
			if err != nil {
				return err
			}

			fmt.Fprintln(f.IO.StdOut, "Ran pipeline (id:", pipe.ID, "), status:", pipe.Status, ", ref:", pipe.Ref, ", weburl: ", pipe.WebURL, ")")
			return nil
		},
	}
	pipelineRunCmd.Flags().StringP("branch", "b", "", "Run pipeline on branch/ref <string>")
	pipelineRunCmd.Flags().StringSliceP("variables", "", []string{}, "Pass variables to pipeline run")

	return pipelineRunCmd
}

package run

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/pkg/api"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

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
		Short:   `Create a new pipeline run`,
		Aliases: []string{"create"},
		Example: heredoc.Doc(`
	$ glab pipeline run
	$ glab pipeline run -b trunk
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

			// TODO: support setting pipeline variables via cli.
			v := []*gitlab.PipelineVariable{
				{
					Key:          "GLAB_CLI_KEY",
					Value:        "GLAB_CLI_VAL",
					VariableType: "env_var",
				},
			}

			c := &gitlab.CreatePipelineOptions{
				Variables: v,
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

	return pipelineRunCmd
}

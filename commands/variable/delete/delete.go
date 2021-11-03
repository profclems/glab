package delete

import (
	"fmt"

	"errors"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/variable/variableutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/iostreams"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type DeleteOpts struct {
	HTTPClient func() (*gitlab.Client, error)
	IO         *iostreams.IOStreams
	BaseRepo   func() (glrepo.Interface, error)

	Key   string
	Scope string
	Group string
}

func NewCmdSet(f *cmdutils.Factory, runE func(opts *DeleteOpts) error) *cobra.Command {
	opts := &DeleteOpts{
		IO: f.IO,
	}

	cmd := &cobra.Command{
		Use:     "delete <key>",
		Short:   "Delete a project or group variable",
		Aliases: []string{"remove"},
		Args:    cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ glab variable delete VAR_NAME
			$ glab variable delete VAR_NAME --scope=prod
			$ glab variable delete VARNAME -g mygroup
		`),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.HTTPClient = f.HttpClient
			opts.BaseRepo = f.BaseRepo
			opts.Key = args[0]

			if !variableutils.IsValidKey(opts.Key) {
				err = cmdutils.FlagError{Err: fmt.Errorf("invalid key provided.\n%s", variableutils.ValidKeyMsg)}
				return
			} else if len(args) != 1 {
				err = cmdutils.FlagError{Err: errors.New("no key provided")}
			}

			if cmd.Flags().Changed("scope") && opts.Group != "" {
				err = cmdutils.FlagError{Err: errors.New("scope is not required for group variables")}
				return
			}

			if runE != nil {
				err = runE(opts)
				return
			}
			err = deleteRun(opts)
			return
		},
	}

	cmd.Flags().StringVarP(&opts.Scope, "scope", "s", "*", "The environment_scope of the variable. All (*), or specific environments")
	cmd.Flags().StringVarP(&opts.Group, "group", "g", "", "Delete variable from a group")

	return cmd

}

func deleteRun(opts *DeleteOpts) error {
	c := opts.IO.Color()
	httpClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}

	baseRepo, err := opts.BaseRepo()
	if err != nil {
		return err
	}

	if opts.Group == "" {
		// Delete project-level variable
		err = api.DeleteProjectVariable(httpClient, baseRepo.FullName(), opts.Key, opts.Scope)
		if err != nil {
			return err
		}

		fmt.Fprintf(opts.IO.StdOut, "%s Deleted variable %s with scope %s for %s\n", c.GreenCheck(), opts.Key, opts.Scope, baseRepo.FullName())
	} else {
		// Delete group-level variable
		err = api.DeleteGroupVariable(httpClient, opts.Group, opts.Key)
		if err != nil {
			return err
		}

		fmt.Fprintf(opts.IO.StdOut, "%s Deleted variable %s for group %s\n", c.GreenCheck(), opts.Key, opts.Group)
	}

	return nil
}

package update

import (
	"errors"
	"fmt"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/variable/variableutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type UpdateOpts struct {
	HTTPClient func() (*gitlab.Client, error)
	IO         *iostreams.IOStreams
	BaseRepo   func() (glrepo.Interface, error)

	Key       string
	Value     string
	Type      string
	Scope     string
	Protected bool
	Masked    bool
	Group     string
}

func NewCmdSet(f *cmdutils.Factory, runE func(opts *UpdateOpts) error) *cobra.Command {
	opts := &UpdateOpts{
		IO: f.IO,
	}

	cmd := &cobra.Command{
		Use:   "update <key> <value>",
		Short: "Update an existing project or group variable",
		Args:  cobra.RangeArgs(1, 2),
		Example: heredoc.Doc(`
			$ glab variable update WITH_ARG "some value"
			$ glab variable update FROM_FLAG -v "some value"
			$ glab variable update FROM_ENV_WITH_ARG "${ENV_VAR}"
			$ glab variable update FROM_ENV_WITH_FLAG -v"${ENV_VAR}"
			$ glab variable update FROM_FILE < secret.txt
			$ cat file.txt | glab variable update SERVER_TOKEN
			$ cat token.txt | glab variable update GROUP_TOKEN -g mygroup --scope=prod
		`),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Supports repo override
			opts.HTTPClient = f.HttpClient
			opts.BaseRepo = f.BaseRepo

			opts.Key = args[0]

			if !variableutils.IsValidKey(opts.Key) {
				err = cmdutils.FlagError{Err: fmt.Errorf("invalid key provided.\n%s", variableutils.ValidKeyMsg)}
				return
			}

			if opts.Value != "" && len(args) == 2 {
				err = cmdutils.FlagError{Err: errors.New("specify value either by second positional argument or --value flag")}
				return
			}

			if cmd.Flags().Changed("scope") && opts.Group != "" {
				err = cmdutils.FlagError{Err: errors.New("scope is not required for group variables")}
				return
			}

			opts.Value, err = variableutils.GetValue(opts.Value, opts.IO, args)
			if err != nil {
				return
			}

			if cmd.Flags().Changed("type") {
				if opts.Type != "env_var" && opts.Type != "file" {
					err = cmdutils.FlagError{Err: fmt.Errorf("invalid type: %s. --type must be one of `env_var` or `file`", opts.Type)}
					return
				}
			}

			if runE != nil {
				err = runE(opts)
				return
			}
			err = updateRun(opts)
			return
		},
	}

	cmd.Flags().StringVarP(&opts.Value, "value", "v", "", "The value of a variable")
	cmd.Flags().StringVarP(&opts.Type, "type", "t", "env_var", "The type of a variable: {env_var|file}")
	cmd.Flags().StringVarP(&opts.Scope, "scope", "s", "*", "The environment_scope of the variable. All (*), or specific environments")
	cmd.Flags().StringVarP(&opts.Group, "group", "g", "", "Set variable for a group")
	cmd.Flags().BoolVarP(&opts.Masked, "masked", "m", false, "Whether the variable is masked")
	cmd.Flags().BoolVarP(&opts.Protected, "protected", "p", false, "Whether the variable is protected")
	return cmd
}

func updateRun(opts *UpdateOpts) error {
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
		// update project-level variable
		updateProjectVarOpts := &gitlab.UpdateProjectVariableOptions{
			Value:            gitlab.String(opts.Value),
			VariableType:     gitlab.VariableType(gitlab.VariableTypeValue(opts.Type)),
			Masked:           gitlab.Bool(opts.Masked),
			Protected:        gitlab.Bool(opts.Protected),
			EnvironmentScope: gitlab.String(opts.Scope),
		}

		_, err = api.UpdateProjectVariable(httpClient, baseRepo.FullName(), opts.Key, updateProjectVarOpts)
		if err != nil {
			return err
		}

		fmt.Fprintf(opts.IO.StdOut,
			"%s Updated variable %s for project %s with scope %s\n",
			c.GreenCheck(),
			opts.Key,
			baseRepo.FullName(),
			opts.Scope)

	} else {
		// update group-level variable
		updateGroupVarOpts := &gitlab.UpdateGroupVariableOptions{
			Value:            gitlab.String(opts.Value),
			VariableType:     gitlab.VariableType(gitlab.VariableTypeValue(opts.Type)),
			Masked:           gitlab.Bool(opts.Masked),
			Protected:        gitlab.Bool(opts.Protected),
			EnvironmentScope: gitlab.String(opts.Scope),
		}

		_, err = api.UpdateGroupVariable(httpClient, opts.Group, opts.Key, updateGroupVarOpts)
		if err != nil {
			return err
		}

		fmt.Fprintf(opts.IO.StdOut, "%s Updated variable %s for group %s\n", c.GreenCheck(), opts.Key, opts.Group)
	}

	return nil

}

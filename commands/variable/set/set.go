package set

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type SetOpts struct {
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

	ValueSet bool
}

func NewCmdSet(f *cmdutils.Factory, runE func(opts *SetOpts) error) *cobra.Command {
	opts := &SetOpts{
		IO: f.IO,
	}

	validKeyMsg := "A valid key must have no more than 255 characters; only A-Z, a-z, 0-9, and _ are allowed"

	cmd := &cobra.Command{
		Use:     "set <key> <value>",
		Short:   "Create a new project or group variable",
		Aliases: []string{"new", "create"},
		Args:    cobra.RangeArgs(1, 2),
		Example: heredoc.Doc(`
			$ glab variable set WITH_ARG "some value"
			$ glab variable set FROM_FLAG -v "some value"
			$ glab variable set FROM_ENV_WITH_ARG "${ENV_VAR}"
			$ glab variable set FROM_ENV_WITH_FLAG -v"${ENV_VAR}"
			$ glab variable set FROM_FILE < secret.txt
			$ cat file.txt | glab variable set SERVER_TOKEN
			$ cat token.txt | glab variable set GROUP_TOKEN -g mygroup --scope=prod
		`),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Supports repo override
			opts.HTTPClient = f.HttpClient
			opts.BaseRepo = f.BaseRepo

			opts.Key = args[0]

			if !isValidKey(opts.Key) {
				err = cmdutils.FlagError{Err: fmt.Errorf("invalid key provided.\n%s", validKeyMsg)}
				return
			}

			if opts.Value != "" && len(args) == 2 {
				if opts.Value != "" {
					err = cmdutils.FlagError{Err: errors.New("specify value either by second positional argument or --value flag")}
					return
				}
				opts.Value = args[1]
			}

			if cmd.Flags().Changed("scope") && opts.Group != "" {
				err = cmdutils.FlagError{Err: errors.New("scope is not required for group variables")}
				return
			}

			opts.ValueSet = cmd.Flags().Changed("value") || len(args) == 2
			opts.Value, err = getValue(opts)
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
			err = setRun(opts)
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

func setRun(opts *SetOpts) error {
	httpClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}
	if opts.Group != "" {
		// creating project-level variable
		createVarOpts := &gitlab.CreateGroupVariableOptions{
			Key:          gitlab.String(opts.Key),
			Value:        gitlab.String(opts.Value),
			VariableType: gitlab.VariableType(gitlab.VariableTypeValue(opts.Type)),
			Masked:       gitlab.Bool(opts.Masked),
			Protected:    gitlab.Bool(opts.Protected),
		}
		_, err = api.CreateGroupVariable(httpClient, opts.Group, createVarOpts)
		if err != nil {
			return err
		}

		fmt.Fprintf(opts.IO.StdOut, "%s Created variable %s for group %s\n", utils.GreenCheck(), opts.Key, opts.Group)
		return nil
	}

	// creating group-level variable
	baseRepo, err := opts.BaseRepo()
	if err != nil {
		return err
	}
	createVarOpts := &gitlab.CreateProjectVariableOptions{
		Key:              gitlab.String(opts.Key),
		Value:            gitlab.String(opts.Value),
		EnvironmentScope: gitlab.String(opts.Scope),
		Masked:           gitlab.Bool(opts.Masked),
		Protected:        gitlab.Bool(opts.Protected),
		VariableType:     gitlab.VariableType(gitlab.VariableTypeValue(opts.Type)),
	}
	_, err = api.CreateProjectVariable(httpClient, baseRepo.FullName(), createVarOpts)
	if err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.StdOut, "%s Created variable %s for %s\n", utils.GreenCheck(), opts.Key, baseRepo.FullName())
	return nil
}

func getValue(opts *SetOpts) (string, error) {
	if opts.Value != "" || opts.ValueSet {
		return opts.Value, nil
	}

	if opts.IO.IsInTTY {
		return "", &cmdutils.FlagError{Err: errors.New("no value specified but nothing on STDIN")}
	}

	// read value from STDIN if not provided
	defer opts.IO.In.Close()
	value, err := ioutil.ReadAll(opts.IO.In)
	if err != nil {
		return "", fmt.Errorf("failed to read value from STDIN: %w", err)
	}
	return strings.TrimSpace(string(value)), nil
}

// isValidKey checks if a key is valid if it follows the following criteria:
// must have no more than 255 characters;
// only A-Z, a-z, 0-9, and _ are allowed
func isValidKey(key string) bool {
	// check if key falls within range of 1-255
	if len(key) > 255 || len(key) < 1 {
		return false
	}
	keyRE := regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	return keyRE.MatchString(key)
}

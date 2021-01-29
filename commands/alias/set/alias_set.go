package set

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/commands/cmdutils"

	"github.com/MakeNowJust/heredoc"
	"github.com/google/shlex"
	"github.com/profclems/glab/internal/config"
	"github.com/spf13/cobra"
)

type SetOptions struct {
	Config    func() (config.Config, error)
	Name      string
	Expansion string
	IsShell   bool
	RootCmd   *cobra.Command
	IO        *iostreams.IOStreams
}

func NewCmdSet(f *cmdutils.Factory, runF func(*SetOptions) error) *cobra.Command {
	opts := &SetOptions{
		Config: f.Config,
	}

	var aliasSetCmd = &cobra.Command{
		Use:   "set <alias name> '<command>' [flags]",
		Short: `Set an alias.`,
		Long: heredoc.Doc(`
			Declare a word as a command alias that will expand to the specified command(s).

			The expansion may specify additional arguments and flags. If the expansion
			includes positional placeholders such as '$1', '$2', etc., any extra arguments
			that follow the invocation of an alias will be inserted appropriately.

			If '--shell' is specified, the alias will be run through a shell interpreter (sh). This allows you
			to compose commands with "|" or redirect with ">". Note that extra arguments following the alias
			will not be automatically passed to the expanded expression. To have a shell alias receive
			arguments, you must explicitly accept them using "$1", "$2", etc., or "$@" to accept all of them.

			Platform note: on Windows, shell aliases are executed via "sh" as installed by Git For Windows. If
			you have installed git on Windows in some other way, shell aliases may not work for you.
			Quotes must always be used when defining a command as in the examples.
		`),
		Example: heredoc.Doc(`
		$ glab alias set mrv 'mr view'
		$ glab mrv -w 123
		#=> glab mr view -w 123

		$ glab alias set createissue 'glab create issue --title "$1"'
		$ glab createissue "My Issue" --description "Something is broken."
		# => glab create issue --title "My Issue" --description "Something is broken."

		$ glab alias set --shell igrep 'glab issue list --assignee="$1" | grep $2'
		$ glab igrep user foo
		#=> glab issue list --assignee="user" | grep "foo"
	`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.RootCmd = cmd.Root()
			opts.Name = args[0]
			opts.Expansion = args[1]
			opts.IO = f.IO

			if runF != nil {
				return runF(opts)
			}
			return setRun(cmd, opts)
		},
	}
	aliasSetCmd.Flags().BoolVarP(&opts.IsShell, "shell", "s", false, "Declare an alias to be passed through a shell interpreter")
	return aliasSetCmd
}

func setRun(cmd *cobra.Command, opts *SetOptions) error {
	c := opts.IO.Color()
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	aliasCfg, err := cfg.Aliases()
	if err != nil {
		return err
	}

	if opts.IO.IsaTTY && opts.IO.IsErrTTY {
		fmt.Fprintf(opts.IO.StdErr, "- Adding alias for %s: %s\n", c.Bold(opts.Name), c.Bold(opts.Expansion))
	}

	expansion := opts.Expansion
	isShell := opts.IsShell
	if isShell && !strings.HasPrefix(expansion, "!") {
		expansion = "!" + expansion
	}
	isShell = strings.HasPrefix(expansion, "!")

	if validCommand(opts.RootCmd, opts.Name) {
		return fmt.Errorf("could not create alias: %q is already a glab command", opts.Name)
	}

	if !isShell && !validCommand(opts.RootCmd, expansion) {
		return fmt.Errorf("could not create alias: %s does not correspond to a glab command", expansion)
	}

	successMsg := fmt.Sprintf("%s Added alias.", c.Green("✓"))
	if oldExpansion, ok := aliasCfg.Get(opts.Name); ok {
		successMsg = fmt.Sprintf("%s Changed alias %s from %s to %s",
			c.Green("✓"),
			c.Bold(opts.Name),
			c.Bold(oldExpansion),
			c.Bold(expansion),
		)
	}

	err = aliasCfg.Set(opts.Name, expansion)
	if err != nil {
		return fmt.Errorf("could not create alias: %s", err)
	}

	fmt.Fprintln(opts.IO.StdErr, successMsg)
	return nil
}

func validCommand(rootCmd *cobra.Command, expansion string) bool {
	split, err := shlex.Split(expansion)
	if err != nil {
		return false
	}

	cmd, _, err := rootCmd.Traverse(split)
	return err == nil && cmd != rootCmd
}

package commands

import (
	"fmt"
	"github.com/google/shlex"
	"strings"

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/utils"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

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
		aliasName := args[0]
		aliasedCommand := args[1]
		// Check if provided alias name is already a glab command
		// err should be <nil> if alias name already exists as a command
		fmt.Fprintf(colorableOut(cmd), "- Adding alias for %s: %s\n", utils.Bold(aliasName), utils.Bold(aliasedCommand))

		isShell, err := cmd.Flags().GetBool("shell")
		if err != nil {
			return err
		}
		if isShell && !strings.HasPrefix(aliasedCommand, "!") {
			aliasedCommand = "!" + aliasedCommand
		}
		isShell = strings.HasPrefix(aliasedCommand, "!")

		if validCommand(RootCmd, aliasName) {
			return fmt.Errorf("could not create alias: \"%s\" is already a glab command", aliasName)
		}

		if !isShell && !validCommand(RootCmd, aliasedCommand) {
			return fmt.Errorf("could not create alias: %s does not correspond to a glab command", aliasedCommand)
		}

		successMsg := fmt.Sprintf("%s Added alias.", utils.GreenCheck())
		if oldExpansion := config.GetAlias(aliasName); oldExpansion != "" {
			successMsg = fmt.Sprintf("%s Changed alias %s from %s to %s",
				utils.Green("✓"),
				utils.Bold(aliasName),
				utils.Bold(oldExpansion),
				utils.Bold(aliasedCommand),
			)
		}
		err = config.SetAlias(aliasName, aliasedCommand)
		if err != nil {
			return fmt.Errorf("could not create alias: %s", err)
		}
		fmt.Fprintln(colorableOut(cmd), successMsg)
		return nil
	},
}

func validCommand(rootCmd *cobra.Command, expansion string) bool {
	split, err := shlex.Split(expansion)
	if err != nil {
		return false
	}

	cmd, _, err := rootCmd.Traverse(split)
	return err == nil && cmd != rootCmd
}

func init() {
	aliasSetCmd.Flags().BoolP("shell", "s", false, "Declare an alias to be passed through a shell interpreter")
	aliasCmd.AddCommand(aliasSetCmd)
}

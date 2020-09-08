package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/profclems/glab/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(completionCmd)
	completionCmd.Flags().StringP("shell", "s", "", "Shell type: {bash|zsh|fish|powershell}")
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for glab commands.
The output of this command will be computer code and is meant to be saved to a
file or immediately evaluated by an interactive shell.
For example, for bash you could add this to your '~/.bash_profile':
	eval "$(glab completion -s bash)"
When installing glab through a package manager, however, it's possible that
no additional shell configuration is necessary to gain completion support. For
Homebrew, see <https://docs.brew.sh/Shell-Completion>
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		shellType, err := cmd.Flags().GetString("shell")
		if err != nil {
			return err
		}

		if shellType == "" {
			out := cmd.OutOrStdout()
			isTTY := false
			if outFile, isFile := out.(*os.File); isFile {
				isTTY = utils.IsTerminal(outFile)
			}

			if isTTY {
				return errors.New("error: the value for `--shell` is required\nsee `glab help completion` for more information")
			}
			shellType = "bash"
		}

		switch shellType {
		case "bash":
			return RootCmd.GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			return RootCmd.GenZshCompletion(cmd.OutOrStdout())
		case "powershell":
			return RootCmd.GenPowerShellCompletion(cmd.OutOrStdout())
		case "fish":
			return RootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
		default:
			return fmt.Errorf("unsupported shell type %q", shellType)
		}
	},
}

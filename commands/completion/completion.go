package completion

import (
	"errors"
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"

	"github.com/profclems/glab/internal/utils"
	"github.com/spf13/cobra"
)

func NewCmdCompletion(io *utils.IOStreams) *cobra.Command {
	var shellType string

	var completionCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for glab commands.

The output of this command will be computer code and is meant to be saved to a
file or immediately evaluated by an interactive shell.

For example, for bash you could add this to your '~/.bash_profile':

	eval "$(glab completion -s bash)"

When installing glab through a package manager, however, it's possible that
no additional shell configuration is necessary to gain completion support. 
For Homebrew, see <https://docs.brew.sh/Shell-Completion>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if shellType == "" {
				if io.IsaTTY && io.IsErrTTY {
					return &cmdutils.FlagError{Err: errors.New("error: the value for `--shell` is required")}
				}
				shellType = "bash"
			}

			out := io.StdOut
			rootCmd := cmd.Parent()

			switch shellType {
			case "bash":
				return rootCmd.GenBashCompletion(out)
			case "zsh":
				return rootCmd.GenZshCompletion(out)
			case "powershell":
				return rootCmd.GenPowerShellCompletion(out)
			case "fish":
				return rootCmd.GenFishCompletion(out, true)
			default:
				return fmt.Errorf("unsupported shell type %q", shellType)
			}
		},
	}

	completionCmd.Flags().StringVarP(&shellType, "shell", "s", "", "Shell type: {bash|zsh|fish|powershell}")
	return completionCmd
}

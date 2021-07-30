package completion

import (
	"fmt"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/spf13/cobra"
)

func NewCmdCompletion(io *iostreams.IOStreams) *cobra.Command {
	var (
		shellType string

		// description will not be added if true
		excludeDesc = false
	)

	var completionCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for glab commands.

The output of this command will be computer code and is meant to be saved to a
file or immediately evaluated by an interactive shell.

For example, for bash you could add this to your '~/.bash_profile':

	eval "$(glab completion -s bash)"

Generate a %[1]s_gh%[1]s completion script and put it somewhere in your %[1]s$fpath%[1]s:
				gh completion -s zsh > /usr/local/share/zsh/site-functions/_gh
			Ensure that the following is present in your %[1]s~/.zshrc%[1]s:
				autoload -U compinit
				compinit -i
			
			Zsh version 5.7 or later is recommended.

When installing glab through a package manager, however, it's possible that
no additional shell configuration is necessary to gain completion support. 
For Homebrew, see <https://docs.brew.sh/Shell-Completion>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := io.StdOut
			rootCmd := cmd.Parent()

			switch shellType {
			case "bash":
				return rootCmd.GenBashCompletionV2(out, !excludeDesc)
			case "zsh":
				return rootCmd.GenZshCompletion(out)
			case "powershell":
				return rootCmd.GenPowerShellCompletion(out)
			case "fish":
				return rootCmd.GenFishCompletion(out, !excludeDesc)
			default:
				return fmt.Errorf("unsupported shell type %q", shellType)
			}
		},
	}

	completionCmd.Flags().StringVarP(&shellType, "shell", "s", "bash", "Shell type: {bash|zsh|fish|powershell}")
	completionCmd.Flags().BoolVarP(&excludeDesc, "no-desc", "", false, "Do not include shell completion description. Only for bash and fish")
	return completionCmd
}

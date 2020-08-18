package commands

import (
	"glab/internal/config"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var aliasSetCmd = &cobra.Command{
	Use:   "set <alias name> '<command>' [flags]",
	Short: `Set an alias.`,
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		aliasName := args[0]
		aliasedCommand := args[1]
		config.SetAlias(aliasName, aliasedCommand)
	},
	Example: heredoc.Doc(`
	$ glab alias set createissue 'glab create issue --title "$1"'
	$ glab createissue "My Issue" --description "Something is broken."
	# => glab create issue --title "My Issue" --description "Something is broken."
	`),
}

func init() {
	aliasCmd.AddCommand(aliasSetCmd)
}

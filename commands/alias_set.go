package commands

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var aliasSetCmd = &cobra.Command{
	Use:   "set <alias name> '<command>' [flags]",
	Short: `Set an alias.`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
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

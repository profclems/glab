package commands

import (
	"fmt"
	"github.com/profclems/glab/internal/utils"

	"github.com/profclems/glab/internal/config"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var aliasSetCmd = &cobra.Command{
	Use:   "set <alias name> '<command>' [flags]",
	Short: `Set an alias.`,
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		aliasName := args[0]
		aliasedCommand := args[1]
		// Check if provided alias name is already a glab command
		// err should be <nil> if alias name already exists as a command
		fmt.Printf("- Adding alias for %s: %s\n", aliasName, aliasedCommand)
		_, _, err := RootCmd.Find(append([]string{""}, aliasName))
		if err == nil {
			return fmt.Errorf("could not create alias: \"%s\" is already a glab command", aliasName)
		}
		err = config.SetAlias(aliasName, aliasedCommand)
		if err != nil {
			return err
		}
		fmt.Println(utils.GreenCheck(), "Alias added")
		return nil
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

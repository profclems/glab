package commands

import "github.com/spf13/cobra"

// mrCmd is merge request command
var labelCmd = &cobra.Command{
	Use:   "label <command> [flags]",
	Short: `Manage labels on remote`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 2 {
			cmd.Help()
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(labelCmd)
}
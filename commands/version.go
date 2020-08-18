package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"glab/internal/git"
)

// Version is the version for glab
var Version string

// Build holds the date bin was released
var Build string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "show glab version information",
	Long:    ``,
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("glab version %s (%s)\n", Version, Build)
		if err := git.RunCmd([]string{"version"}); err != nil {
			fmt.Println(err)
		}
		fmt.Println("Made with ❤ by Clement Sam <clementsam75@gmail.com> and contributors")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

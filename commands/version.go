package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"glab/internal/git"
)

// Version is set at build
var Version string
var build string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "show glab version information",
	Long:    ``,
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("glab version %s (%s)\n", Version, build)
		if err := git.RunCmd([]string{"version"}); err != nil  {
			fmt.Println(err)
		}
		fmt.Println("Made with ❤ by Clement Sam <clementsam75@gmail.com")
	},
}

func init()  {
	RootCmd.AddCommand(versionCmd)
}
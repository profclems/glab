package commands

import (
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"glab/internal/git"
	"glab/internal/manip"
)

var mrDeleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Short:   `Delete merge requests`,
	Long:    ``,
	Aliases: []string{"del"},
	Args:    cobra.ExactArgs(1),
	Example: "$ glab delete 123",
	RunE:    deleteMergeRequest,
}

func deleteMergeRequest(cmd *cobra.Command, args []string) error {

	if len(args) > 0 {
		mergeID := strings.Trim(args[0], " ")
		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo = r
		}
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Deleting Merge Request #" + i2)
			issue, err := gitlabClient.MergeRequests.DeleteMergeRequest(repo, manip.StringToInt(i2))

			if issue != nil {
				if issue.StatusCode == 204 {
					fmt.Println(aurora.Green("Merge Request Deleted Successfully"))
				} else if issue.StatusCode == 404 {
					fmt.Println(aurora.Red("Merge Request does not exist"))
				} else if issue.StatusCode == 401 {
					fmt.Println(aurora.Red("You are not authorized to perform this action"))
				} else {
					fmt.Println(aurora.Red("Could not complete request."))
				}
			} else if err != nil {
				return err
			}
		}
	} else {
		cmd.Help()
	}
	return nil
}

func init() {
	mrCmd.AddCommand(mrDeleteCmd)
}

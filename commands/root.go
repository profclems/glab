package commands

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"glab/internal/config"
	"glab/internal/git"
	"glab/internal/update"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

// RootCmd is the main root/parent command
var RootCmd = &cobra.Command{
	Use:           "glab <command> <subcommand> [flags]",
	Short:         "A GitLab CLI Tool",
	Long:          `GLab is an open source Gitlab Cli tool bringing GitLab to your command line`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Example: heredoc.Doc(`
	$ glab issue create
	$ glab mr list --merged
	$ glab pipeline list
	`),
	Annotations: map[string]string{
		"help:environment": heredoc.Doc(`
			GITLAB_TOKEN: an authentication token for API requests. Setting this avoids being
			prompted to authenticate and overrides any previously stored credentials.
			Can be set with glab config --token=<YOUR-GITLAB-ACCESS-TOKEN>

			GITLAB_REPO: specify the Gitlab repository in "OWNER/REPO" format for commands that
			otherwise operate on a local repository. (Depreciated in v1.6.2) 
			Can be set with glab config --repo=OWNER/REPO

			GITLAB_URI: specify the url of the gitlab server if self hosted (eg: https://gitlab.example.com). Default is https://gitlab.com.
			Can be set with glab config --url=gitlab.example.com

			GIT_REMOTE_URL_VAR: git remote variable that contains the gitlab url. Defaults is origin
			Can be set with glab config --remote-var=VARIABLE
		`),
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Printf("Unknown command: %s\n", args[0])
			_ = cmd.Usage()
			return
		} else if ok, _ := cmd.Flags().GetBool("version"); ok {
			versionCmd.Run(cmd, args)
			return
		}

		_ = cmd.Help()
	},
}

// Execute executes the root command.
func Execute() (*cobra.Command, error) {
	return RootCmd.ExecuteC()
}

// versionCmd represents the version command
var updateCmd = &cobra.Command{
	Use:     "check-update",
	Short:   "Check for latest glab releases",
	Long:    ``,
	Aliases: []string{"update", ""},
	Run:     checkForUpdate,
}

var configCmd = &cobra.Command{
	Use:     "config [flags]",
	Short:   `Configuration`,
	Long:    ``,
	Aliases: []string{"conf"},
	Args:    cobra.MaximumNArgs(2),
	Run:     config.Set,
}

func init() {
	RootCmd.Flags().BoolP("version", "v", false, "show glab version information")
	RootCmd.AddCommand(updateCmd)
	initConfigCmd()
	RootCmd.AddCommand(configCmd)
}

func er(msg interface{}) {
	color.Error.Println("Error:", msg)
	os.Exit(1)
}
func cmdErr(cmd *cobra.Command, args []string) {
	color.Error.Println("Error: Unknown command:")
	_ = cmd.Usage()
}

func initConfigCmd() {
	configCmd.Flags().BoolP("global", "g", false, "Set configuration globally")
	configCmd.Flags().StringP("url", "u", "", "specify the url of the gitlab server if self hosted (eg: https://gitlab.example.com).")
	configCmd.Flags().StringP("remote-var", "o", "", "Shorthand name for the remote repository. An example of a remote shorthand name is `origin`")
	configCmd.Flags().StringP("token", "t", "", "an authentication token for API requests.")
}

func isSuccessful(code int) bool {
	if code >= 200 && code < 300 {
		return true
	}
	return false
}

func checkForUpdate(*cobra.Command, []string) {

	releaseInfo, err := update.CheckForUpdate()
	if err != nil {
		er("Could not check for update! Make sure you have a stable internet connection")
		return
	}
	latestVersion := strings.TrimSpace(releaseInfo.Name)
	Version = strings.TrimSpace(Version)
	if isLatestVersion(latestVersion, Version) {
		color.Green.Println("You are already using the latest version of glab")
	} else {
		color.Printf("<yellow>A new version of glab has been released:</> <red>%s</> → <green>%s</>\n%s\n",
			Version, latestVersion, releaseInfo.HTMLUrl)
	}
}

func isLatestVersion(latestVersion, appVersion string)  bool {
	latestVersion = strings.TrimSpace(latestVersion)
	appVersion = strings.TrimSpace(appVersion)
	vo, v1e := version.NewVersion(appVersion)
	vn, v2e := version.NewVersion(latestVersion)
	return v1e == nil && v2e == nil && vo.LessThan(vn)
}

// ListInfo represents the parameters required to display a list result.
type ListInfo struct {
	// Name of the List to be used in constructing Description and EmptyMessage if not provided.
	Name string
	// List of columns to display
	Columns []string
	// Total number of record. Ideally size of the List.
	Total int
	// Function to pick a cell value from cell index
	GetCellValue func(int, int) interface{}
	// Optional. Description of the List. If not provided, default one constructed from list Name.
	Description string
	// Optional. EmptyMessage to display when List is empty. If not provided, default one constructed from list Name.
	EmptyMessage string
	// TableWrap wraps the contents when the column length exceeds the maximum width
	TableWrap bool
}

// Prints the list data on console
func DisplayList(lInfo ListInfo) {
	table := uitable.New()
	table.MaxColWidth = 70
	table.Wrap = lInfo.TableWrap
	fmt.Println()

	if lInfo.Total > 0 {
		description := lInfo.Description
		if description == "" {
			description = fmt.Sprintf("Showing %s %d of %d on %s\n\n", lInfo.Name, lInfo.Total, lInfo.Total, git.GetRepo())
		}
		fmt.Println(description)
		header := make([]interface{}, len(lInfo.Columns))
		for ci, c := range lInfo.Columns {
			header[ci] = c
		}
		table.AddRow(header...)

		for ri := 0; ri < lInfo.Total; ri++ {
			row := make([]interface{}, len(lInfo.Columns))
			for ci := range lInfo.Columns {
				row[ci] = lInfo.GetCellValue(ri, ci)
			}
			table.AddRow(row...)
		}

		fmt.Println(table)
	} else {
		emptyMessage := lInfo.EmptyMessage
		if emptyMessage == "" {
			emptyMessage = fmt.Sprintf("No %s available on %s", lInfo.Name, git.GetRepo())
		}
		fmt.Println(emptyMessage)
	}

}

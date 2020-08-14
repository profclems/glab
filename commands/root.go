package commands

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"glab/internal/config"
	"glab/internal/update"
	"os"
	"strings"
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
func Execute() error {
	RootCmd.Flags().BoolP("version", "v", false, "show glab version information")
	return RootCmd.Execute()
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
	cobra.OnInitialize(initConfig)
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

func initConfig() {
	config.SetGlobalPathDir()
	config.UseGlobalConfig = true
	if config.GetEnv("GITLAB_URI") == "NOTFOUND" || config.GetEnv("GITLAB_URI") == "OK" {
		config.SetEnv("GITLAB_URI", "https://gitlab.com")
	}
	if config.GetEnv("GIT_REMOTE_URL_VAR") == "NOTFOUND" || config.GetEnv("GIT_REMOTE_URL_VAR") == "OK" {
		config.SetEnv("GIT_REMOTE_URL_VAR", "origin")
	}
	config.UseGlobalConfig = false
}

func initConfigCmd() {
	configCmd.Flags().BoolP("global", "g", false, "Set configuration globally")
	configCmd.Flags().StringP("url", "u", "", "specify the url of the gitlab server if self hosted (eg: https://gitlab.example.com).")
	configCmd.Flags().StringP("remote-var", "o", "", "delete merge request <id>")
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
	var latestVersion string
	if err != nil {
		er("Could not check for update! Make sure you have a stable internet connection")
		return
	}
	latestVersion = strings.TrimSpace(releaseInfo.Name)
	Version = strings.TrimSpace(Version)
	if latestVersion == Version {
		color.Green.Println("You are already using the latest version of glab")
	} else {
		color.Printf("<yellow>A new version of glab has been released:</> <red>%s</> → <green>%s</>\n", Version, latestVersion)
		fmt.Println(releaseInfo.HTMLUrl)
	}
}

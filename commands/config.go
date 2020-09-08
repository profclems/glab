package commands

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/utils"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [flags]",
	Short: `Set and get glab settings`,
	Long: heredoc.Doc(`Get and set key/value strings.
		Current respected settings:
		- gitlab_token: Your gitlab access token, defaults to environment variables
		- gitlab_uri: if unset, defaults to https://gitlab.com
		- git_remote_url_var, if unset, defaults to origin
		- browser: if unset, defaults to environment variables
		- editor: if unset, defaults to environment variables.
		- visual: alternative for editor. if unset, defaults to environment variables.
		- glamour_style: Your desired markdown renderer style. Options are dark, light, notty. Custom styles are allowed set a custom style
https://github.com/charmbracelet/glamour#styles
	`),
	Aliases: []string{"conf"},
	Args:    cobra.MaximumNArgs(2),
	RunE:    configFunc,
}

func init() {
	configCmd.Flags().BoolP("global", "g", false, "Set configuration globally")
	configCmd.Flags().StringP("url", "u", "", "specify the url of the gitlab server if self hosted (eg: https://gitlab.example.com).")
	configCmd.Flags().StringP("remote-var", "o", "", "Shorthand name for the remote repository. An example of a remote shorthand name is `origin`")
	configCmd.Flags().StringP("token", "t", "", "an authentication token for API requests.")

	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configInitCmd)

	configGetCmd.Flags().StringP("host", "", "", "Get per-host setting")
	configSetCmd.Flags().StringP("host", "", "", "Set per-host setting")
	configSetCmd.Flags().BoolP("global", "g", false, "write to global ~/.config/glab-cli/.env file rather than the repository .glab-cli/config/.env")
	configGetCmd.Flags().BoolP("global", "g", false, "Read from global config file (~/.config/glab-cli/.env). [Default: looks through OS → Local → Global]")

	// TODO reveal and add usage once we properly support multiple hosts
	_ = configGetCmd.Flags().MarkHidden("host")
	// TODO reveal and add usage once we properly support multiple hosts
	_ = configSetCmd.Flags().MarkHidden("host")
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Prints the value of a given configuration key",
	Long: `Get the value for a given configuration key.
Examples:
  $ glab config get gitlab_uri
  https://gitlab.com
  $ glab config get git_remote_url_var
  origin
`,
	Args: cobra.ExactArgs(1),
	RunE: configGet,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Updates configuration with the value of a given key",
	Long: `Update the configuration by setting a key to a value.
Examples:
  $ glab config set editor vim
  $ glab config set git_remote_url_var origin
`,
	Args: cobra.ExactArgs(2),
	RunE: configSet,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Shows a prompt to set basic glab configuration",
	Long: `Update the configuration by setting a key to a value.
Examples:
  $ glab config init
  ? Enter default Gitlab Host (Current Value: https://gitlab.com): |
`,
	Args: cobra.ExactArgs(2),
	RunE: configSet,
}

func configGet(cmd *cobra.Command, args []string) error {
	key := strings.ToUpper(args[0])
	/*
		hostname, err := cmd.Flags().GetString("host")
		if err != nil {
			return err
		}
	*/

	global, err := cmd.Flags().GetBool("global")
	if err != nil {
		return err
	}
	if global {
		config.UseGlobalConfig = true
	}
	val := config.GetEnv(key)

	if val != "" {
		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "%s\n", val)
	}

	return nil
}

func configSet(cmd *cobra.Command, args []string) error {
	key := strings.ToUpper(args[0])
	value := args[1]
	/*
		hostname, err := cmd.Flags().GetString("host")
		if err != nil {
			return err
		}
	*/

	global, err := cmd.Flags().GetBool("global")
	if err != nil {
		return err
	}
	if global {
		config.UseGlobalConfig = true
	}

	err = config.SetEnv(key, value)
	if err != nil {
		return fmt.Errorf("failed to set %q to %q: %w", key, value, err)
	}

	return nil
}

func configInit(cmd *cobra.Command, args []string) error {
	_, err := config.PromptAndSetEnv(fmt.Sprintf("Enter default Gitlab Host (Current Value: %s): ",
		config.GetEnv("GITLAB_URI")), "GITLAB_URI")
	if err != nil {
		return err
	}
	_, err = config.PromptAndSetEnv("Enter default Gitlab Token: ", "GITLAB_TOKEN")
	if err != nil {
		return err
	}
	_, err = config.PromptAndSetEnv(fmt.Sprintf("Enter Git remote url variable (Current Value: %s): ",
		config.GetEnv("GIT_REMOTE_URL_VAR")), "GIT_REMOTE_URL_VAR")
	if err != nil {
		return err
	}
	fmt.Fprintf(colorableOut(cmd), "%s Configuration updated", utils.GreenCheck())
	return nil
}

// TODO --flag=value config format is set to be deprecated in v2.0.0
func configFunc(cmd *cobra.Command, args []string) error {
	var isUpdated bool
	if b, _ := cmd.Flags().GetBool("global"); b {
		config.UseGlobalConfig = true
	}
	if b, _ := cmd.Flags().GetString("token"); b != "" {
		config.SetEnv("GITLAB_TOKEN", b)
		isUpdated = true
	}
	if b, _ := cmd.Flags().GetString("url"); b != "" {
		config.SetEnv("GITLAB_URI", b)
		isUpdated = true
	}
	if b, _ := cmd.Flags().GetString("remote-var"); b != "" {
		config.SetEnv("GIT_REMOTE_URL_VAR", b)
		isUpdated = true
	}
	if b, _ := cmd.Flags().GetString("pid"); b != "" {
		config.SetEnv("GITLAB_PROJECT_ID", b)
		isUpdated = true
	}
	if !isUpdated {
		err := configInit(cmd, args)
		if err != nil {
			return err
		}
	}
	if isUpdated {
		// Add depreciation warning
		fmt.Fprintf(colorableOut(cmd), "%s flag=value config format is set to be deprecated in later releases. Use 'glab config set <key> <value>'\n", utils.Yellow("!WARNING:"))
		fmt.Fprintf(colorableOut(cmd), "%s Configuration updated\n", utils.GreenCheck())
	}
	return nil
}

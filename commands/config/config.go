package config

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glinstance"
	"github.com/profclems/glab/internal/utils"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var isGlobal bool

func NewCmdConfig(f *cmdutils.Factory) *cobra.Command {
	var configCmd = &cobra.Command{
		Use:   "config [flags]",
		Short: `Set and get glab settings`,
		Long: heredoc.Doc(`Get and set key/value strings.

		Current respected settings:

		- token: Your gitlab access token, defaults to environment variables
		- gitlab_uri: if unset, defaults to https://gitlab.com
		- browser: if unset, defaults to environment variables
		- editor: if unset, defaults to environment variables.
		- visual: alternative for editor. if unset, defaults to environment variables.
		- glamour_style: Your desired markdown renderer style. Options are dark, light, notty. Custom styles are allowed set a custom style
https://github.com/charmbracelet/glamour#styles
	`),
		Aliases: []string{"conf"},
	}


	configCmd.Flags().BoolVarP(&isGlobal, "global", "g", false, "use global config file")

	configCmd.AddCommand(NewCmdConfigGet(f))
	configCmd.AddCommand(NewCmdConfigSet(f))
	configCmd.AddCommand(NewCmdConfigInit(f))

	return configCmd
}

func NewCmdConfigGet(f *cmdutils.Factory) *cobra.Command {
	var hostname string

	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Prints the value of a given configuration key",
		Long:  `Get the value for a given configuration key.`,
		Example: `
  $ glab config get editor
  vim
  $ glab config get glamour_style
  notty
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			val, err := cfg.Get(hostname, args[0])
			if err != nil {
				return err
			}

			if val != "" {
				fmt.Fprintf(utils.ColorableOut(cmd), "%s\n", val)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&hostname, "host", "h", "", "Get per-host setting")
	cmd.Flags().BoolP("global", "g", false, "Read from global config file (~/.config/glab-cli/config.yml). [Default: looks through Environment variables → Local → Global]")

	return cmd
}

func NewCmdConfigSet(f *cmdutils.Factory) *cobra.Command {
	var hostname string

	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Updates configuration with the value of a given key",
		Long: `Update the configuration by setting a key to a value.
Use glab config set --global if you want to set a global config. 
Specifying the --hostname flag also saves in the global config file
`,
		Example: `
  $ glab config set editor vim
  $ glab config set token xxxxx -h gitlab.com
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			localCfg, _ := cfg.Local()

			key, value := args[0], args[1]
			if isGlobal || hostname != "" {
				err = cfg.Set(hostname, key, value)
			} else {
				err = localCfg.Set(key, value)
			}

			if err != nil {
				return fmt.Errorf("failed to set %q to %q: %w", key, value, err)
			}

			if isGlobal || hostname != "" {
				err = cfg.Write()
			} else {
				err = localCfg.Write()
			}

			if err != nil {
				return fmt.Errorf("failed to write config to disk: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&hostname, "host", "h", "", "Set per-host setting")
	cmd.Flags().BoolVarP(&isGlobal, "global", "g", false, "write to global ~/.config/glab-cli/config.yml file rather than the repository .glab-cli/config/config")
	return cmd
}

func NewCmdConfigInit(f *cmdutils.Factory) *cobra.Command {
	var configInitCmd = &cobra.Command{
		Use:   "init",
		Short: "Shows a prompt to set basic glab configuration",
		Long: `Update the configuration by setting a key to a value.
Examples:
  $ glab config init
  ? Enter default Gitlab Host (Current Value: https://gitlab.com): |
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return configInit(cmd, f)
		},
	}
	return configInitCmd
}

func configInit(cmd *cobra.Command, f *cmdutils.Factory) error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}
	baseRepo, err := f.BaseRepo()
	if err != nil {
		return err
	}
	host := baseRepo.RepoHost()
	host, err = config.Prompt(fmt.Sprintf("Enter Gitlab Host (Current Value: %s): ", host), host)
	if err != nil {
		return err
	}
	host, protocol := glinstance.StripHostProtocol(host)
	err = cfg.Set(host, "api_protocol", protocol)
	if err != nil {
		return err
	}

	token, _ := cfg.Get(host, "token")
	token, err = config.Prompt("Enter Gitlab Token: ", token)
	if err != nil {
		return err
	}
	err = cfg.Set(host, "token", token)
	if err != nil {
		return nil
	}

	if cfg.Write() != nil {
		return err
	}
	fmt.Fprintf(utils.ColorableOut(cmd), "%s Configuration updated", utils.GreenCheck())
	return nil
}

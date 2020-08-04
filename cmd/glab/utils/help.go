package utils

import "fmt"

// PrintHelpHelp is exported
func PrintHelpHelp() {

	fmt.Println(`Work seamlessly with Gitlab from the command line.

USAGE
  glab <command> <subcommand> [flags]

CORE COMMANDS
  issue:      Create, view and manage issues
  repo:       Create, manage repositories
  mr:         Create, view, approve and merge merge requests
  pipeline:         view ,delete  piplines

ADDITIONAL COMMANDS
  config:     Manage configuration for glab
  help:       Help about any command

FLAGS
  --help      Show help for command
  --version   Show glab version

EXAMPLES
  $ glab issue create
  $ glab repo clone profclems/glab
  $ glab pr checkout 321

ENVIRONMENT VARIABLES
  GITLAB_TOKEN: an authentication token for API requests. Setting this avoids being
  prompted to authenticate and overrides any previously stored credentials.

  GITLAB_REPO: specify the Gitlab repository in "OWNER/REPO" format for commands that
  otherwise operate on a local repository.

  GITLAB_URI: specify the url of the gitlab server if self hosted (eg: gitlab.example.com)

LEARN MORE
  Use "glab <command> <subcommand> --help" for more information about a command.
  Read the manual at https://glab.clementsam.tech

FEEDBACK
  Open an issue using “glab issue create -R profclems/glab”


Built with ❤ by Clement Sam <https://clementsam.tech>`)

}

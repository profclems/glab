# GLab [Currently BETA]
An open source custom Gitlab Cli tool written in Go (golang) to help work seamlessly with Gitlab from the command line.

## Installation
Download a binary suitable for your OS at https://github.com/profclems/glab/releases.

## Usage
  ```bash
  glab <command> <subcommand> [flags]
  ```

## Core Commands
  ```bash
  issue:      Create, view and manage issues
  repo:       Create, manage repositories
  mr:         Create, view, approve and merge merge requests
  ```

## Additional Commands
  
  ```bash
  config:     Manage configuration for glab
  help:       Help about any command
  ```

## Flags
  ```bash
  --help      Show help for command
  --version   Show glab version
  ```

## Examples
  ```bash
  $ glab issue create
  $ glab issue list --closed
  $ glab repo clone profclems/glab
  $ glab pr checkout 321
  ```

## Envronment Variables
  ```bash
  GITLAB_TOKEN: an authentication token for API requests. Setting this avoids being
  prompted to authenticate and overrides any previously stored credentials.

  GITLAB_REPO: specify the Gitlab repository in "OWNER/REPO" format for commands that
  otherwise operate on a local repository.

  GITLAB_URI: specify the url of the gitlab server if self hosted (eg: gitlab.example.com)
  ```
  
## Learn More
  Use "glab <command> --help" for more information about a command.
  Read the manual at https://glab.clementsam.tech

## Feedback
  Open an issue using `glab issue create -R profclems/glab` to submit an issue on through Gitlab or open a PR on Github


Built with ❤ by Clement Sam <https://clementsam.tech>

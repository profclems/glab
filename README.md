# GLab
GLab open source custom Gitlab Cli tool written in Go (golang) to help work seamlessly with Gitlab from the command line.

## Usage
  ```bash
  glab <command> <subcommand> [flags]
  ```

### Core Commands

- `glab mr [list, create, close, reopen, delete]`
- `glab issue [list, create, close, reopen, delete]`
- `glab config [set]`
- `glab help`


### Examples
  ```bash
  $ glab issue create --title="This is an issue title" --description="This is a really long description"
  $ glab issue list --closed
  ```

## Installation
Download a binary suitable for your OS at https://github.com/profclems/glab/releases/latest.

### Windows
Available as an installable executable file or a Portable archived file in tar and zip formats at the [releases page](https://github.com/profclems/glab/releases/latest).
Download and install now at the [releases page](https://github.com/profclems/glab/releases/latest).

The installable executable file sets the PATH automatically.

### Linux
Download the zip, unzip and install:

1. Download the `.zip` file from the [releases page][]
2. `unzip glab-*-linux-amd64.zip` to unzip the downloaded file 
3. `sudo cp glab-*-linux-amd64/glab /usr/bin` to move to the bin path so you can execute `glab` globally

Or download the tar ball, untar and install:

1. Download the `.tar.gz` file from the [releases page][]
2. `unzip glab-*-linux-amd64.tar.gz` to unzip the downloaded file 
3. `sudo cp glab-*-linux-amd64/glab /usr/bin`

### MacOS
1. Download the `.tar.gz` or `.zip` file from the [releases page][] and unzip or untar
2. ls /usr/local/bin/ || sudo mkdir /usr/local/bin/; to make sure the bin folder exists
3. `sudo cp glab-*-darwin-amd64/glab /usr/bin`

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
  Open an issue on Github


Built with ❤ by Clement Sam <https://clementsam.tech>

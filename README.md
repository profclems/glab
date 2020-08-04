[![Gitter](https://badges.gitter.im/glabcli/community.svg)](https://gitter.im/glabcli/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

# GLab
GLab open source custom Gitlab Cli tool written in Go (golang) to help work seamlessly with Gitlab from the command line.

![image](https://user-images.githubusercontent.com/41906128/88968573-0b556400-d29f-11ea-8504-8ecd9c292263.png)

## Usage
  ```bash
  glab <command> <subcommand> [flags]
  ```

### Core Commands

- `glab mr [list, create, close, reopen, delete]`
- `glab issue [list, create, close, reopen, delete]`
- `glab config`
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

1. Download the `.zip` file from the [releases page](https://github.com/profclems/glab/releases/latest)
2. `unzip glab-*-linux-amd64.zip` to unzip the downloaded file 
3. `sudo mv glab-*-linux-amd64/glab /usr/bin` to move to the bin path so you can execute `glab` globally

Or download the tar ball, untar and install:

1. Download the `.tar.gz` file from the [releases page](https://github.com/profclems/glab/releases/latest)
2. `unzip glab-*-linux-amd64.tar.gz` to unzip the downloaded file 
3. `sudo mv glab-*-linux-amd64/glab /usr/bin`

### MacOS
1. Download the `.tar.gz` or `.zip` file from the [releases page](https://github.com/profclems/glab/releases/latest) and unzip or untar
2. ls /usr/local/bin/ || sudo mkdir /usr/local/bin/; to make sure the bin folder exists
3. `sudo mv glab-*-darwin-amd64/glab /usr/bin`

### Building From Source
If a supported binary for your OS is not found at the [releases page](https://github.com/profclems/glab/releases/latest), you can build from source:

1. Verify that you have Go 1.13.8+ installed

   ```sh
   $ go version
   go version go1.14
   ```

   If `go` is not installed, follow instructions on [the Go website](https://golang.org/doc/install).

2. Clone this repository

   ```sh
   $ git clone https://github.com/profclems/glab.git glab-cli
   $ cd glab-cli
   ```

   or 

   ```sh
   $ git clone https://gitlab.com/profclems/glab.git
   $ cd glab-cli
   ```

3. Build the project

   ```
   $ make build
   ```

4. Move the resulting `bin/glab` executable to somewhere in your PATH

   ```sh
   $ sudo mv ./bin/glab /usr/local/bin/
   ```
   or
   ```sh
   $ sudo mv ./bin/glab /usr/bin/
   ```

4. Run `glab help` to check if it worked.


## Setting Up
After successfull installation, run:
```bash
glab config --token=<YOUR-GITLAB-ACCESS-TOKEN> --url=https://gitlab.com --pid=<YOUR-GITLAB-PROJECT-ID> --repo=OWNER/REPO
```
### Example
```bash
glab config --token=sometoken --url=https://gitlab.com --pid=someprojectid --repo=profclems/glab
```
**NB**: Change gitlab.com to company or group's gitlab url if self-hosted

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
Read the documentation at https://clementsam.tech/glab


## Contributions
Thanks for considering contributing to this project!

Please read the [contibutions guide](.github/contributions.md) and [Code of conduct](.github/CODE_OF_CONDUCT.md). 

Feel free to open an issue or submit a pull request!


## License
[MIT](LICENSE)


## Author
Built with ❤ by [Clement Sam](https://clementsam.tech)

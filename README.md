![GLab](https://user-images.githubusercontent.com/9063085/90530075-d7a58580-e14a-11ea-9727-4f592f7dcf2e.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/profclems/glab)](https://goreportcard.com/report/github.com/profclems/glab)
[![Documentation Status](https://readthedocs.org/projects/glab/badge/?version=latest)](https://glab.readthedocs.io/en/latest/?badge=latest)
[![Gitter](https://badges.gitter.im/glabcli/community.svg)](https://gitter.im/glabcli/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![License](https://img.shields.io/github/license/profclems/glab)](LICENSE)
[![Twitter](https://img.shields.io/twitter/url?style=social&url=https%3A%2F%2Fgithub.com%2Fprofclems%2Fglab)](https://twitter.com/intent/tweet?text=Take%20Gitlab%20to%20the%20command%20line%20with%20%23glab,%20an%20open-source%20GitLab%20CLI%20tool:&url=https%3A%2F%2Fgithub.com%2Fprofclems%2Fglab)


GLab is an open source Gitlab Cli tool written in Go (golang) to help work seamlessly with Gitlab from the command line. Work with issues, merge requests, **watch running pipelines directly from your CLI** among other features. 
Inspired by `gh`, [the official GitHub CLI tool](https://github.com/cli/cli).

![image](https://user-images.githubusercontent.com/41906128/88968573-0b556400-d29f-11ea-8504-8ecd9c292263.png)

## Usage
  ```bash
  glab <command> <subcommand> [flags]
  ```

### Core Commands
```bash
  auth:        Manage glab's authentication state
  issue:       Work with GitLab issues
  label:       Manage labels on remote
  mr:          Create, view and manage merge requests
  pipeline:    Manage pipelines
  release:     Manage GitLab releases
  repo:        Work with GitLab repositories and projects
  
```

### Additional Commands
```bash
  alias:       Create, list and delete aliases
  check-update: Check for latest glab releases
  completion:  Generate shell completion scripts
  config:      Set and get glab settings
  help:        Help about any command
  version:     show glab version information
```

### Examples
  ```bash
  $ glab issue create --title="This is an issue title" --description="This is a really long description"
  $ glab issue list --closed
  $ glab pipeline ci view -b master    # to watch the latest pipeline on master
  $ glab pipeline status    # classic ci view
  ```
  
## Demo
[![asciicast](https://asciinema.org/a/368622.svg)](https://asciinema.org/a/368622)
  
## Learn More
Read the [documentation](https://glab.readthedocs.io/) for more information on this tool.

## Installation
Download a binary suitable for your OS at the [releases page](https://github.com/profclems/glab/releases/latest).

### Quick Install (Bash)
You can install or update `glab` with:
```bash
curl -sL https://j.mp/glab-i | sudo bash
```
or
```bash
curl -s https://raw.githubusercontent.com/profclems/glab/trunk/scripts/quick_install.sh | sudo bash
```
*Installs into `usr/local/bin`*

**NOTE**: Please take care when running scripts in this fashion. Consider peaking at the install script itself and verify that it works as intended.

### Windows
Available for download on [scoop](https://scoop.sh) or manually as an installable executable file or a Portable archived file in tar and zip formats at the [releases page](https://github.com/profclems/glab/releases/latest).

#### Scoop
```sh
scoop bucket add profclems-bucket https://github.com/profclems/scoop-bucket.git
scoop install glab
```
Updating:
```sh
scoop update glab
```

### Linux
Prebuilt binaries available at the [releases page](https://github.com/profclems/glab/releases/latest).

#### Linuxbrew (Homebrew)
```sh
brew install glab
```
Updating:
```sh
brew upgrade glab
```
#### Snapcraft
[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/glab)

Make sure you have snap installed on your Linux Distro (https://snapcraft.io/docs/installing-snapd).
1. `sudo snap install --edge glab`
1. `sudo snap connect glab:ssh-keys` to grant ssh access

#### Arch Linux
`glab` is available through the [gitlab-glab-bin](https://aur.archlinux.org/packages/gitlab-glab-bin/) package on the AUR or download and install an archive from the [releases page](https://github.com/profclems/glab/releases/latest). Arch Linux also supports [snap](https://snapcraft.io/docs/installing-snap-on-arch-linux).
```sh
pacman -Sy gitlab-glab-bin
```

#### KISS Linux
`glab` is available on the [KISS Linux Community Repo](https://github.com/kisslinux/community) as `gitlab-glab`.
If you already have the community repo configured in your `KISS_PATH` you can install `glab` through your terminal.
```sh
kiss b gitlab-glab && kiss i gitlab-glab
```
If you do not have the community repo configured in your `KISS_PATH`, follow the guide on the official guide [Here](https://k1ss.org/install#3.0) to learn how to setup it up.

### MacOS
`glab` is available via Homebrew

#### Homebrew
```sh
brew install glab
```
Updating:
```sh
brew upgrade glab
```

### Building From Source
If a supported binary for your OS is not found at the [releases page](https://github.com/profclems/glab/releases/latest), you can build from source:

#### Prerequisites for building from source are:
- `make`
- Go 1.13+

1. Verify that you have Go 1.13+ installed

   ```sh
   $ go version
   go version go1.14
   ```

   If `go` is not installed, follow instructions on [the Go website](https://golang.org/doc/install).

2. Clone this repository

   ```sh
   $ git clone https://github.com/profclems/glab.git
   $ cd glab
   ```
   If you have $GOPATH/bin or $GOBIN in your $PATH, you can just install with `make install` (install glab in $GOPATH/bin) and **skip steps 3 and 4**.

3. Build the project
   ```
   $ make
   ```

4. Move the resulting `bin/glab` executable to somewhere in your PATH

   ```sh
   $ sudo mv ./bin/glab /usr/local/bin/
   ```

4. Run `glab version` to check if it worked and `glab config init` to set up

## Authentication

Get a GitLab access token at https://gitlab.com/profile/personal_access_tokens or https://gitlab.example.com/profile/personal_access_tokens if self-hosted

- start interactive setup
```sh
$ glab auth login
```

- authenticate against gitlab.com by reading the token from a file
```sh
$ glab auth login --stdin < myaccesstoken.txt
```

- authenticate against a self-hosted GitLab instance by reading from a file
```sh
$ glab auth login --hostname salsa.debian.org --stdin < myaccesstoken.txt
```

- authenticate with token and hostname (Not recommended for shared environments)
```sh
$ glab auth login --hostname gitlab.example.org --token xxxxx
```

## Configuration

`glab` follows the XDG Base Directory [Spec](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html): global configuration file is saved at `~/.config/glab-cli`. Local configuration file is saved at the root of the working git directory and automatically added to `.gitignore`.

**To set configuration globally**

```sh
$ glab config set --global editor vim
```

**To set configuration for current directory (must be a git repository)**

```sh
$ glab config set editor vim
```

**To set configuration for a specific host**

Use the `--host` flag to set configuration for a specific host. This is always stored in the global config file with or without the `global` flag.

```sh
$ glab config set editor vim --host gitlab.example.org
```


## Environment Variables
  ```sh
  GITLAB_TOKEN: an authentication token for API requests. Setting this avoids being
  prompted to authenticate and overrides any previously stored credentials.
  Can be set in the config with 'glab config set token xxxxxx'

  GITLAB_URI or GITLAB_HOST: specify the url of the gitlab server if self hosted (eg: https://gitlab.example.com). Default is https://gitlab.com.

  REMOTE_ALIAS or GIT_REMOTE_URL_VAR: git remote variable or alias that contains the gitlab url.
  Can be set in the config with 'glab config set remote_alias origin'

  VISUAL, EDITOR (in order of precedence): the editor tool to use for authoring text.
  Can be set in the config with 'glab config set editor vim'

  BROWSER: the web browser to use for opening links.
  Can be set in the config with 'glab config set browser mybrowser'

  GLAMOUR_STYLE: environment variable to set your desired markdown renderer style
  Available options are (dark|light|notty) or set a custom style
  https://github.com/charmbracelet/glamour#styles
  
  NO_COLOR: set to any value to avoid printing ANSI escape sequences for color output. 
  ```

## Issues
If you have an issue: report it on the [issue tracker](https://github.com/profclems/glab/issues)

## Contributing
Feel like contributing? That's awesome! We have a [contributing guide](https://github.com/profclems/glab/blob/trunk/.github/CONTRIBUTING.md) and [Code of conduct](https://github.com/profclems/glab/blob/trunk/.github/CODE_OF_CONDUCT.md) to help guide you

### Support `glab` 💖
By donating $5 or more you can support the ongoing development of this project. We'll appreciate some support. Thank you to all our supporters! 🙏 [[Contribute](https://opencollective.com/glab/contribute)]

#### Individuals

This project exists thanks to all the people who contribute. [[Contribute](https://github.com/profclems/glab/blob/trunk/.github/CONTRIBUTING.md)].
<a href="https://opencollective.com/glab/contribute"><img src="https://opencollective.com/glab/contributors.svg?width=890" /></a>

#### Backers

Thank you to all our backers! 🙏 [[Become a backer](https://opencollective.com/glab/contribute)]
<a href="https://opencollective.com/glab#backers" target="_blank"><img src="https://opencollective.com/glab/backers.svg?width=890"></a>

## License
Copyright © [Clement Sam](https://clementsam.tech)

`glab` is open-sourced software licensed under the [MIT](LICENSE) license.

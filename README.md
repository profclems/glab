# GLab

![GLab](https://user-images.githubusercontent.com/9063085/90530075-d7a58580-e14a-11ea-9727-4f592f7dcf2e.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/profclems/glab)](https://goreportcard.com/report/github.com/profclems/glab)
[![Gitter](https://badges.gitter.im/glabcli/community.svg)](https://gitter.im/glabcli/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![Reddit](https://img.shields.io/reddit/subreddit-subscribers/glab_cli?style=social)](https://reddit.com/r/glab_cli)
[![Twitter Follow](https://img.shields.io/twitter/follow/glab_cli?style=social)](https://twitter.com/glab_cli)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#version-control)

GLab is an open source GitLab CLI tool bringing GitLab to your terminal next to where you are already working with `git` and your code without switching between windows and browser tabs. Work with issues, merge requests, **watch running pipelines directly from your CLI** among other features.
Inspired by [gh], the official GitHub CLI tool.

`glab` is available for repositories hosted on GitLab.com and self-hosted GitLab Instances. `glab` supports multiple authenticated GitLab instances and automatically detects the authenticated hostname from the remotes available in the working git directory.

![image](https://user-images.githubusercontent.com/41906128/88968573-0b556400-d29f-11ea-8504-8ecd9c292263.png)

Table of Contents
=================

* [Usage](#usage)
* [Demo](#demo)
* [Documentation](#documentation)
* [Installation](#installation)
  * [Quick Install (Bash)](#quick-install-bash)
  * [Windows](#windows)
    * [WinGet](#winget)
    * [Scoop](#scoop)
    * [EXE Installer](#exe-installer)
  * [Linux](#linux)
    * [Linuxbrew (Homebrew)](#linuxbrew-homebrew)
    * [Snapcraft](#snapcraft)
    * [Arch Linux](#arch-linux)
    * [KISS Linux](#kiss-linux)
    * [Alpine Linux](#alpine-linux)
      * [Install a pinned version from edge](#install-a-pinned-version-from-edge)
      * [Alpine Linux Docker-way](#alpine-linux-docker-way)
    * [Nix/NixOS](#nixnixos)
  * [macOS](#macos)
    * [Homebrew](#homebrew)
    * [MacPorts](#macports)
  * [Building From Source](#building-from-source)
    * [Prerequisites](#prerequisites-for-building-from-source-are)
* [Authentication](#authentication)
* [Configuration](#configuration)
* [Environment Variables](#environment-variables)
* [What about lab](#what-about-lab)
* [Issues](#issues)
* [Contributing](#contributing)
  * [Support glab <g-emoji class="g-emoji" alias="sparkling_heart" fallback-src="https://github.githubassets.com/images/icons/emoji/unicode/1f496.png">💖</g-emoji>](#support-glab-)
    * [Individuals](#individuals)
    * [Backers](#backers)
* [License](#license)

## Usage
  ```bash
  glab <command> <subcommand> [flags]
  ```

## Demo
[![asciicast](https://asciinema.org/a/368622.svg)](https://asciinema.org/a/368622)

## Documentation
Read the [documentation](https://glab.readthedocs.io/) for usage instructions.

## Installation
Download a binary suitable for your OS at the [releases page](https://github.com/profclems/glab/releases/latest).

### Quick Install
**Supported Platforms**: Linux and macOS

#### Homebrew

```sh
brew install glab
```

Updating (Homebrew):
```sh
brew upgrade glab
```

Alternatively, you can install `glab` by shell script:
```sh
curl -sL https://j.mp/glab-cli | sudo sh
```
or
```sh
curl -s https://raw.githubusercontent.com/profclems/glab/trunk/scripts/install.sh | sudo sh
```
*Installs into `usr/bin`*

**NOTE**: Please take care when running scripts in this fashion. Consider peeking at the install script itself and verify that it works as intended.

### Windows
Available for download via [WinGet](https://github.com/microsoft/winget-cli), [scoop](https://scoop.sh), or downloadable EXE installer file.

#### WinGet
```sh
winget install glab
```

Updating (WinGet):
```sh
winget install glab
```

#### Scoop
```sh
scoop install glab
```

Updating (Scoop):
```sh
scoop update glab
```
#### EXE Installer

EXE installers are available for download on the [releases page](https://github.com/profclems/glab/releases/latest).

### Linux
Prebuilt binaries available at the [releases page](https://github.com/profclems/glab/releases/latest).

#### Linuxbrew (Homebrew)
```sh
brew install glab
```

Updating (Homebrew):
```sh
brew upgrade glab
```
#### Snapcraft
[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/glab)

Make sure you have [snap installed on your Linux Distro](https://snapcraft.io/docs/installing-snapd).
1. `sudo snap install --edge glab`
1. `sudo snap connect glab:ssh-keys` to grant ssh access

#### Arch Linux
`glab` is available through the [community/glab](https://archlinux.org/packages/community/x86_64/glab/) package or download and install an archive from the [releases page](https://github.com/profclems/glab/releases/latest). Arch Linux also supports [snap](https://snapcraft.io/docs/installing-snap-on-arch-linux).
```sh
yay -Sy glab
```
or any other [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers) of your choice.

#### KISS Linux
> WARNING: It seems that KISS Linux may no longer be actively maintained, so links to its web domain have been removed from this README.

`glab` is available on the [KISS Linux Community Repo](https://github.com/kisslinux/community) as `gitlab-glab`.
If you already have the community repo configured in your `KISS_PATH` you can install `glab` through your terminal.
```sh
kiss b gitlab-glab && kiss i gitlab-glab
```

#### Alpine Linux

`glab` is available on the [Alpine Community Repo](https://git.alpinelinux.org/aports/tree/community/glab?h=master) as `glab`.

##### Install

We use `--no-cache` so we don't need to do an `apk update` before.

```sh
apk add --no-cache glab
```

##### Install a pinned version from edge

To ensure that by default edge will be used to get the latest updates. We need the edge repository under `/etc/apk/repositories`.

Afterwards you can install it with `apk add --no-cache glab@edge`

We use `--no-cache` so we don't need to do an `apk update` before.

```sh
echo "@edge http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
apk add --no-cache glab@edge
```

##### Alpine Linux Docker-way

Use edge directly

```sh
FROM alpine:3.13
RUN apk add --no-cache glab
```

Fetching latest glab version from edge

```sh
FROM alpine:3.13
RUN echo "@edge http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN apk add --no-cache glab@edge
```

### Nix/NixOS

Nix/NixOS users can install from [nixpkgs](https://search.nixos.org/packages?channel=unstable&show=glab&from=0&size=30&sort=relevance&query=glab):

```bash
nix-env -iA nixos.glab
```

### macOS

#### Homebrew

`glab` is available via [Homebrew](https://formulae.brew.sh/formula/glab)

```sh
brew install glab
```
Updating:
```sh
brew upgrade glab
```

#### MacPorts

`glab`is also available via [MacPorts](https://ports.macports.org/port/glab/summary)

```sh
sudo port install glab
```

Updating:
```sh
sudo port selfupdate && sudo port upgrade glab
```

### Building From Source
If a supported binary for your OS is not found at the [releases page](https://github.com/profclems/glab/releases/latest), you can build from source:

#### Prerequisites for building from source
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
   git clone https://github.com/profclems/glab.git
   cd glab
   ```
   If you have $GOPATH/bin or $GOBIN in your $PATH, you can just install with `make install` (install glab in $GOPATH/bin) and **skip steps 3 and 4**.

3. Build the project
   ```sh
   make
   ```

4. Change PATH to find newly compiled `glab`

   ```sh
   export PATH=$PWD/bin:$PATH
   ```

4. Run `glab version` to confirm that it worked

## Authentication

Get a GitLab access token at <https://gitlab.com/-/profile/personal_access_tokens> or <https://gitlab.example.com/-/profile/personal_access_tokens> if self-hosted

- start interactive setup
  ```sh
  glab auth login
  ```

- authenticate against gitlab.com by reading the token from a file
  ```sh
  glab auth login --stdin < myaccesstoken.txt
  ```

- authenticate against a self-hosted GitLab instance by reading from a file
  ```sh
  glab auth login --hostname salsa.debian.org --stdin < myaccesstoken.txt
  ```

- authenticate with token and hostname (Not recommended for shared environments)
  ```sh
  glab auth login --hostname gitlab.example.org --token xxxxx
  ```

## Configuration

By default, `glab` follows the XDG Base Directory [Spec](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html): global configuration file is saved at `~/.config/glab-cli`. Local configuration file is saved at `.git/glab-cli` in the current working git directory. Advanced workflows may override the location of the global configuration by setting the `GLAB_CONFIG_DIR` environment variable.

**To set configuration globally**

```sh
glab config set --global editor vim
```

**To set configuration for current directory (must be a git repository)**

```sh
glab config set editor vim
```

**To set configuration for a specific host**

Use the `--host` flag to set configuration for a specific host. This is always stored in the global config file with or without the `global` flag.

```sh
glab config set editor vim --host gitlab.example.org
```


## Environment Variables
  ```sh
  GITLAB_TOKEN: an authentication token for API requests. Setting this avoids being
  prompted to authenticate and overrides any previously stored credentials.
  Can be set in the config with 'glab config set token xxxxxx'

  GITLAB_URI or GITLAB_HOST: specify the url of the gitlab server if self hosted (eg: https://gitlab.example.com). Default is https://gitlab.com.

  GITLAB_API_HOST: specify the host where the API endpoint is found. Useful when there are separate [sub]domains or hosts for git and the API endpoint: defaults to the hostname found in the git URL

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

  FORCE_HYPERLINKS: set to 1 to force hyperlinks to be output, even when not outputing to a TTY
  ```

## What about [Lab]?

Both `glab` and [lab] are open-source tools with the same goal of bringing GitLab to your command line and simplifying the developer workflow. In many ways `lab` is to [hub], while `glab` is to [gh].

If you want a tool that’s more opinionated and intended to help simplify your GitLab workflows from the command line, then `glab` is for you. However, If you're looking for a tool like [hub] that feels like using git and allows you to interact with GitLab, you might consider using [lab].

Some `glab` commands such as `ci view` and `ci trace` were adopted from [lab].

[gh]:https://github.com/cli/cli
[hub]:https://github.com/github/hub
[lab]:https://github.com/zaquestion/lab

## Issues
If you have an issue: report it on the [issue tracker](https://github.com/profclems/glab/issues)

## Contributing
Feel like contributing? That's awesome! We have a [contributing guide](https://github.com/profclems/glab/blob/trunk/CONTRIBUTING.md) and [Code of conduct](https://github.com/profclems/glab/blob/trunk/CODE_OF_CONDUCT.md) to help guide you

### Contributors

#### Individuals

This project exists thanks to all the people who contribute. [[Contribute](https://github.com/profclems/glab/blob/trunk/.github/CONTRIBUTING.md)].
<a href="https://opencollective.com/glab/contribute"><img src="https://opencollective.com/glab/contributors.svg?width=890" /></a>

#### Organizations

[![Fosshost.org](https://fosshost.org/img/fosshost-logo.png)](https://fosshost.org)

## License
Copyright © [Clement Sam](https://clementsam.tech)

`glab` is open-sourced software licensed under the [MIT](LICENSE) license.

![GLab](https://user-images.githubusercontent.com/9063085/90530075-d7a58580-e14a-11ea-9727-4f592f7dcf2e.png)
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-19-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

[![Go Report Card](https://goreportcard.com/badge/github.com/profclems/glab)](https://goreportcard.com/report/github.com/profclems/glab)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/profclems/glab/goreleaser)
![.github/workflows/build_docs.yml](https://github.com/profclems/glab/workflows/.github/workflows/build_docs.yml/badge.svg)
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
- `glab mr [list, create, close, reopen, delete, ...]`
- `glab issue [list, create, close, reopen, delete, ...]`
- `glab pipeline [list, delete, ci status, ci view, ...]`
- `glab release`
- `glab repo`
- `glab label`
- `glab alias`


### Examples
  ```bash
  $ glab issue create --title="This is an issue title" --description="This is a really long description"
  $ glab issue list --closed
  $ glab pipeline ci view -b master    # to watch the latest pipeline on master
  $ glab pipeline status    # classic ci view
  ```
  
## Learn More
Read the [documentation](https://clementsam.tech/glab) for more information on this tool.

## Support `glab` 💖
By donating $5 or more you can support the ongoing development of this project. We'll appreciate some support. Thank you to all our supporters! 🙏 [[Contribute](https://opencollective.com/glab/contribute)]

### Individuals

This project exists thanks to all the people who contribute. [[Contribute](https://github.com/profclems/glab/blob/trunk/.github/CONTRIBUTING.md)].
<a href="https://opencollective.com/glab/contribute"><img src="https://opencollective.com/glab/contributors.svg?width=890" /></a>

### Backers

Thank you to all our backers! 🙏 [[Become a backer](https://opencollective.com/glab/contribute)]
<a href="https://opencollective.com/glab#backers" target="_blank"><img src="https://opencollective.com/glab/backers.svg?width=890"></a>

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
Available for download on scoop or manually as an installable executable file or a Portable archived file in tar and zip formats at the [releases page](https://github.com/profclems/glab/releases/latest).
Download and install now at the [releases page](https://github.com/profclems/glab/releases/latest).

The installable executable file sets the PATH automatically.

#### Scoop
```sh
scoop bucket add profclems-bucket https://github.com/profclems/scoop-bucket.git
scoop install glab
```

### Linux
Downloads available via linuxbrew (homebrew) and tar balls

#### Linuxbrew (Homebrew)
```sh
brew install profclems/tap/glab
```
Updating:
```sh
brew upgrade glab
```

#### Arch Linux
`glab` is available through the [gitlab-glab-bin](https://aur.archlinux.org/packages/gitlab-glab-bin/) package on the AUR.

#### Manual Installation
Download the tar ball, untar and install:

1. Download the `.tar.gz` file from the [releases page](https://github.com/profclems/glab/releases/latest)
2. `unzip glab-*-linux-amd64.tar.gz` to unzip the downloaded file 
3. `sudo mv glab-*-linux-amd64/glab /usr/bin`

### MacOS
`glab` is available via Homebrew or you can manually install

#### Homebrew
```sh
brew install profclems/tap/glab
```
Updating:
```sh
brew upgrade glab
```

#### Installing manually
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

4. Run `glab version` to check if it worked and `glab config -g` to set up


## Configuration
Get a GitLab access token at https://gitlab.com/profile/personal_access_tokens or https://gitlab.example.com/profile/personal_access_tokens if self-hosted.
```sh
glab config init # Will be prompted for details and stored in the global directory
```
**To set configuration globally**
```sh
glab config set -g token xxxxxx --host=gitlab.com
```
**To set configuration for current directory (must be a git repository)**
```sh
glab config set token xxxxxx --host=gitlab.com
```

### Example
```sh
glab config 
```
**NB**: Change gitlab.com to company or group's gitlab url (eg. gitlab.example.com) if self-hosted

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
  ```
  
## ToDo
Aside adding more features, the biggest thing this tool still needs is tests 😞

## Issues
If you have an issue: report it on the [issue tracker](https://github.com/profclems/glab/issues)

## Contributing
Feel like contributing? That's awesome! We have a [contributing guide](https://github.com/profclems/glab/blob/trunk/.github/CONTRIBUTING.md) and [Code of conduct](https://github.com/profclems/glab/blob/trunk/.github/CODE_OF_CONDUCT.md) to help guide you

## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://bredley.co.uk"><img src="https://avatars3.githubusercontent.com/u/32489229?v=4" width="100px;" alt=""/><br /><sub><b>Bradley Garrod</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=BreD1810" title="Code">💻</a> <a href="#platform-BreD1810" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/profclems/glab/commits?author=BreD1810" title="Documentation">📖</a></td>
    <td align="center"><a href="https://twitter.com/tetheusmeuneto"><img src="https://avatars2.githubusercontent.com/u/9063085?v=4" width="100px;" alt=""/><br /><sub><b>Matheus Lugon</b></sub></a><br /><a href="#design-matheuslugon" title="Design">🎨</a></td>
    <td align="center"><a href="https://github.com/princeselasi"><img src="https://avatars2.githubusercontent.com/u/59126177?v=4" width="100px;" alt=""/><br /><sub><b>Opoku-Dapaah </b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=princeselasi" title="Documentation">📖</a> <a href="#design-princeselasi" title="Design">🎨</a></td>
    <td align="center"><a href="https://github.com/pgollangi"><img src="https://avatars3.githubusercontent.com/u/6123002?v=4" width="100px;" alt=""/><br /><sub><b>Prasanna Kumar Gollangi</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=pgollangi" title="Code">💻</a> <a href="#maintenance-pgollangi" title="Maintenance">🚧</a></td>
    <td align="center"><a href="https://github.com/sirlatrom"><img src="https://avatars3.githubusercontent.com/u/425633?v=4" width="100px;" alt=""/><br /><sub><b>Sune Keller</b></sub></a><br /><a href="#financial-sirlatrom" title="Financial">💵</a> <a href="https://github.com/profclems/glab/commits?author=sirlatrom" title="Code">💻</a></td>
    <td align="center"><a href="https://sattellite.me"><img src="https://avatars1.githubusercontent.com/u/322910?v=4" width="100px;" alt=""/><br /><sub><b>sattellite</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=sattellite" title="Code">💻</a> <a href="https://github.com/profclems/glab/issues?q=author%3Asattellite" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/abakermi"><img src="https://avatars1.githubusercontent.com/u/60294727?v=4" width="100px;" alt=""/><br /><sub><b>Abdelhak Akermi</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=abakermi" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="http://patrickmcmichael.org"><img src="https://avatars0.githubusercontent.com/u/3779458?v=4" width="100px;" alt=""/><br /><sub><b>Patrick McMichael</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=Saturn" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/wolffc"><img src="https://avatars3.githubusercontent.com/u/1393783?v=4" width="100px;" alt=""/><br /><sub><b>Christian Wolff</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=wolffc" title="Documentation">📖</a></td>
    <td align="center"><a href="https://www.linkedin.com/in/lwpamihiranga/"><img src="https://avatars3.githubusercontent.com/u/39789194?v=4" width="100px;" alt=""/><br /><sub><b>Amith Mihiranga</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=lwpamihiranga" title="Documentation">📖</a></td>
    <td align="center"><a href="https://clementsam.tech"><img src="https://avatars0.githubusercontent.com/u/41906128?v=4" width="100px;" alt=""/><br /><sub><b>Clement Sam</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=profclems" title="Code">💻</a> <a href="#maintenance-profclems" title="Maintenance">🚧</a> <a href="#platform-profclems" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="https://github.com/j-mcavoy"><img src="https://avatars1.githubusercontent.com/u/17990820?v=4" width="100px;" alt=""/><br /><sub><b>John McAvoy</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=j-mcavoy" title="Code">💻</a></td>
    <td align="center"><a href="http://docs.vue2.net"><img src="https://avatars1.githubusercontent.com/u/8638857?v=4" width="100px;" alt=""/><br /><sub><b>wiwi</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=Baiang" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/bgraf"><img src="https://avatars2.githubusercontent.com/u/2063428?v=4" width="100px;" alt=""/><br /><sub><b>Benjamin Graf</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=bgraf" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://qa.debian.org/developer.php?login=ah&comaint=yes"><img src="https://avatars1.githubusercontent.com/u/3367571?v=4" width="100px;" alt=""/><br /><sub><b>andhe</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=andhe" title="Code">💻</a> <a href="#security-andhe" title="Security">🛡️</a></td>
    <td align="center"><a href="https://zacharyspringer.com/"><img src="https://avatars3.githubusercontent.com/u/22923676?v=4" width="100px;" alt=""/><br /><sub><b>Zachary Springer</b></sub></a><br /><a href="#financial-Zachcodes" title="Financial">💵</a></td>
    <td align="center"><a href="https://github.com/zemzale"><img src="https://avatars3.githubusercontent.com/u/14844365?v=4" width="100px;" alt=""/><br /><sub><b>zemzale</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=zemzale" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/anjali-sharma"><img src="https://avatars3.githubusercontent.com/u/13588177?v=4" width="100px;" alt=""/><br /><sub><b>Anjali S.</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=anjali-sharma" title="Code">💻</a> <a href="https://github.com/profclems/glab/issues?q=author%3Aanjali-sharma" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/sunzoje"><img src="https://avatars3.githubusercontent.com/u/3515704?v=4" width="100px;" alt=""/><br /><sub><b>sanjju simha</b></sub></a><br /><a href="https://github.com/profclems/glab/commits?author=sunzoje" title="Code">💻</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!

## License
Copyright © [Clement Sam](https://clementsam.tech)

`glab` is open-sourced software licensed under the [MIT](LICENSE) license.

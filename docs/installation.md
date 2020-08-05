---
layout: page
title: Installation
---

## Installation
Download a binary suitable for your OS at the [releases page](https://github.com/profclems/glab/releases/latest).

### Quick Install (Bash)
```sh
curl -s https://raw.githubusercontent.com/profclems/glab/trunk/scripts/quick_install.sh | sudo bash
```
*Installs into `usr/local/bin`*

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

4. Run `glab help` to check if it worked.

  
## Links
[Issues]({{ '/issues' | absolute_url }})

[Merge Requests]({{ '/mr' | absolute_url }})

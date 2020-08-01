---
layout: page
title: Overview
---

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
    

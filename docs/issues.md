---
layout: default
title: Issues
---

# glab issue
Create and manage issues

## Usage
  ```bash
  glab issue <subcommand> [flags]
  ```

### Sub Commands

- `create`
- `list`
- `close`
- `reopen`
- `delete`
- `subscribe`
- `unsubscribe`


### Examples
  ```bash
  $ glab issue create --title="This is an issue title" --description="This is a really long description" --labels=bug,refactor
  $ glab issue list --closed
  $ glab issue close 34
  $ glab issue reopen 34
  $ glab issue delete 34
  $ glab issue delete 34,56,7 
  $ glab issue unsubscribe 45
  $ glab issue subscribe 45
  ```

## Creating an issue
### Usage
  ```bash
  glab issue create [flags]
  ```

### Flags
  ```bash
  --title           Supply a title. Otherwise, you will prompt for one. (--title="string")
  --description     Supply a description. Otherwise, you will prompt for one. (--description="string")
  --label           Add label by name. Multiple labels should be comma separated. Otherwise, you will prompt for one, though optional (--label="string,string")
  --assigns         Assign issue to people by their ID. Multiple values should be comma separated (--assigned=value,value)
  --milestone       Add the issue to a milestone by id. (--milestone=value)
  --confidential    Set issue as confidential. Optional boolean value (--confidential) or (--confidential=true)
  --mr              Link issue to a merge request by ID. (--mr=id)
  --weight          Set weight of issue
  --epic          
  ```

## Installation
[Installation Guide]({{ '/installation' | absolute_url }})
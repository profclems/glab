---
layout: page
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
--assigns         Assign issue to people by their ID. Multiple values should be comma separated (--assigns=value,value)
--milestone       Add the issue to a milestone by id. (--milestone=value)
--confidential    Set issue as confidential. Optional boolean value (--confidential) or (--confidential=true)
--mr              Link issue to a merge request by ID. (--mr=id)
--weight          Set weight of issue
--epic          
  ```

## Listing Issues
### Usage
  ```bash
glab issue list [flags]
  ```
#### Alias: `ls`

### Flags
  ```bash
--all             Show all opened and closed issues
--closed          Get the list of closed issues
--opened          Get all opened issues (default)            
  ```
### Example
```sh
# Get all opened issues
glab issue list

# Get closed issues
glab issue list --closed

# Get all issues
glab issue list --all
```

## Closing Issues
### Usage
To close a single issue
  ```bash
glab issue close <issueId>
  ```
To close multiple issues
  ```bash
glab issue close <comma,separated,ids>
  ```

### Example
```sh
glab issue close 23
glab issue close 67,34,21
```

## Reopening Issues
### Usage
To reopen a single issue
  ```bash
glab issue reopen <issueId>
  ```
To reopen multiple issues
  ```bash
glab issue reopen <comma,separated,ids>
  ```

### Example
```sh
glab issue reopen 23
glab issue reopen 67,34,21
```

## Subscribe to an Issues
### Usage
To subscribe to a single issue
  ```bash
glab issue unsubscribe<issueId>
  ```
To subscribe to multiple issues
  ```bash
glab issue subscribe <comma,separated,ids>
  ```

### Example
```sh
glab issue subscribe 23
glab issue subscribe 67,34,21
```

## Unsubscribe to an Issues
### Usage
To unsubscribe to a single issue
  ```bash
glab issue unsubscribe <issueId>
  ```
To unsubscribe to multiple issues
  ```bash
glab issue unsubscribe <comma,separated,ids>
  ```

### Example
```sh
glab issue unsubscribe 23
glab issue unsubscribe 67,34,21
```

## Deleting Issues
### Usage
To delete to a single issue
  ```bash
glab issue delete <issueId>
  ```
To delete to multiple issues
  ```bash
glab issue delete <comma,separated,ids>
  ```

### Example
```sh
glab issue delete 23
glab issue delete 67,34,21
```

## Links
[Installation Guide]({{ '/installation' | absolute_url }})

[Managing Merge Requests]({{ '/mr' | absolute_url }})
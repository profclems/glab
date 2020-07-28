---
layout: page
title: Merge Requests
---

# glab mr
Create and manage merge requests

## Usage
  ```bash
glab mr <subcommand> [flags]
  ```

### Sub Commands

- `create`
- `merge`, `accept`
- `list`
- `close`
- `reopen`
- `delete`
- `subscribe`
- `unsubscribe`


### Examples
  ```bash
$ glab mr create --title="This is an merge request title" --description="This is a really long description" --labels=bug,refactor
$ glab mr list --closed
$ glab mr merge 45
$ glab mr close 34
$ glab mr reopen 34
$ glab mr delete 34
$ glab mr delete 34,56,7 
$ glab mr unsubscribe 45
$ glab mr subscribe 45
  ```

## Creating merge requests
### Usage
  ```bash
glab mr create [flags]
  ```

### Flags
  ```bash
--title           Supply a title. Otherwise, you will prompt for one. (--title="string")
--description     Supply a description. Otherwise, you will prompt for one. (--description="string")
--source          Supply the source branch. Otherwise, you will prompt for one. (--source="string")
--target          Supply the target branch. Otherwise, you will prompt for one. (--target="string")
--label           Add label by name. Multiple labels should be comma separated. Otherwise, you will prompt for one, though optional (--label="string,string")
--assigns         Assign merge request to people by their ID. Multiple values should be comma separated (--assigns=value,value)
--milestone       Add the merge request to a milestone by id. (--milestone=value)
--weight          Set weight of merge request
--epic          
--allow-collaboration
--remove-source-branch

  ```
## Merging/Approving Merge Requests
### Usage
  ```bash
glab mr merge <mrID>         
  ```
### Example
```sh
glab mr merge 56
```

## Listing Merge Requests
### Usage
  ```bash
glab mr list [flags]
  ```
#### Alias: `ls`

### Flags
  ```bash
--all             Show all opened and closed merge requests
--closed          Get the list of closed merge requests
--opened          Get all opened merge requests (default)            
  ```
### Example
```sh
# Get all opened merge requests
glab mr list

# Get closed merge requests
glab mr list --closed

# Get all merge requests
glab mr list --all
```

## Closing Merge Requests
### Usage
To close a single merge request
  ```bash
glab mr close <mergeRequestId>
  ```
To close multiple merge requests
  ```bash
glab mr close <comma,separated,ids>
  ```

### Example
```sh
glab mr close 23
glab mr close 67,34,21
```

## Reopening Merge Requests
### Usage
To reopen a single merge request
  ```bash
glab mr reopen <mergeRequestId>
  ```
To reopen multiple merge requests
  ```bash
glab mr reopen <comma,separated,ids>
  ```

### Example
```sh
glab mr reopen 23
glab mr reopen 67,34,21
```

## Subscribe to Merge Requests
### Usage
To subscribe to a single merge request
  ```bash
glab mr unsubscribe<mergeRequestId>
  ```
To subscribe to multiple merge requests
  ```bash
glab mr subscribe <comma,separated,ids>
  ```

### Example
```sh
glab mr subscribe 23
glab mr subscribe 67,34,21
```

## Unsubscribe to Merge Requests
### Usage
To unsubscribe to a single merge request
  ```bash
glab mr unsubscribe <mergeRequestId>
  ```
To unsubscribe to multiple merge requests
  ```bash
glab mr unsubscribe <comma,separated,ids>
  ```

### Example
```sh
glab mr unsubscribe 23
glab mr unsubscribe 67,34,21
```

## Deleting Merge Requests
### Usage
To delete to a single merge request
  ```bash
glab mr delete <mergeRequestId>
  ```
To delete to multiple merge requests
  ```bash
glab mr delete <comma,separated,ids>
  ```

### Example
```sh
glab mr delete 23
glab mr delete 67,34,21
```

## Links
[Installation Guide]({{ '/installation' | absolute_url }})

[Managing Issues]({{ '/issues' | absolute_url }})
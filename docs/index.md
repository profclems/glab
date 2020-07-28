# GLab
GLab open source custom Gitlab Cli tool written in Go (golang) to help work seamlessly with Gitlab from the command line.

![image](https://user-images.githubusercontent.com/41906128/88602028-613cc880-d061-11ea-84c1-71b6e1e02611.png)

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

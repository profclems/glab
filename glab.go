package glab

import (
	"fmt"
	"strings"

	"glab/cmd/glab/utils"
	"glab/commands"

	"github.com/logrusorgru/aurora"
)

var (
	version string
	build   string
	commit  string
)

func printVersion(_ map[string]string, _ map[int]string) {
	fmt.Println()
	fmt.Println("GLab version", version)
	fmt.Println("Build:", build)
	fmt.Println("Commit:", commit)
	fmt.Println("https://github.com/profclems/glab")
	fmt.Println()
	fmt.Println("Made with ❤ by Clement Sam <https://clementsam.tech>")
	fmt.Println()
}

// Help is exported
func Help(args map[string]string, arrCmd map[int]string) {
	utils.PrintHelpHelp()
}

func config(cmdArgs map[string]string, arrCmd map[int]string) {
	cmdHelpList := map[string]string{
		"uri":        "GITLAB_URI",
		"url":        "GITLAB_URI",
		"token":      "GITLAB_TOKEN",
		"repo":       "GITLAB_REPO",
		"pid":        "GITLAB_PROJECT_ID",
		"remote-var": "GIT_REMOTE_URL_VAR",
		"origin":     "GIT_REMOTE_URL_VAR",
		"origin-var": "GIT_REMOTE_URL_VAR",
	}

	commands.UseGlobalConfig = true
	if commands.VariableExists("GITLAB_URI") == "NOTFOUND" || commands.VariableExists("GITLAB_URI") == "OK" {
		commands.SetEnv("GITLAB_URI", "https://gitlab.com")
	}
	if commands.VariableExists("GIT_REMOTE_URL_VAR") == "NOTFOUND" || commands.VariableExists("GIT_REMOTE_URL_VAR") == "OK" {
		commands.SetEnv("GIT_REMOTE_URL_VAR", "origin")
	}
	commands.UseGlobalConfig = false

	var isUpdated bool
	if arrCmd[0] == "global" {
		commands.UseGlobalConfig = true
	}
	for i := 0; i < len(arrCmd); i++ {
		if commands.CommandArgExists(cmdArgs, arrCmd[i]) && commands.CommandArgExists(cmdHelpList, arrCmd[i]) {
			commands.SetEnv(cmdHelpList[arrCmd[i]], cmdArgs[arrCmd[i]])
			isUpdated = true
		} else if arrCmd[0] != "global" {
			fmt.Println(aurora.Red(arrCmd[i] + ": invalid flag"))
		}
	}

	if isUpdated {
		fmt.Println(aurora.Green("Environment variable(s) updated"))
	}
}

// Exec is exported
func Exec(cmd string, cmdArgs map[string]string, arrCmd map[int]string) {
	commandList := map[string]func(map[string]string, map[int]string){
		"issue":     commands.ExecIssue,
		"mr":        commands.ExecMergeRequest,
		"label":     commands.ExecLabel,
		"pipeline":  commands.ExecPipeline,
		"repo":      commands.ExecRepo,
		"help":      Help,
		"config":    config,
		"version":   printVersion,
		"--version": printVersion,
		"-v":        printVersion,
	}
	cmd = strings.Trim(cmd, " ")
	if cmd == "" {
		Help(cmdArgs, arrCmd)
	}
	if commands.CommandExists(commandList, cmd) {

		if len(cmdArgs) == 1 {
			if cmdArgs["help"] == "true" {
				cmdHelpList := map[string]func(){
					"help":  utils.PrintHelpHelp,
					"issue": utils.PrintHelpIssue,
					"mr":    utils.PrintHelpMr,
					"repo":  utils.PrintHelpRepo,
					"pipeline": utils.PrintHelpPipeline,
				}
				cmdHelpList[cmd]()
				return
			}
		}
		commandList[cmd](cmdArgs, arrCmd)
	} else {
		fmt.Println(cmd + ": command not found")
		fmt.Println()
		Help(cmdArgs, arrCmd)
	}
}

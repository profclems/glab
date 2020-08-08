package glab

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

func readAndTrimString(defaultVal string) string  {
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')
	str = strings.TrimSuffix(str, "\n")
	if str == "" && defaultVal != "" {
		return defaultVal
	}
	return str
}

func readAndSetEnv(env string) string  {
	var envDefVal string
	if env == "" {
		envDefVal = commands.GetEnv(env)
	}
	envVal := readAndTrimString(envDefVal)
	commands.SetEnv(env, envVal)
	return envVal
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
			if arrCmd[0] != "unset" {
				commands.SetEnv(cmdHelpList[arrCmd[i]], "")
			} else {
				commands.SetEnv(cmdHelpList[arrCmd[i]], cmdArgs[arrCmd[i]])
			}
			isUpdated = true
		} else if arrCmd[0] != "global" && arrCmd[0] != "init" {
			fmt.Println(aurora.Red(arrCmd[i] + ": invalid flag"))
		}
	}
	if !isUpdated {
		fmt.Printf("Enter default Gitlab Host (Current Value: %s): ", commands.GetEnv("GITLAB_URI"))
		readAndSetEnv("GITLAB_URI")
		fmt.Print("Enter default Gitlab Token: ")
		readAndSetEnv("GITLAB_TOKEN")
		fmt.Printf("Enter Git remote url variable (Current Value: %s): ", commands.GetEnv("GIT_REMOTE_URL_VAR"))
		readAndSetEnv("GIT_REMOTE_URL_VAR")
	}
	fmt.Println(aurora.Green("Environment variable(s) updated"))
}

// Exec is exported
func Exec(cmd string, cmdArgs map[string]string, arrCmd map[int]string) {
	commandList := map[string]func(map[string]string, map[int]string){
		"issue":     commands.ExecIssue,
		"mr":        commands.ExecMergeRequest,
		"label":     commands.ExecLabel,
		"pipeline":  commands.ExecPipeline,
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

		if len(cmdArgs) > 0 {
			if cmdArgs["help"] == "true" {
				cmdHelpList := map[string]func(){
					"help":  utils.PrintHelpHelp,
					"issue": utils.PrintHelpIssue,
					"mr":    utils.PrintHelpMr,
					"repo":  utils.PrintHelpRepo,
				}
				if _, ok := cmdHelpList[cmd]; ok {
					cmdHelpList[cmd]()
				} else {
					log.Fatal("Invalid command")
				}
			}
		}
		commandList[cmd](cmdArgs, arrCmd)
	} else {
		log.Fatal(cmd + ": command not found")
	}
}

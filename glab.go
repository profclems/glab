package glab

import (
	"bufio"
	"fmt"
	"glab/commands"
	"log"
	"os"
	"strings"
)

func Version(_ map[string]string, _ map[int]string)  {
	version := "v0.1.0"
	fmt.Println("GLab version", version)
	fmt.Println("https://github.com/profclems/glab")
	fmt.Println()
	fmt.Println("Made with ❤ by Clement Sam <https://clementsam.tech>")
	fmt.Println()
}

func OpenFile(filename string)  {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func Issue(cmdArgs map[string]string, arrCmd map[int]string)  {
	commands.ExecIssue(cmdArgs, arrCmd)
	return
}

func MergeRequest(cmdArgs map[string]string, arrCmd map[int]string)  {
	commands.MakeRequest(`{}`,"projects/20131402/issues/1","GET")
	return
}

func Help(args map[string]string, arrCmd map[int]string) {
	OpenFile("./utils/help.txt")
	return
}

func Config(cmdArgs map[string]string, arrCmd map[int]string)  {
	if commands.CommandArgExists(cmdArgs, "uri") {
		commands.SetEnv("GITLAB_URI", cmdArgs["uri"])
	}
	if commands.CommandArgExists(cmdArgs, "uri") {
		commands.SetEnv("GITLAB_TOKEN", cmdArgs["token"])
	}
	if commands.CommandArgExists(cmdArgs, "uri") {
		commands.SetEnv("GITLAB_REPO", cmdArgs["repo"])
	}
	fmt.Println("Environment variable(s) updated")
}

func Exec(cmd string, cmdArgs map[string]string, arrCmd map[int]string)  {
	commandList := map[string]func(map[string]string, map[int]string) {
		"issue": Issue,
		"mr" : MergeRequest,
		"help" : Help,
		"config" : Config,
		"version" : Version,
		"--version" : Version,
		"-v" : Version,
	}
	cmd = strings.Trim(cmd, " ")
	if cmd == "" {
		Help(cmdArgs, arrCmd)
		return
	}
	if commands.CommandExists(commandList, cmd) {

		if len(cmdArgs)>0 {
			if cmdArgs["help"] == "true" {
				OpenFile("./utils/"+cmd+".txt")
				return
			}
		}
		commandList[cmd](cmdArgs, arrCmd)
	} else {
		fmt.Println(cmd, ":command not found")
		fmt.Println()
		Help(cmdArgs, arrCmd)
	}
	return
}

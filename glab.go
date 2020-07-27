package glab

import (
	"bufio"
	"fmt"
	. "github.com/logrusorgru/aurora"
	"glab/cmd/glab/utils"
	"glab/commands"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func Version(_ map[string]string, _ map[int]string)  {
	version := "v1.5.0"
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
	commands.ExecMergeRequest(cmdArgs, arrCmd)
	return
}

func Help(args map[string]string, arrCmd map[int]string) {
	//OpenFile("./utils/help.go")
	utils.PrintHelpHelp()
	return
}

func ConfigEnv(key string, value string)  {
	data, _ := ioutil.ReadFile("./config/.env")

	file := string(data)
	line := 0
	temp := strings.Split(file, "\n")
	newData := ""
	keyExists := false
	newConfig := key+"="+(value)+"\n"
	for _, item := range temp {
		//fmt.Println("[",line,"]",item)
		env := strings.Split(item, "=")
		justString := fmt.Sprint(item)
		if env[0] == key {
			newData += newConfig
			keyExists = true
		} else {
			newData += justString + "\n"
		}
		line++
	}
	if !keyExists {
		newData += newConfig
	}
	_ = os.Mkdir("./config", 0700)
	f, _ := os.Create("./config/.env")// Create a writer
	w := bufio.NewWriter(f)
	_, _ = w.WriteString(strings.Trim(newData, "\n"))
	_ = w.Flush()
}

func Config(cmdArgs map[string]string, arrCmd map[int]string)  {
	cmdHelpList := map[string]string {
		"uri" : "GITLAB_URI",
		"url" : "GITLAB_URI",
		"token" : "GITLAB_TOKEN",
		"repo" : "GITLAB_REPO",
		"pid" : "GITLAB_PROJECT_ID",
	}
	isUpdated := false

	fmt.Println() //Upper Space
	for i:=0; i < len(arrCmd); i++ {
		if commands.CommandArgExists(cmdArgs, arrCmd[i]) && commands.CommandArgExists(cmdHelpList, arrCmd[i]) {
			ConfigEnv(cmdHelpList[arrCmd[i]], cmdArgs[arrCmd[i]])
			isUpdated = true
		} else {
			fmt.Println(Red(arrCmd[i]+": command not found"))
		}
	}

	if isUpdated {
		fmt.Println(Green("Environment variable(s) updated"))
	}
	fmt.Println() //ending space
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
				cmdHelpList := map[string]func() {
					"help" : utils.PrintHelpHelp,
					"issue" : utils.PrintHelpIssue,
					"mr" : utils.PrintHelpMr,
					"repo" : utils.PrintHelpRepo,
				}
				//OpenFile("./utils/"+cmd+".txt")
				cmdHelpList[cmd]()
				return
			}
		}
		commandList[cmd](cmdArgs, arrCmd)
	} else {
		fmt.Println(cmd+": command not found")
		fmt.Println()
		Help(cmdArgs, arrCmd)
	}
	return
}

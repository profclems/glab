package main

import (
	"fmt"
	"glab"
	"glab/commands"
	"os"
	"strings"
)

func main() {
	// take the command and arguments
	args := os.Args
	cmdArgs := args[1:]

	arr := make(map[string]string)
	arrCmd := make(map[int]string)

	// check if command was passed
	if len(cmdArgs) == 0 {
		glab.Help(arr, arrCmd)
		return
	}

	cmd := args[1] //Get the command
	argLen := len(cmdArgs)

	// Parse the arguments in a map
	for i := 1; i < argLen; i++ {
		sp := strings.Split(strings.TrimLeft(cmdArgs[i], "-"), "=") // Trim preceding dashes and split key and value into map
		if sp[0] == "h" {
			sp[0] = "help"
		}
		if len(sp) > 0 {
			if len(sp) > 1 {
				arr[sp[0]] = sp[1]
			} else {
				arr[sp[0]] = "true" // set flags with empty values to true
			}
			arrCmd[(i - 1)] = sp[0]
		}
	}
	// Set Global Env Path
	commands.SetGlobalPathDir()
	fmt.Println() //Upper Space
	// Execute Command
	glab.Exec(cmd, arr, arrCmd)
	fmt.Println() //ending space
}

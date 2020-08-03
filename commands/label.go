package commands

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
	"strings"
)

func createLabel(cmdArgs map[string]string, _ map[int]string) {
	reader := bufio.NewReader(os.Stdin)
	var labelTitle string
	var labelColor string
	if !CommandArgExists(cmdArgs, "name") {
		fmt.Print(aurora.Cyan("Name" + "\n" + "-> "))
		labelTitle, _ = reader.ReadString('\n')
	} else {
		labelTitle = strings.Trim(cmdArgs["title"], " ")
	}
	if !CommandArgExists(cmdArgs, "color") {
		fmt.Print(aurora.Cyan("Color" + "\n" + "-> "))
		labelColor, _ = reader.ReadString('\n')
	} else {
		labelColor = strings.Trim(cmdArgs["label"], "[] ")
	}
	git, repo := InitGitlabClient()
	// Create new label
	l := &gitlab.CreateLabelOptions{
		Name:  gitlab.String(labelTitle),
		Color: gitlab.String(labelColor),
	}
	label, _, err := git.Labels.CreateLabel(repo, l)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created label: %s\nWith color: %s\n", label.Name, label.Color)
}


func listLabels(cmdArgs map[string]string, _ map[int]string) {
	git, repo := InitGitlabClient()
	// List all labels
	labels, _, err := git.Labels.ListLabels(repo, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Showing label %d of %d on %s", len(labels), len(labels), repo)
	fmt.Println()
	for _, label := range labels {
		fmt.Println(label.Name)
	}
}

// ExecRepo is ...
func ExecLabel(cmdArgs map[string]string, arrCmd map[int]string) {
	commandList := map[interface{}]func(map[string]string, map[int]string){
		"create":      	createLabel,
		"new":     		createLabel,
		"list":        	listLabels,
	}
	if _, ok := commandList[arrCmd[0]]; ok {
		commandList[arrCmd[0]](cmdArgs, arrCmd)
	} else {
		fmt.Println(arrCmd[0]+":", "Invalid Command")
	}
}


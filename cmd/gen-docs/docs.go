package main

import (
	"log"

	"github.com/profclems/glab/commands"

	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(commands.RootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}

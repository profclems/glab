package main

import (
	"log"

	"github.com/spf13/cobra/doc"
	"glab/commands"
)

func main() {
	err := doc.GenMarkdownTree(commands.RootCmd, "./gen-docs")
	if err != nil {
		log.Fatal(err)
	}
}

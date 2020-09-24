package main

import (
	"log"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/root"

	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(root.NewCmdRoot(&cmdutils.Factory{}, "", ""), "./docs")
	if err != nil {
		log.Fatal(err)
	}
}

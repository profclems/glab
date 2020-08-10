package main

import (
	"github.com/spf13/cobra/doc"
	"glab/commands"
	"log"
	"path"
	"strings"
)

/*
const fmTemplate = `---
layout: page
title: "%s"
---
`
 */

func main() {
	filePrepender := func(filename string) string {
		/*
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		return fmt.Sprintf(fmTemplate, strings.Replace(base, "_", " ", -1))
		 */
		return ""
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "/" + strings.ToLower(base) + "/"
	}

	err := doc.GenMarkdownTreeCustom(commands.RootCmd, "./docs", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}

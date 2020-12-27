package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/profclems/glab/commands"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/utils"
	"github.com/spf13/cobra"
)

var tocTree = `.. toctree::
   :glob:
   :maxdepth: 0

%s

`

func main() {
	//err := doc.GenMarkdownTree(commands.NewCmdRoot(&cmdutils.Factory{}, "", ""), "./docs")
	//if err != nil {
	//	log.Fatal(err)
	//}

	docLoc := "./docs/source/"
	var err error

	ioS, _, _, _ := utils.IOTest()
	glabCli := commands.NewCmdRoot(&cmdutils.Factory{IO: ioS}, "", "")
	glabCli.DisableAutoGenTag = true
	cmds := glabCli.Commands()

	for _, cmd := range cmds {
		fmt.Println("Generating docs for " + cmd.Name())
		// create directories for parent commands
		_ = os.MkdirAll(docLoc+cmd.Name(), 0750)

		// Generate parent command
		out := new(bytes.Buffer)
		err = GenReSTCustom(cmd, out)
		if err != nil {
			fatal(err)
		}

		// Generate children commands
		for _, cmdC := range cmd.Commands() {
			err = GenReSTTreeCustom(cmdC, docLoc+cmd.Name())
			if err != nil {
				fatal(err)
			}
		}

		err = config.WriteFile(docLoc+cmd.Name()+"/index.rst", out.Bytes(), 0755)
		if err != nil {
			fatal(err)
		}
	}
}

func printSubcommands(cmd *cobra.Command, buf *bytes.Buffer) {
	if len(cmd.Commands()) < 1 {
		return
	}

	var tree string
	// Generate children commands
	for _, cmdC := range cmd.Commands() {
		if cmdC.Name() != "help" {
			tree += fmt.Sprintf("%s <%s>\n", cmdC.Name(), cmdC.Name())
		}
	}

	subcommands := ""
	if tree != "" {
		tree = utils.Indent(tree, "   ")
		subcommands = fmt.Sprintf(tocTree, tree)
	}

	if subcommands != "" {
		buf.WriteString("Subcommands\n")
		buf.WriteString("~~~~~~~~~~~\n")
		buf.WriteString(subcommands)
		buf.WriteString("\n")
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

// adapted from : github.com/spf13/cobra/doc/rest_docs.go

// GenReSTTreeCustom is the the same as GenReSTTree, but
// with custom filePrepender and linkHandler.
func GenReSTTreeCustom(cmd *cobra.Command, dir string) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := GenReSTTreeCustom(c, dir); err != nil {
			return err
		}
	}

	basename := cmd.Name() + ".rst"
	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := GenReSTCustom(cmd, f); err != nil {
		return err
	}
	return nil
}

// GenReSTCustom creates custom reStructured Text output.
func GenReSTCustom(cmd *cobra.Command, w io.Writer) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	short := cmd.Short
	long := cmd.Long
	if len(long) == 0 {
		long = short
	}
	ref := strings.Replace(name, " ", "_", -1)

	buf.WriteString(".. _" + ref + ":\n\n")
	buf.WriteString(name + "\n")
	buf.WriteString(strings.Repeat("-", len(name)) + "\n\n")
	buf.WriteString(short + "\n\n")
	buf.WriteString("Synopsis\n")
	buf.WriteString("~~~~~~~~\n\n")
	buf.WriteString("\n" + long + "\n\n")

	if cmd.Runnable() {
		buf.WriteString(fmt.Sprintf("::\n\n  %s\n\n", cmd.UseLine()))
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("Examples\n")
		buf.WriteString("~~~~~~~~\n\n")
		buf.WriteString(fmt.Sprintf("::\n\n%s\n\n", utils.Indent(cmd.Example, "  ")))
	}

	if err := printOptionsReST(buf, cmd, name); err != nil {
		return err
	}

	printSubcommands(cmd, buf)

	_, err := buf.WriteTo(w)
	return err
}

// adapted from : github.com/spf13/cobra/doc/rest_docs.go
func printOptionsReST(buf *bytes.Buffer, cmd *cobra.Command, name string) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)
	if flags.HasAvailableFlags() {
		buf.WriteString("Options\n")
		buf.WriteString("~~~~~~~\n\n::\n\n")
		flags.PrintDefaults()
		buf.WriteString("\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)
	if parentFlags.HasAvailableFlags() {
		buf.WriteString("Options inherited from parent commands\n")
		buf.WriteString("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n\n::\n\n")
		parentFlags.PrintDefaults()
		buf.WriteString("\n")
	}
	return nil
}

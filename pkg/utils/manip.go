package utils

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/profclems/glab/pkg/prompt"

	"github.com/AlecAivazis/survey/v2"
)

type EditorOptions struct {
	FileName      string
	Label         string
	Help          string
	Default       string
	AppendDefault bool
	HideDefault   bool
}

func Editor(opts EditorOptions) string {
	var container string
	editor := &survey.Editor{
		Renderer:      survey.Renderer{},
		Message:       opts.Label,
		Default:       opts.Default,
		Help:          opts.Help + "Uses the editor defined by the $VISUAL or $EDITOR environment variables). If neither of those are present, notepad (on Windows) or vim (Linux or Mac) is used",
		HideDefault:   opts.HideDefault,
		AppendDefault: opts.AppendDefault,
		FileName:      opts.FileName,
	}
	err := prompt.AskOne(editor, &container)
	if err != nil {
		log.Fatal(err)
	}
	return container
}

// ReplaceNonAlphaNumericChars : Replaces non alpha-numeric values with provided char/string
func ReplaceNonAlphaNumericChars(words, replaceWith string) string {
	reg := regexp.MustCompile("[^A-Za-z0-9]+")
	newStr := reg.ReplaceAllString(strings.Trim(words, " "), replaceWith)
	return newStr
}

func StringToInt(str string) int {
	strInt, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return strInt
}

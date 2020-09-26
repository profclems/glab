package utils

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	
	"github.com/AlecAivazis/survey/v2"
)

func AskQuestionWithInput(question, defaultVal string, isRequired bool) string {
	str := ""
	prompt := &survey.Input{
		Message: question,
	}
	var err error
	if isRequired {
		err = survey.AskOne(prompt, &str, survey.WithValidator(survey.Required))
	} else {
		err = survey.AskOne(prompt, &str)
	}
	if err != nil {
		log.Fatal(err)
	}
	str = strings.TrimSuffix(str, "\n")
	if str == "" && defaultVal != "" {
		return defaultVal
	}
	return str
}

// Confirm prompts user for a confirmation and returns a bool value
func Confirm(question string) (confirmed bool, err error) {
	confirmed = false
	prompt := &survey.Confirm{
		Message: question,
	}
	err = survey.AskOne(prompt, &confirmed)
	return
}

func AskQuestionWithMultiSelect(question string, options []string) []string {
	labels := []string{}
	prompt := &survey.MultiSelect{
		Message: question,
		Options: options,
	}
	err := survey.AskOne(prompt, &labels)
	if err != nil {
		log.Fatal(err)
	}
	return labels
}

func AskQuestionMultiline(question string, defaultVal string) string {
	str := ""
	prompt := &survey.Multiline{
		Message: question,
	}
	err := survey.AskOne(prompt, &str)
	if err != nil {
		log.Fatal(err)
	}
	str = strings.TrimSuffix(str, "\n")
	if str == "" && defaultVal != "" {
		return defaultVal
	}
	return str
}

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
	prompt := &survey.Editor{
		Renderer:      survey.Renderer{},
		Message:       opts.Label,
		Default:       opts.Default,
		Help:          opts.Help + "Uses the editor defined by the $VISUAL or $EDITOR environment variables). If neither of those are present, notepad (on Windows) or vim (Linux or Mac) is used",
		HideDefault:   opts.HideDefault,
		AppendDefault: opts.AppendDefault,
		FileName:      opts.FileName,
	}
	err := survey.AskOne(prompt, &container)
	if err != nil {
		log.Fatal(err)
	}
	return container
}

// ReplaceNonAlphaNumericChars : Replaces non alpha-numeric values with provided char/string
func ReplaceNonAlphaNumericChars(words, replaceWith string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
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

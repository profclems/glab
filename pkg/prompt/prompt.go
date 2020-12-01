package prompt

import (
	"github.com/AlecAivazis/survey/v2"
)

func AskQuestionWithInput(response interface{}, question, defaultVal string, isRequired bool) error {
	prompt := &survey.Input{
		Message: question,
		Default: defaultVal,
	}
	var err error
	if isRequired {
		err = survey.AskOne(prompt, response, survey.WithValidator(survey.Required))
	} else {
		err = survey.AskOne(prompt, response)
	}
	if err != nil {
		return err
	}
	return nil
}

func MultiSelect(response interface{}, question string, options []string, opts ...survey.AskOpt) error {
	prompt := &survey.MultiSelect{
		Message: question,
		Options: options,
	}
	err := AskOne(prompt, response, opts...)
	if err != nil {
		return err
	}
	return nil
}

func AskMultiline(response interface{}, question string, defaultVal string) error {
	prompt := &survey.Multiline{
		Message: question,
		Default: defaultVal,
	}
	err := survey.AskOne(prompt, response)
	if err != nil {
		return err
	}
	return nil
}

var AskOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	return survey.AskOne(p, response, opts...)
}

var Ask = func(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
	return survey.Ask(qs, response, opts...)
}

var Confirm = func(result *bool, prompt string, defaultVal bool) error {
	p := &survey.Confirm{
		Message: prompt,
		Default: defaultVal,
	}
	return survey.AskOne(p, result)
}

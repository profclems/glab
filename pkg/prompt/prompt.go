package prompt

import (
	"github.com/AlecAivazis/survey/v2"
)

func AskQuestionWithInput(response interface{}, name, question, defaultVal string, isRequired bool) error {
	prompt := []*survey.Question{
		{
			Name: name,
			Prompt: &survey.Input{
				Message: question,
				Default: defaultVal,
			},
		},
	}
	var err error
	if isRequired {
		err = Ask(prompt, response, survey.WithValidator(survey.Required))
	} else {
		err = Ask(prompt, response)
	}
	if err != nil {
		return err
	}
	return nil
}

func MultiSelect(response interface{}, name, question string, options []string, opts ...survey.AskOpt) error {
	prompt := []*survey.Question{
		{
			Name: name,
			Prompt: &survey.MultiSelect{
				Message: question,
				Options: options,
			},
		},
	}
	err := Ask(prompt, response, opts...)
	if err != nil {
		return err
	}
	return nil
}

func AskMultiline(response interface{}, name, question string, defaultVal string) error {
	prompt := []*survey.Question{
		{
			Name: name,
			Prompt: &survey.Multiline{
				Message: question,
				Default: defaultVal,
			},
		},
	}
	err := Ask(prompt, response)
	if err != nil {
		return err
	}
	return nil
}

func Select(response interface{}, name string, question string, options []string, opts ...survey.AskOpt) error {
	prompt := []*survey.Question{
		{
			Name: name,
			Prompt: &survey.Select{
				Message: question,
				Options: options,
			},
		},
	}
	err := Ask(prompt, response, opts...)
	if err != nil {
		return err
	}
	return nil
}

var AskOne = survey.AskOne

var Ask = survey.Ask

var Confirm = func(result *bool, prompt string, defaultVal bool) error {
	p := &survey.Confirm{
		Message: prompt,
		Default: defaultVal,
	}
	return survey.AskOne(p, result)
}

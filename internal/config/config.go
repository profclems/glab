package config

import (
	"github.com/profclems/glab/internal/utils"
)

// PromptAndSetEnv : prompts user for value and returns default value if empty
func Prompt(question, defaultVal string) (envVal string, err error) {
	envVal = utils.AskQuestionWithInput(question, defaultVal, false)
	return
}

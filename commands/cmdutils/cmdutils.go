package cmdutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/profclems/glab/internal/config"

	"github.com/profclems/glab/internal/git"
)

const (
	IssueTemplate        = "issue_templates"
	MergeRequestTemplate = "merge_request_templates"
)

// LoadGitLabTemplate finds and loads the GitLab template from the working git directory
// Follows the format officially supported by GitLab
// https://docs.gitlab.com/ee/user/project/description_templates.html#setting-a-default-template-for-issues-and-merge-requests.
//
// TODO: load from remote repository if repo is overriden by -R flag
func LoadGitLabTemplate(tmplType, tmplName string) (string, error) {
	wdir, err := git.ToplevelDir()
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(tmplName, ".md") {
		tmplName = tmplName + ".md"
	}

	tmplFile := filepath.Join(wdir, ".gitlab", tmplType, tmplName)
	f, err := os.Open(tmplFile)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	tmpl, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(tmpl)), nil
}

func ListGitLabTemplates(tmplType string) ([]string, error) {
	wdir, err := git.ToplevelDir()
	tmplFolder := filepath.Join(wdir, ".gitlab", tmplType)
	var files []string
	f, err := os.Open(tmplFolder)
	if err != nil {
		return files, err
	}
	fileNames, err := f.Readdirnames(-1)
	defer f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileNames {
		files = append(files, strings.TrimSuffix(file, ".md"))
	}
	return files, nil
}

func GetEditor(cf func() (config.Config, error)) (string, error) {
	cfg, err := cf()
	if err != nil {
		return "", fmt.Errorf("could not read config: %w", err)
	}
	// will search in the order glab_editor, visual, editor first from the env before the config file
	editorCommand, _ := cfg.Get("", "editor")

	return editorCommand, nil
}

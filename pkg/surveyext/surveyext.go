// this is a wrapper for https://github.com/AlecAivazis/survey package but with
// additional extensions and customizations for glab
// adapted from https://github.com/cli/cli/blob/trunk/pkg/surveyext/editor.go
package surveyext

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/kballard/go-shellquote"
	"github.com/profclems/glab/pkg/execext"
)

var (
	bom           = []byte{0xef, 0xbb, 0xbf}
	defaultEditor = "nano" // EXTENDED to switch from vim as a default editor
	editorSkipKey = 's'
	editorOpenKey = 'e'
)

type showable interface {
	Show()
}

func init() {
	if runtime.GOOS == "windows" {
		defaultEditor = "notepad"
	} else if g := os.Getenv("GIT_EDITOR"); g != "" {
		defaultEditor = g
	} else if v := os.Getenv("VISUAL"); v != "" {
		defaultEditor = v
	} else if e := os.Getenv("EDITOR"); e != "" {
		defaultEditor = e
	}
}

// EXTENDED to enable different prompting behavior
type GLabEditor struct {
	*survey.Editor
	EditorCommand string
	BlankAllowed  bool
}

func (e *GLabEditor) editorCommand() string {
	if e.EditorCommand == "" {
		return defaultEditor
	}

	return e.EditorCommand
}

// EXTENDED to change prompt text
var EditorQuestionTemplate = `
{{- if .ShowHelp }}{{- color .Config.Icons.Help.Format }}{{ .Config.Icons.Help.Text }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color .Config.Icons.Question.Format }}{{ .Config.Icons.Question.Text }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .ShowAnswer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ .Config.HelpInput }} for help]{{color "reset"}} {{end}}
  {{- if and .Default (not .HideDefault)}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
	{{- color "cyan"}}[(e) or Enter to launch {{ .EditorCommand }}{{- if .BlankAllowed }}, (s) or Esc to skip{{ end }}] {{color "reset"}}
{{- end}}`

// EXTENDED to pass editor name (to use in prompt)
type EditorTemplateData struct {
	survey.Editor
	EditorCommand string
	BlankAllowed  bool
	Answer        string
	ShowAnswer    bool
	ShowHelp      bool
	Config        *survey.PromptConfig
}

// EXTENDED to augment prompt text and keypress handling
func (e *GLabEditor) prompt(initialValue string, config *survey.PromptConfig) (interface{}, error) {
	err := e.Render(
		EditorQuestionTemplate,
		// EXTENDED to support printing editor in prompt and BlankAllowed
		EditorTemplateData{
			Editor:        *e.Editor,
			BlankAllowed:  e.BlankAllowed,
			EditorCommand: filepath.Base(e.editorCommand()),
			Config:        config,
		},
	)
	if err != nil {
		return "", err
	}

	// start reading runes from the standard in
	rr := e.NewRuneReader()
	_ = rr.SetTermMode()
	defer func() { _ = rr.RestoreTermMode() }()

	cursor := e.NewCursor()
	cursor.Hide()
	defer cursor.Show()

	for {
		// EXTENDED to handle the Enter or e to edit / s or Esc to skip behavior + BlankAllowed
		r, _, err := rr.ReadRune()
		if err != nil {
			return "", err
		}
		if r == editorOpenKey || r == '\r' || r == '\n' {
			break
		}
		if r == editorSkipKey || r == terminal.KeyEscape {
			if e.BlankAllowed {
				return "", nil
			} else {
				continue
			}
		}
		if r == terminal.KeyInterrupt {
			return "", terminal.InterruptErr
		}
		if r == terminal.KeyEndTransmission {
			break
		}
		if string(r) == config.HelpInput && e.Help != "" {
			err = e.Render(
				EditorQuestionTemplate,
				EditorTemplateData{
					// EXTENDED to support printing editor in prompt, BlankAllowed
					Editor:        *e.Editor,
					BlankAllowed:  e.BlankAllowed,
					EditorCommand: filepath.Base(e.editorCommand()),
					ShowHelp:      true,
					Config:        config,
				},
			)
			if err != nil {
				return "", err
			}
		}
		continue
	}

	stdio := e.Stdio()
	text, err := Edit(e.editorCommand(), e.FileName, initialValue, stdio.In, stdio.Out, stdio.Err, cursor)
	if err != nil {
		return "", err
	}

	// check length, return default value on empty
	if text == "" && !e.AppendDefault {
		return e.Default, nil
	}

	return text, nil
}

// EXTENDED This is straight copypasta from survey to get our overridden prompt called.;
func (e *GLabEditor) Prompt(config *survey.PromptConfig) (interface{}, error) {
	initialValue := ""
	if e.Default != "" && e.AppendDefault {
		initialValue = e.Default
	}
	return e.prompt(initialValue, config)
}

func Edit(editorCommand, fn, initialValue string, stdin io.Reader, stdout io.Writer, stderr io.Writer, cursor showable) (string, error) {
	// prepare the temp file
	pattern := fn
	if pattern == "" {
		pattern = "survey*.txt"
	}
	f, err := ioutil.TempFile("", pattern)
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())

	// write utf8 BOM header
	// The reason why we do this is because notepad.exe on Windows determines the
	// encoding of an "empty" text file by the locale, for example, GBK in China,
	// while golang string only handles utf8 well. However, a text file with utf8
	// BOM header is not considered "empty" on Windows, and the encoding will then
	// be determined utf8 by notepad.exe, instead of GBK or other encodings.
	if _, err := f.Write(bom); err != nil {
		return "", err
	}

	// write initial value
	if _, err := f.WriteString(initialValue); err != nil {
		return "", err
	}

	// close the fd to prevent the editor unable to save file
	if err := f.Close(); err != nil {
		return "", err
	}

	if editorCommand == "" {
		editorCommand = defaultEditor
	}
	args, err := shellquote.Split(editorCommand)
	if err != nil {
		return "", err
	}
	args = append(args, f.Name())

	editorExe, err := execext.LookPath(args[0])
	if err != nil {
		return "", err
	}

	cmd := exec.Command(editorExe, args[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if cursor != nil {
		cursor.Show()
	}

	// open the editor
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// raw is a BOM-unstripped UTF8 byte slice
	raw, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return "", err
	}

	// strip BOM header
	return string(bytes.TrimPrefix(raw, bom)), nil
}

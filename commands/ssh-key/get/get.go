package get

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/iostreams"
	"github.com/profclems/glab/pkg/prompt"
	"github.com/profclems/glab/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type GetOpts struct {
	HTTPClient func() (*gitlab.Client, error)
	IO         *iostreams.IOStreams
	BaseRepo   func() (glrepo.Interface, error)

	KeyID int
}

func NewCmdGet(f *cmdutils.Factory, runE func(*GetOpts) error) *cobra.Command {
	opts := &GetOpts{
		IO: f.IO,
	}
	cmd := &cobra.Command{
		Use:   "get <key-id>",
		Short: "Gets a single key",
		Long:  "Returns a single SSH key specified by the ID",
		Example: heredoc.Doc(`
		# Get ssh key with ID as argument
		$ glab ssh-key get 7750633

		# Interactive
		$ glab ssh-key get
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.HTTPClient = f.HttpClient
			opts.BaseRepo = f.BaseRepo

			if len(args) == 0 && !opts.IO.PromptEnabled() {
				return cmdutils.FlagError{Err: errors.New("<key-id> argument is required when running in non-ttys")}
			}

			if len(args) == 1 {
				opts.KeyID = utils.StringToInt(args[0])
			}

			if runE != nil {
				return runE(opts)
			}

			return getRun(opts)
		},
	}

	return cmd
}

func getRun(opts *GetOpts) error {
	httpClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}

	if opts.KeyID == 0 {
		opts.KeyID, err = keySelectPrompt(httpClient)
		if err != nil {
			return cmdutils.WrapError(err, "failed to prompt")
		}
	}

	key, _, err := httpClient.Users.GetSSHKey(opts.KeyID)
	if err != nil {
		return cmdutils.WrapError(err, "failed to get ssh key")
	}

	opts.IO.LogInfo(key.Key)

	return nil
}

func keySelectPrompt(client *gitlab.Client) (int, error) {
	keys, _, err := client.Users.ListSSHKeys()
	if err != nil {
		return 0, err
	}

	keyOpts := map[string]int{}
	surveyOpts := make([]string, 0, len(keys))
	for _, key := range keys {
		keyOpts[key.Title] = key.ID
		surveyOpts = append(surveyOpts, key.Title)
	}

	keySelectQuestion := &survey.Select{
		Message: "Select Key",
		Options: surveyOpts,
	}

	var result string
	err = prompt.AskOne(keySelectQuestion, &result)
	return keyOpts[result], err
}

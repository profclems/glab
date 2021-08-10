package add

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/iostreams"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type AddOpts struct {
	HTTPClient func() (*gitlab.Client, error)
	IO         *iostreams.IOStreams
	BaseRepo   func() (glrepo.Interface, error)

	Title     string
	Key       string
	ExpiresAt string

	KeyFile string
}

func NewCmdAdd(f *cmdutils.Factory, runE func(*AddOpts) error) *cobra.Command {
	opts := &AddOpts{
		IO: f.IO,
	}
	cmd := &cobra.Command{
		Use:   "add <title> [key-file]",
		Short: "Add an SSH key to your GitLab account",
		Long: heredoc.Doc(`
		Creates a new SSH key owned by the currently authenticated user.

		The --title flag is always required
		`),
		Example: heredoc.Doc(`
		# Read ssh key from stdin and upload
		$ glab ssh-key add -t "my title"

		# Read ssh key from specified key file and upload
		$ glab ssh-key add ~/.ssh/id_ed25519.pub -t "my title"
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.HTTPClient = f.HttpClient
			opts.BaseRepo = f.BaseRepo

			if len(args) == 0 {
				if opts.IO.IsOutputTTY() && opts.IO.IsInTTY {
					return &cmdutils.FlagError{Err: errors.New("missing key file")}
				}
				opts.KeyFile = "-"
			} else {
				opts.KeyFile = args[0]
			}

			if runE != nil {
				return runE(opts)
			}

			return addRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "New SSH key's title")
	cmd.Flags().StringVarP(&opts.ExpiresAt, "expires-at", "e", "", "The expiration date of the SSH key in ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)")

	_ = cmd.MarkFlagRequired("title")

	return cmd
}

func addRun(opts *AddOpts) error {
	httpClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}

	var keyFileReader io.Reader
	if opts.KeyFile == "-" {
		keyFileReader = opts.IO.In
		defer opts.IO.In.Close()
	} else {
		f, err := os.Open(opts.KeyFile)
		if err != nil {
			return err
		}
		defer f.Close()

		keyFileReader = f
	}

	keyInBytes, err := ioutil.ReadAll(keyFileReader)
	if err != nil {
		return cmdutils.WrapError(err, "failed to read ssh key file")
	}

	opts.Key = string(keyInBytes)

	err = UploadSSHKey(httpClient, opts.Title, opts.Key, opts.ExpiresAt)
	if err != nil {
		return cmdutils.WrapError(err, "failed to add new ssh public key")
	}

	if opts.IO.IsOutputTTY() {
		cs := opts.IO.Color()
		opts.IO.Logf("%s New SSH public key added to your account\n", cs.GreenCheck())
	}

	return nil
}

// adapted from https://github.com/cli/cli/blob/trunk/pkg/cmd/pr/diff/diff.go
package diff

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"syscall"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/xanzy/go-gitlab"

	"github.com/spf13/cobra"
)

type DiffOptions struct {
	factory *cmdutils.Factory
	IO      *utils.IOStreams

	Args     []string
	UseColor string
}

func NewCmdDiff(f *cmdutils.Factory, runF func(*DiffOptions) error) *cobra.Command {
	opts := &DiffOptions{
		factory: f,
		IO:      f.IO,
	}

	cmd := &cobra.Command{
		Use:   "diff [<id> | <branch>]",
		Short: "View changes in a merge request",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			if repoOverride, _ := cmd.Flags().GetString("repo"); repoOverride != "" && len(args) == 0 {
				return &cmdutils.FlagError{Err: errors.New("argument required when using the --repo flag")}
			}

			if len(args) > 0 {
				opts.Args = args
			}

			if !validColorFlag(opts.UseColor) {
				return &cmdutils.FlagError{Err: fmt.Errorf("did not understand color: %q. Expected one of always, never, or auto", opts.UseColor)}
			}

			if opts.UseColor == "auto" && !opts.IO.IsaTTY {
				opts.UseColor = "never"
			}

			if runF != nil {
				return runF(opts)
			}
			return diffRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.UseColor, "color", "auto", "Use color in diff output: {always|never|auto}")

	return cmd
}

func diffRun(opts *DiffOptions) error {
	apiClient, err := opts.factory.HttpClient()
	if err != nil {
		return err
	}
	mr, baseRepo, err := mrutils.MRFromArgs(opts.factory, opts.Args)
	if err != nil {
		return err
	}

	diffs, _, err := apiClient.MergeRequests.GetMergeRequestDiffVersions(baseRepo.FullName(), mr.IID, &gitlab.GetMergeRequestDiffVersionsOptions{})
	if err != nil {
		return fmt.Errorf("could not find merge request diffs: %w", err)
	}

	diffOut := &bytes.Buffer{}
	for _, diff := range diffs {
		// the diffs are not included in the GetMergeRequestDiffVersions so we query for each diff version
		diffVersion, _, err := apiClient.MergeRequests.GetSingleMergeRequestDiffVersion(baseRepo.FullName(), mr.IID, diff.ID)
		if err != nil {
			return fmt.Errorf("could not find merge request diff: %w", err)
		}
		for _, diffLine := range diffVersion.Diffs {
			if diffLine.RenamedFile {
				diffOut.WriteString("-" + diffLine.OldPath + "\n")
			}
			if diffLine.NewFile || diffLine.RenamedFile {
				diffOut.WriteString("+" + diffLine.NewPath + "\n")
			} else {
				diffOut.WriteString(diffLine.OldPath + "\n")
			}

			diffOut.WriteString(diffLine.Diff)
		}
	}

	defer diffOut.Reset()

	err = opts.IO.StartPager()
	if err != nil {
		return err
	}
	defer opts.IO.StopPager()

	if opts.UseColor == "never" {
		_, err = io.Copy(opts.IO.StdOut, diffOut)
		if errors.Is(err, syscall.EPIPE) {
			return nil
		}
		return err
	}

	diffLines := bufio.NewScanner(diffOut)
	for diffLines.Scan() {
		diffLine := diffLines.Text()
		switch {
		case isHeaderLine(diffLine):
			fmt.Fprintf(opts.IO.StdOut, "\x1b[1;38m%s\x1b[m\n", diffLine)
		case isAdditionLine(diffLine):
			fmt.Fprintf(opts.IO.StdOut, "\x1b[32m%s\x1b[m\n", diffLine)
		case isRemovalLine(diffLine):
			fmt.Fprintf(opts.IO.StdOut, "\x1b[31m%s\x1b[m\n", diffLine)
		default:
			fmt.Fprintln(opts.IO.StdOut, diffLine)
		}
	}

	if err := diffLines.Err(); err != nil {
		return fmt.Errorf("error reading merge request diff: %w", err)
	}

	return nil
}

var diffHeaderPrefixes = []string{"+++", "---", "diff", "index"}

func isHeaderLine(dl string) bool {
	for _, p := range diffHeaderPrefixes {
		if strings.HasPrefix(dl, p) {
			return true
		}
	}
	return false
}

func isAdditionLine(dl string) bool {
	return strings.HasPrefix(dl, "+")
}

func isRemovalLine(dl string) bool {
	return strings.HasPrefix(dl, "-")
}

func validColorFlag(c string) bool {
	return c == "auto" || c == "always" || c == "never"
}

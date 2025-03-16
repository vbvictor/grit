package stat

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/git"
)

var churnOpts = &git.ChurnOptions{
	SortBy:       git.Changes,
	Top:          git.DefaultTop,
	Extensions:   nil,
	Since:        time.Time{},
	Until:        time.Time{},
	Path:         "",
	ExcludeRegex: nil,
	OutputFormat: "",
}

var (
	extensionList     []string
	since             string
	until             string
	repoPath          string
	excludeChurnRegex string
)

var ChurnCmd = &cobra.Command{ //nolint:exhaustruct // no need to set all fields
	Use:   "churn [flags] <repository>",
	Short: "Finds files with the most changes in git repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		var err error
		repoPath, err = filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("error getting absolute path: %w", err)
		}

		flag.LogIfVerbose("Processing repository: %s\n", repoPath)

		if err := git.PopulateOpts(churnOpts, extensionList, since, until, repoPath, excludeChurnRegex); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}

		churns, err := git.ReadGitChurn(repoPath, churnOpts)
		if err != nil {
			return fmt.Errorf("error getting churn metrics: %w", err)
		}

		churns = git.SortAndLimit(churns, churnOpts.SortBy, churnOpts.Top)

		return printChurnStats(churns, os.Stdout, churnOpts)
	},
}

func init() {
	flags := ChurnCmd.PersistentFlags()

	flag.SortFlag(flags, &churnOpts.SortBy, git.Commits,
		fmt.Sprintf("Specify churn sort type: [%s, %s, %s, %s]", git.Changes, git.Additions, git.Deletions, git.Commits))
	flag.TopFlag(flags, &churnOpts.Top)
	flag.VerboseFlag(flags, &flag.Verbose)
	flag.OutputFormatFlag(flags, &churnOpts.OutputFormat)
	flag.ExcludeRegexFlag(flags, &excludeChurnRegex)
	flag.ExtensionsFlag(flags, &extensionList)
	flag.SinceFlag(flags, &since)
	flag.UntilFlag(flags, &until)

	ChurnCmd.Flag(flag.LongUntil).DefValue = flag.DefaultUntil
	ChurnCmd.Flag(flag.LongSince).DefValue = flag.DefaultSince
}

func printChurnStats(results []*git.ChurnChunk, out io.Writer, opts *git.ChurnOptions) error {
	switch opts.OutputFormat {
	case flag.CSV:
		git.PrintCSV(results, out, opts)
	case flag.Tabular:
		git.PrintTable(results, out, opts)
	default:
		return fmt.Errorf("unsupported output format: %s", opts.OutputFormat)
	}

	return nil
}

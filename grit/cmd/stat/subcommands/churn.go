package stat

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/git"
)

var churnOpts = &git.ChurnOptions{
	SortBy:      git.Changes,
	Top:         git.DefaultTop,
	Extensions:  nil,
	Since:       time.Time{},
	Until:       time.Time{},
	Path:        "",
	ExcludePath: "",
}

var (
	extensionList []string
	since         string
	until         string
	repoPath      string
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

		if err := git.PopulateOpts(churnOpts, extensionList, since, until, repoPath); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}

		churns, err := git.ReadGitChurn(repoPath, churnOpts)
		if err != nil {
			return fmt.Errorf("error getting churn metrics: %w", err)
		}

		churns = git.SortAndLimit(churns, churnOpts.SortBy, churnOpts.Top)

		return git.PrintStats(churns, os.Stdout, churnOpts)
	},
}

func init() {
	flags := ChurnCmd.PersistentFlags()

	flags.StringVar(&churnOpts.SortBy, flag.LongSort, git.Commits,
		fmt.Sprintf("Specify churn sort type: [%s, %s, %s, %s]", git.Changes, git.Additions, git.Deletions, git.Commits))
	flags.IntVarP(&churnOpts.Top, flag.LongTop, flag.ShortTop, git.DefaultTop, "Number of top files to display")
	flags.BoolVarP(&flag.Verbose, flag.LongVerbose, flag.ShortVerbose, false, "Show detailed progress")
	flags.StringVar(&churnOpts.ExcludePath, flag.LongExclude, "", "Exclude files matching regex pattern")
	flags.StringSliceVarP(&extensionList, flag.LongExtensions, flag.ShortExt, nil,
		"Only include files with given extensions in comma-separated list, e.g. 'go,h,c'")
	flags.StringVarP(&since, flag.LongSince, flag.ShortSince, "", "Start date for analysis in format 'YYYY-MM-DD'")
	flags.StringVarP(&until, flag.LongUntil, flag.ShortUntil, "", "End date for analysis in format 'YYYY-MM-DD'")

	ChurnCmd.Flag(flag.LongUntil).DefValue = flag.DefaultUntil
	ChurnCmd.Flag(flag.LongSince).DefValue = flag.DefaultSince
}

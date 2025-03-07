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
	SortBy:       git.Changes,
	Top:          git.DefaultTop,
	Extensions:   nil,
	Since:        time.Time{},
	Until:        time.Time{},
	OutputFormat: flag.JSON,
	Path:         "",
	ExcludePath:  "",
}

var (
	extensionList []string
	since         string
	until         string
	repoPath      string
)

var ChurnCmd = &cobra.Command{ //nolint:exhaustruct // no need to set all fields
	Use:   "churn <path>",
	Short: "Finds files with the most changes in git repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		var err error
		repoPath, err = filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("error getting absolute path: %w", err)
		}

		flag.LogIfVerbose("Processing repository: %s\n", repoPath)

		if err := populateOptsFromFlags(); err != nil {
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

	flags.StringVar(&churnOpts.SortBy, flag.LongSort, "commits",
		fmt.Sprintf("Sort by: %s, %s, %s, %s", git.Changes, git.Additions, git.Deletions, git.Commits))
	flags.IntVarP(&churnOpts.Top, flag.LongTop, flag.ShortTop, git.DefaultTop, "Number of top files to display")
	flags.BoolVarP(&flag.Verbose, flag.LongVerbose, flag.ShortVerbose, false, "Show detailed progress")
	flags.StringVar(&churnOpts.ExcludePath, flag.LongExclude, "", "Exclude files matching regex pattern")
	flags.StringSliceVarP(&extensionList, flag.LongExtensions, flag.ShortExt, nil,
		"Only include files with given extensions in comma-separated list. For example go,h,c")
	flags.StringVarP(&since, flag.LongSince, flag.ShortSince, "", "Start date for analysis (YYYY-MM-DD)")
	flags.StringVarP(&until, flag.LongUntil, flag.ShortUntil, "", "End date for analysis (YYYY-MM-DD)")

	ChurnCmd.Flag(flag.LongUntil).DefValue = flag.DefaultUntil
	ChurnCmd.Flag(flag.LongSince).DefValue = flag.DefaultSince
}

func populateOptsFromFlags() error {
	churnOpts.Path = repoPath

	if err := setSinceOpt(churnOpts, since); err != nil {
		return fmt.Errorf("error setting since option: %w", err)
	}

	if err := setUntilOpt(churnOpts, until); err != nil {
		return fmt.Errorf("error setting until option: %w", err)
	}

	if flag.Extensions != nil {
		churnOpts.Extensions = flag.GetExtMap(flag.Extensions)
	}

	return nil
}

func setSinceOpt(opts *git.ChurnOptions, since string) error {
	if since != "" {
		var err error

		opts.Since, err = time.Parse(time.DateOnly, since)
		if err != nil {
			return fmt.Errorf("error parsing since date: %w", err)
		}
	} else {
		opts.Since = time.Now().AddDate(-1, 0, 0)
	}

	return nil
}

func setUntilOpt(opts *git.ChurnOptions, until string) error {
	if until != "" {
		var err error

		opts.Until, err = time.Parse(time.DateOnly, until)
		if err != nil {
			return fmt.Errorf("error parsing until date: %w", err)
		}
	} else {
		opts.Until = time.Now()
	}

	return nil
}

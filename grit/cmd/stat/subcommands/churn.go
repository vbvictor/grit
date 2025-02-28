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

var ChurnCmd = &cobra.Command{
	Use:   "churn <path>",
	Short: "Find files with the most changes in git repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		repoPath, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("error getting absolute path: %w", err)
		}

		if flag.Verbose {
			fmt.Printf("Processing repository: %s\n", repoPath)
		}

		opts, err := ChurnOptsFromFlags()
		if err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}
		opts.Path = repoPath

		churns, err := git.ReadGitChurn(repoPath, opts)
		if err != nil {
			return fmt.Errorf("error getting churn metrics: %w", err)
		}

		churns = git.SortAndLimit(churns, opts.SortBy, opts.Top)

		return git.PrintStats(churns, os.Stdout, opts)
	},
}

func init() {
	flags := ChurnCmd.PersistentFlags()

	flags.StringVar(&flag.SortBy, flag.LongSort, "commits",
		fmt.Sprintf("Sort by: %s, %s, %s, %s", git.Changes, git.Additions, git.Deletions, git.Commits))
	flags.IntVarP(&flag.Top, flag.LongTop, flag.ShortTop, git.DefaultTop, "Number of top files to display")
	flags.BoolVarP(&flag.Verbose, flag.LongVerbose, flag.ShortVerbose, false, "Show detailed progress")
	flags.StringVar(&flag.ExcludePath, flag.LongExclude, "", "Exclude files matching regex pattern")
	flags.StringSliceVarP(&flag.Extensions, flag.LongExtensions, flag.ShortExt, nil,
		"Only include files with given extensions in comma-separated list. For example go,h,c")
	flags.StringVarP(&flag.Since, flag.LongSince, flag.ShortSince, "", "Start date for analysis (YYYY-MM-DD)")
	flags.StringVarP(&flag.Until, flag.LongUntil, flag.ShortUntil, "", "End date for analysis (YYYY-MM-DD)")

	ChurnCmd.Flag(flag.LongUntil).DefValue = flag.DefaultUntil
	ChurnCmd.Flag(flag.LongSince).DefValue = flag.DefaultSince
}

func ChurnOptsFromFlags() (*git.ChurnOptions, error) {
	opts := new(git.ChurnOptions)

	opts.SortBy = flag.SortBy
	opts.Top = flag.Top
	opts.ExcludePath = flag.ExcludePath
	opts.OutputFormat = flag.OutputFormat

	if flag.Since != "" {
		var err error

		opts.Since, err = time.Parse(time.DateOnly, flag.Since)
		if err != nil {
			return opts, fmt.Errorf("error parsing since date: %w", err)
		}
	} else {
		opts.Since = time.Now().AddDate(-1, 0, 0)
	}

	if flag.Until != "" {
		var err error

		opts.Until, err = time.Parse(time.DateOnly, flag.Until)
		if err != nil {
			return opts, fmt.Errorf("error parsing since date: %w", err)
		}
	} else {
		opts.Until = time.Now()
	}

	if flag.Extensions != nil {
		opts.Extensions = flag.GetExtMap(flag.Extensions)
	}

	return opts, nil
}

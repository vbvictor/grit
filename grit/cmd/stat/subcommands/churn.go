package stat

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/git"
)

var ChurnCmd = &cobra.Command{
	Use:   "churn <repository>",
	Short: "Get churn metrics of a repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		repoPath, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("error getting absolute path: %w", err)
		}

		if flag.Verbose {
			fmt.Printf("Processing repository: %s\n", repoPath)
		}

		opts, err := createOptsFromFlags()
		if err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}
		opts.Path = repoPath

		return git.PrintRepoStats(repoPath, opts)
	},
}

func init() {
	flags := ChurnCmd.PersistentFlags()

	flags.StringVar(&flag.SortBy, "sort", "commits",
		fmt.Sprintf("Sort by: %s, %s, %s, %s", git.Changes, git.Additions, git.Deletions, git.Commits))
	flags.IntVarP(&flag.Top, "top", "t", git.DefaultTop, "Number of top files to display")
	flags.BoolVarP(&flag.Verbose, "verbose", "v", false, "Show detailed progress")
	flags.StringVar(&flag.ExcludePath, "exclude", "", "Exclude files matching regex pattern")
	flags.StringSliceVarP(&flag.Extensions, "ext", "e", nil,
		"Only include files with given extensions in comma-separated list. For example go,h,c")
	flags.StringVarP(&flag.Since, "since", "s", "", "Start date for analysis (YYYY-MM-DD)")
	flags.StringVarP(&flag.Until, "until", "u", "", "End date for analysis (YYYY-MM-DD)")
	flags.StringVarP(&flag.OutputFormat, "format", "f", flag.Tabular,
		fmt.Sprintf("Output format [%s]", strings.Join(flag.AvailableOutputFormats, ", ")))

	ChurnCmd.Flag("until").DefValue = "current date"
	ChurnCmd.Flag("since").DefValue = "one year ago"
}

func createOptsFromFlags() (git.ChurnOptions, error) {
	opts := git.ChurnOptions{}

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

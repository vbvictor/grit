package stat

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/complexity"
	"github.com/vbvictor/grit/pkg/git"
)

var ComplexityCmd = &cobra.Command{
	Use:   "complexity <path>",
	Short: "Show files that have the most complexity in repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		repoPath, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("error getting absolute path: %w", err)
		}

		if flag.Verbose {
			fmt.Printf("Processing repository: %s\n", repoPath)
		}

		opts, err := complexityOptsFromFlags()
		if err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}

		fileStat, err := complexity.RunComplexity(repoPath, opts)
		if err != nil {
			return fmt.Errorf("error running complexity analysis: %w", err)
		}

		fileStat = sortAndLimit(fileStat, opts)

		complexity.PrintTabular(fileStat, os.Stdout)

		return nil
	},
}

func init() {
	flags := ComplexityCmd.PersistentFlags()

	flags.StringVarP(&flag.Engine, flag.LongEngine, flag.ShortEngine, complexity.Gocyclo,
		"Complexity calculation engine to use: gocyclo or gocognit")
	flags.IntVarP(&flag.Top, flag.LongTop, flag.ShortTop, git.DefaultTop, "Number of top files to display")
	flags.BoolVarP(&flag.Verbose, flag.LongVerbose, flag.ShortVerbose, false, "Show detailed progress")
	flags.StringVar(&flag.ExcludePath, flag.LongExclude, "", "Exclude files matching regex pattern")
}

func complexityOptsFromFlags() (complexity.Options, error) { //nolint:unparam // May be used in the future
	opts := complexity.Options{}

	opts.Top = flag.Top
	opts.ExcludePath = flag.ExcludePath
	opts.Engine = flag.Engine

	return opts, nil
}

func sortAndLimit(fileStat []*complexity.FileStat, opts complexity.Options) []*complexity.FileStat {
	slices.SortFunc(fileStat, func(a, b *complexity.FileStat) int {
		if b.AvgComplexity > a.AvgComplexity {
			return 1
		}

		if b.AvgComplexity < a.AvgComplexity {
			return -1
		}

		return 0
	})

	if opts.Top > 0 && opts.Top < len(fileStat) {
		return fileStat[:opts.Top]
	}

	return fileStat
}

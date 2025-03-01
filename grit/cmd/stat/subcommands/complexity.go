package stat

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/complexity"
	"github.com/vbvictor/grit/pkg/git"
)

var ComplexityCmd = &cobra.Command{
	Use:   "complexity <path>",
	Short: "Finds the most complex files",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		repoPath, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("error getting absolute path: %w", err)
		}

		if flag.Verbose {
			fmt.Printf("Processing repository: %s\n", repoPath)
		}

		opts, err := ComplexityOptsFromFlags()
		if err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}

		fileStat, err := complexity.RunComplexity(repoPath, opts)
		if err != nil {
			return fmt.Errorf("error running complexity analysis: %w", err)
		}

		fileStat = complexity.SortAndLimit(fileStat, opts)

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

func ComplexityOptsFromFlags() (complexity.Options, error) { //nolint:unparam // error return is reserved for future
	opts := complexity.Options{}

	opts.Top = flag.Top
	opts.ExcludePath = flag.ExcludePath
	opts.Engine = flag.Engine

	return opts, nil
}

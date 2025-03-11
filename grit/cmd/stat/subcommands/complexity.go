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

var complexityOpts = complexity.Options{
	Engine:       complexity.Gocyclo,
	ExcludeRegex: nil,
	Top:          10, //nolint:mnd // default value
}

var excludeComplexityRegex string

var ComplexityCmd = &cobra.Command{ //nolint:exhaustruct // no need to set all fields
	Use:   "complexity [flags] <path>",
	Short: "Finds the most complex files",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		repoPath, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("error getting absolute path: %w", err)
		}

		flag.LogIfVerbose("Processing repository: %s\n", repoPath)

		if err := complexity.PopulateOpts(&complexityOpts, excludeComplexityRegex); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}

		fileStat, err := complexity.RunComplexity(repoPath, &complexityOpts)
		if err != nil {
			return fmt.Errorf("error running complexity analysis: %w", err)
		}

		fileStat = complexity.SortAndLimit(fileStat, complexityOpts)

		complexity.PrintTabular(fileStat, os.Stdout)

		return nil
	},
}

func init() {
	flags := ComplexityCmd.PersistentFlags()

	flags.StringVarP(&complexityOpts.Engine, flag.LongEngine, flag.ShortEngine, complexity.Gocyclo,
		fmt.Sprintf("Specify complexity calculation engine: [%s, %s]", complexity.Gocyclo, complexity.Gocognit))
	flags.IntVarP(&complexityOpts.Top, flag.LongTop, flag.ShortTop, git.DefaultTop, "Number of top files to display")
	flags.BoolVarP(&flag.Verbose, flag.LongVerbose, flag.ShortVerbose, false, "Show detailed progress")
	flags.StringVar(&excludeComplexityRegex, flag.LongExclude, "", "Exclude files matching regex pattern")
}

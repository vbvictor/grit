package stat

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/complexity"
)

var complexityOpts = complexity.Options{
	Engine:       complexity.Gocyclo,
	ExcludeRegex: nil,
	Top:          10, //nolint:mnd // default value
	OutputFormat: "",
}

var excludeComplexityRegex string

var ComplexityCmd = &cobra.Command{ //nolint:exhaustruct // no need to set all fields
	Use:   "complexity [flags] <path>",
	Short: "Finds the most complex files",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		path := filepath.Clean(args[0])
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("repository does not exist: %w", err)
		}

		flag.LogIfVerbose("Processing repository: %s\n", path)

		if err := complexity.PopulateOpts(&complexityOpts, excludeComplexityRegex); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}

		fileStat, err := complexity.RunComplexity(path, &complexityOpts)
		if err != nil {
			return fmt.Errorf("error running complexity analysis: %w", err)
		}

		fileStat = complexity.SortAndLimit(fileStat, complexityOpts)

		return printComplexityStats(fileStat, os.Stdout, &complexityOpts)
	},
}

func init() {
	flags := ComplexityCmd.PersistentFlags()

	flag.ComplexityEngineFlag(flags, &complexityOpts.Engine)
	flag.TopFlag(flags, &complexityOpts.Top)
	flag.VerboseFlag(flags, &flag.Verbose)
	flag.ExcludeRegexFlag(flags, &excludeComplexityRegex)
	flag.OutputFormatFlag(flags, &complexityOpts.OutputFormat)
}

func printComplexityStats(results []*complexity.FileStat, out io.Writer, opts *complexity.Options) error {
	switch opts.OutputFormat {
	case flag.CSV:
		complexity.PrintCSV(results, out)
	case flag.Tabular:
		complexity.PrintTabular(results, out)
	default:
		return fmt.Errorf("unsupported output format: %s", opts.OutputFormat)
	}

	return nil
}

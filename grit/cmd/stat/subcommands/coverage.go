package stat

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/coverage"
)

var coverageOpts = &coverage.Options{
	SortBy:           coverage.Worst,
	Top:              10, //nolint:mnd // default value
	ExcludeRegex:     nil,
	RunCoverage:      "",
	CoverageFilename: "coverage.out",
	OutputFormat:     "",
}

var excludeCoverageRegex string

var CoverageCmd = &cobra.Command{ //nolint:exhaustruct // no need to set all fields
	Use:   "coverage [flags] <path>",
	Short: "Finds files with the least unit-test coverage",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		repoPath := filepath.ToSlash(filepath.Clean(args[0]))

		flag.LogIfVerbose("Processing directory: %s\n", repoPath)

		if err := coverage.PopulateOpts(coverageOpts, excludeCoverageRegex); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}

		covData, err := coverage.GetCoverageData(repoPath, coverageOpts)
		if err != nil {
			return fmt.Errorf("failed to get coverage data: %w", err)
		}

		covData = coverage.SortAndLimit(covData, coverageOpts.SortBy, coverageOpts.Top)

		return printCoverageStats(covData, os.Stdout, coverageOpts)
	},
}

func init() {
	flags := CoverageCmd.PersistentFlags()

	flag.SortFlag(flags, &coverageOpts.SortBy, coverage.Worst,
		fmt.Sprintf("Specify sort type: [%s, %s]", coverage.Worst, coverage.Best))
	flag.RunCoverageFlag(flags, &coverageOpts.RunCoverage)
	flag.CoverageFilenameFlag(flags, &coverageOpts.CoverageFilename)
	flag.VerboseFlag(flags, &flag.Verbose)
	flag.TopFlag(flags, &coverageOpts.Top)
	flag.ExcludeRegexFlag(flags, &excludeCoverageRegex)
	flag.OutputFormatFlag(flags, &coverageOpts.OutputFormat)
}

func printCoverageStats(results []*coverage.FileCoverage, out io.Writer, opts *coverage.Options) error {
	switch opts.OutputFormat {
	case flag.CSV:
		coverage.PrintCSV(results, out)
	case flag.Tabular:
		coverage.PrintTabular(results, out)
	default:
		return fmt.Errorf("unsupported output format: %s", opts.OutputFormat)
	}

	return nil
}

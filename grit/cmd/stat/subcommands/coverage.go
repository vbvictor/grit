package stat

import (
	"errors"
	"fmt"
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
		repoPath, err := filepath.Abs(args[0])
		if err != nil {
			return errors.Join(&flag.AbsRepoPathError{Path: args[0]}, err)
		}

		flag.LogIfVerbose("Processing directory: %s\n", repoPath)

		if err := coverage.PopulateOpts(coverageOpts, excludeCoverageRegex); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}

		covData, err := coverage.GetCoverageData(repoPath, coverageOpts)
		if err != nil {
			return fmt.Errorf("failed to get coverage data: %w", err)
		}

		covData = coverage.SortAndLimit(covData, coverageOpts.SortBy, coverageOpts.Top)

		return coverage.PrintStats(covData, os.Stdout, coverageOpts)
	},
}

func init() {
	flags := CoverageCmd.PersistentFlags()

	flags.StringVar(&coverageOpts.SortBy, flag.LongSort, coverage.Worst,
		fmt.Sprintf("Specify sort type: [%s, %s]", coverage.Worst, coverage.Best))
	flags.StringVarP(&coverageOpts.RunCoverage, flag.LongRunCoverage, flag.ShortRunCoverage, flag.Auto,
		`Specify tests run format:
  'Auto' will run unit tests if coverage file is not found
  'Always' will run unit tests on every invoke 'coverage' command
  'Never' will never run unit tests and always look for present test-coverage file
`)
	flags.StringVarP(&coverageOpts.CoverageFilename, flag.LongFileCoverage, flag.ShortFileCoverage, "coverage.out",
		"Name of code coverage file to read or create")
	flags.BoolVarP(&flag.Verbose, flag.LongVerbose, flag.ShortVerbose, false, "Enable verbose output")
	flags.IntVarP(&coverageOpts.Top, flag.LongTop, flag.ShortTop, flag.DefaultTop, "Number of top files to display")
	flags.StringVar(&excludeCoverageRegex, flag.LongExclude, "", "Exclude files matching regex pattern")
	flags.StringVarP(&coverageOpts.OutputFormat, flag.LongFormat, flag.ShortFormat, flag.Tabular,
		fmt.Sprintf("Specify output format: [%s, %s]", flag.Tabular, flag.CSV))
}

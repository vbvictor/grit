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

var coverageOpts = coverage.Options{
	SortBy:      coverage.Worst,
	Top:         10, //nolint:mnd // default value
	ExcludePath: "",
}

var (
	runCoverage  string
	coverageFile string
)

var CoverageCmd = &cobra.Command{ //nolint:exhaustruct // no need to set all fields
	Use:   "coverage <path>",
	Short: "Finds files with the least unit-test coverage",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		repoPath, err := filepath.Abs(args[0])
		if err != nil {
			return errors.Join(&flag.AbsRepoPathError{Path: args[0]}, err)
		}

		flag.LogIfVerbose("Processing directory: %s\n", repoPath)

		covData, err := GetCoverageData(repoPath, flag.CoverageFile, coverageOpts)
		if err != nil {
			return fmt.Errorf("failed to get coverage data: %w", err)
		}

		covData = coverage.SortAndLimit(covData, coverageOpts.SortBy, coverageOpts.Top)

		coverage.PrintTabular(covData, os.Stdout)

		return nil
	},
}

func init() {
	flags := CoverageCmd.PersistentFlags()

	flags.StringVar(&coverageOpts.SortBy, flag.LongSort, coverage.Worst, "Specify sort type")
	flags.StringVarP(&runCoverage, flag.LongRunCoverage, flag.ShortRunCoverage, flag.Auto, `Specify tests run format:
'Auto' will run unit tests if coverage file is not found
'Always' will run unit tests on every invoke 'stat coverage'
'Never' will never run unit tests and always look for present coverage file`)
	flags.StringVarP(&coverageFile, flag.LongFileCoverage, flag.ShortFileCoverage, "coverage.out",
		"Name of code coverage file to read or create")
	flags.BoolVarP(&flag.Verbose, flag.LongVerbose, flag.ShortVerbose, false, "Enable verbose output")
	flags.IntVarP(&coverageOpts.Top, flag.LongTop, flag.ShortTop, flag.DefaultTop, "Number of top files to display")
	flags.StringVar(&coverageOpts.ExcludePath, flag.LongExclude, "", "Exclude files matching regex")
}

func GetCoverageData(repoPath, coverageFile string, coverageOpts coverage.Options) ([]*coverage.FileCoverage, error) {
	coveragePath := filepath.Join(repoPath, coverageFile)

	_, err := os.Stat(coveragePath)
	if os.IsNotExist(err) {
		flag.LogIfVerbose("Coverage file %s not found\n", coveragePath)

		if runCoverage != flag.Never {
			flag.LogIfVerbose("Running test suite\n\n")

			if err = coverage.RunCoverage(repoPath, coverageFile); err != nil {
				return nil, errors.Join(flag.ErrRunCoverage, err)
			}

			flag.LogIfVerbose("Coverage file %s created\n", coveragePath)
		} else {
			flag.LogIfVerbose("Go test did not run since flag run='Never'\n")

			return nil, errors.Join(flag.ErrCoverageNotFound, err)
		}
	} else if runCoverage == flag.Always {
		flag.LogIfVerbose("Removing previous test coverage file %s\n", coveragePath)
		os.Remove(coveragePath)
		flag.LogIfVerbose("Running test suite\n\n")

		if err = coverage.RunCoverage(repoPath, coverageFile); err != nil {
			return nil, errors.Join(flag.ErrRunCoverage, err)
		}

		flag.LogIfVerbose("Coverage file %s created\n", coveragePath)
	}

	covData, err := coverage.ReadCoverage(filepath.Join(repoPath, flag.CoverageFile), coverageOpts)
	if err != nil {
		return nil, errors.Join(flag.ErrReadCoverage, err)
	}

	return covData, nil
}

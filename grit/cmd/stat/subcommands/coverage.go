package stat

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/coverage"
	"github.com/vbvictor/grit/pkg/process"
)

var CoverageCmd = &cobra.Command{
	Use:   "coverage <path>",
	Short: "Show files unit-test coverage in repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		repoPath, err := filepath.Abs(args[0])
		if err != nil {
			return errors.Join(&process.ErrAbsRepoPath{Path: args[0]}, err)
		}

		flag.LogIfVerbose("Processing directory: %s\n", repoPath)

		opts := CoverageOptsFromFlags()

		covData, err := GetCoverageData(repoPath, flag.CoverageFile, opts)
		if err != nil {
			return fmt.Errorf("failed to get coverage data: %w", err)
		}

		covData = coverage.SortAndLimit(covData, opts.SortBy, opts.Top)

		coverage.PrintTabular(covData, os.Stdout)

		return nil
	},
}

func init() {
	flags := CoverageCmd.PersistentFlags()

	flags.StringVar(&flag.SortBy, flag.LongSort, coverage.Worst, "Specify sort type")
	flags.StringVarP(&flag.RunCoverage, flag.LongRunCoverage, flag.ShortRunCoverage, flag.Auto, `Specify tests run format:
'Auto' will run tests if coverage file is not found
'Always' will run tests on every invoke 'stat coverage'
'Never' will never run tests and always look for present coverage file`)
	flags.StringVarP(&flag.CoverageFile, flag.LongFileCoverage, flag.ShortFileCoverage, "coverage.out",
		"Name of code coverage file to read or create")
	flags.BoolVarP(&flag.Verbose, flag.LongVerbose, flag.ShortVerbose, false, "Enable verbose output")
	flags.IntVarP(&flag.Top, flag.LongTop, flag.ShortTop, flag.DefaultTop, "Number of top files to display")
	flags.StringVar(&flag.ExcludePath, flag.LongExclude, "", "Exclude files matching regex")
}

func CoverageOptsFromFlags() coverage.Options {
	opts := coverage.Options{
		SortBy:      flag.SortBy,
		Top:         flag.Top,
		ExcludePath: flag.ExcludePath,
	}

	return opts
}

func GetCoverageData(repoPath, coverageFile string, opts coverage.Options) ([]*coverage.FileCoverage, error) {
	coveragePath := filepath.Join(repoPath, coverageFile)

	_, err := os.Stat(coveragePath)
	if os.IsNotExist(err) {
		flag.LogIfVerbose("Coverage file %s not found\n", coveragePath)

		if flag.RunCoverage != flag.Never {
			flag.LogIfVerbose("Running test suite\n\n")

			if err = coverage.RunCoverage(repoPath, coverageFile); err != nil {
				return nil, errors.Join(process.ErrRunCoverage, err)
			}

			flag.LogIfVerbose("Coverage file %s created\n", coveragePath)
		} else {
			flag.LogIfVerbose("Go test did not run since flag run='Never'\n")

			return nil, errors.Join(process.ErrCoverageNotFound, err)
		}
	} else if flag.RunCoverage == flag.Always {
		flag.LogIfVerbose("Removing previous test coverage file %s\n", coveragePath)
		os.Remove(coveragePath)
		flag.LogIfVerbose("Running test suite\n\n")

		if err = coverage.RunCoverage(repoPath, coverageFile); err != nil {
			return nil, errors.Join(process.ErrRunCoverage, err)
		}

		flag.LogIfVerbose("Coverage file %s created\n", coveragePath)
	}

	covData, err := coverage.ReadCoverage(filepath.Join(repoPath, flag.CoverageFile), opts)
	if err != nil {
		return nil, errors.Join(process.ErrReadCoverage, err)
	}

	return covData, nil
}

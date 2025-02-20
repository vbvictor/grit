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

		if flag.Verbose {
			fmt.Printf("Processing repository: %s\n", repoPath)
		}

		coveragePath := filepath.Join(repoPath, flag.CoverageFile)

		_, err = os.Stat(coveragePath)
		if os.IsNotExist(err) {
			if flag.RunCoverage {
				if err = coverage.RunCoverage(repoPath, flag.CoverageFile); err != nil {
					return errors.Join(process.ErrRunCoverage, err)
				}
			} else {
				return errors.Join(process.ErrCoverageNotFound, err)
			}
		}

		opts := coverageOptsFromFlags()

		covData, err := coverage.ReadCoverage(filepath.Join(repoPath, flag.CoverageFile), opts)
		if err != nil {
			return errors.Join(process.ErrReadCoverage, err)
		}

		coverage.PrintTabular(covData, os.Stdout)

		return nil
	},
}

func init() {
	flags := CoverageCmd.PersistentFlags()

	flags.StringVar(&flag.SortBy, flag.LongSort, coverage.Worst, "Specify")
	flags.BoolVarP(&flag.RunCoverage, flag.LongRunCoverage, flag.ShortRunCoverage, true,
		"Specify if tests with coverage should run")
	flags.StringVarP(&flag.CoverageFile, flag.LongFileCoverage, flag.ShortFileCoverage, "coverage.out",
		"Name of code coverage file to read or create")
	flags.BoolVarP(&flag.Verbose, flag.LongVerbose, flag.ShortVerbose, false, "Enable verbose output")
	flags.IntVarP(&flag.Top, flag.LongTop, flag.ShortTop, flag.DefaultTop, "Number of top files to display")
	flags.StringVar(&flag.ExcludePath, flag.LongExclude, "", "Exclude files matching regex")
}

func coverageOptsFromFlags() coverage.Options {
	opts := coverage.Options{
		SortBy:      flag.SortBy,
		Top:         flag.Top,
		ExcludePath: flag.ExcludePath,
	}

	return opts
}

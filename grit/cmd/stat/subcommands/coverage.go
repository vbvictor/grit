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
	Short: "Get coverage metrics of a repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		repoPath, err := filepath.Abs(args[0])
		if err != nil {
			return errors.Join(&process.ErrAbsRepoPath{Path: args[0]}, err)
		}

		if flag.Verbose {
			fmt.Printf("Processing repository: %s\n", repoPath)
		}

		if _, err := os.Stat(flag.CoverageFile); os.IsNotExist(err) {
			if flag.RunCoverage {
				if err = coverage.RunCoverage(repoPath, flag.CoverageFile); err != nil {
					return errors.Join(process.ErrRunCoverage, err)
				}
			} else {
				return errors.Join(process.ErrCoverageNotFound, err)
			}
		}

		covData, err := coverage.ReadCoverage(filepath.Join(repoPath, flag.CoverageFile), coverage.Options{})
		if err != nil {
			return errors.Join(process.ErrReadCoverage, err)
		}

		coverage.PrintTabular(covData, os.Stdout)

		return nil
	},
}

func init() {
	flags := ComplexityCmd.PersistentFlags()

	flags.BoolVar(&flag.RunCoverage, "cov-run", true, "Specify if tests with coverage should run")
	flags.StringVar(&flag.CoverageFile, "cov-file", "", "Path to coverage file")
}

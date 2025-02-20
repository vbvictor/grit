package report

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/process"
	"github.com/vbvictor/grit/pkg/report"
)

type factorError struct {
	factor string
	value  float64
}

func (e *factorError) Error() string {
	return fmt.Sprintf("%s factor is lower that 0: %f", e.factor, e.value)
}

var ReportCmd = &cobra.Command{
	Use:           "plot [flags] <repository>",
	Short:         "Compare code complexity and churn metrics",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	PreRunE: func(_ *cobra.Command, _ []string) error {
		return validateFactors(&report.ReportOpts)
	},
	RunE: func(_ *cobra.Command, args []string) error {
		repoPath, err := filepath.Abs(args[0])
		if err != nil {
			return errors.Join(&process.ErrAbsRepoPath{Path: args[0]}, err)
		}

		if flag.Verbose {
			fmt.Printf("Processing repository: %s\n", repoPath)
		}

		/*
			churns, err := git.MostGitChurnFiles(repoPath)
			if err != nil {
				return fmt.Errorf("error getting churn metrics: %w", err)
			}

			fileStat, err := complexity.RunComplexity(repoPath, complexity.Opts)
			if err != nil {
				return fmt.Errorf("error running complexity analysis: %w", err)
			}
		*/

		return nil
	},
}

func init() {
	flags := ReportCmd.PersistentFlags()

	flags.Float64VarP(&report.ReportOpts.ChurnFactor, "churn-mult", "c", 1.0, "Churn multiplier")
	flags.Float64VarP(&report.ReportOpts.ComplexityFactor, "comp-mult", "k", 1.0, "Complexity multiplier")
	flags.Float64VarP(&report.ReportOpts.CoverageFactor, "cov-mult", "v", 1.0, "Coverage multiplier")
	flags.BoolVar(&flag.RunCoverage, "cov-run", true, "Specify if tests with coverage should run")
	flags.StringVar(&flag.CoverageFile, "cov-file", "", "Path to coverage file")
}

func validateFactors(opts *report.Options) error {
	if opts.ChurnFactor < 0 {
		return &factorError{"Churn", opts.ChurnFactor}
	}

	if opts.ComplexityFactor < 0 {
		return &factorError{"Complexity", opts.ComplexityFactor}
	}

	if opts.CoverageFactor < 0 {
		return &factorError{"Coverage", opts.CoverageFactor}
	}

	return nil
}

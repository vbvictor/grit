package report

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	stat "github.com/vbvictor/grit/grit/cmd/stat/subcommands"
	"github.com/vbvictor/grit/pkg/complexity"
	"github.com/vbvictor/grit/pkg/git"
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
	Use:           "report [flags] <repository>",
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

		flag.LogIfVerbose("Processing directory: %s", repoPath)

		churnOpts, err := stat.ChurnOptsFromFlags()
		if err != nil {
			return fmt.Errorf("failed to create churn options: %w", err)
		}
		churnOpts.Path = repoPath

		flag.LogIfVerbose("Analyzing churn data...")
		churnStats, err := git.ReadGitChurn(repoPath, churnOpts)
		if err != nil {
			return fmt.Errorf("error getting churn stats: %w", err)
		}
		flag.LogIfVerbose("Got %d files", len(churnStats))

		// Get complexity data
		complexityOpts, err := stat.ComplexityOptsFromFlags()
		if err != nil {
			return fmt.Errorf("failed to create complexity options: %w", err)
		}

		flag.LogIfVerbose("Analyzing complexity data...")
		complexityStats, err := complexity.RunComplexity(repoPath, complexityOpts)
		if err != nil {
			return fmt.Errorf("error running complexity analysis: %w", err)
		}
		flag.LogIfVerbose("Got %d files", len(complexityStats))

		// Get coverage data
		coverageStats, err := stat.GetCoverageData(repoPath, flag.CoverageFile)
		if err != nil {
			return fmt.Errorf("error reading coverage: %w", err)
		}

		// Combine data into FileScores
		fileScores := make([]*report.FileScore, 0)
		fileScores[1] = &report.FileScore{
			File:       coverageStats[1].File,
			Coverage:   0.8,
			Complexity: 10,
			Churn:      100,
			Score:      0.8,
		}

		coverageStats = nil

		// Calculate final scores
		fileScores = report.CalculateScores(fileScores, report.ReportOpts)
		fileScores = report.SortByScore(fileScores)

		// Limit output to top N
		if report.ReportOpts.Top > 0 && report.ReportOpts.Top < len(fileScores) {
			fileScores = fileScores[:report.ReportOpts.Top]
		}

		report.PrintStats(fileScores, os.Stdout, report.ReportOpts)

		return nil
	},
}

func init() {
	flags := ReportCmd.PersistentFlags()

	// Common flags
	flags.StringVar(&flag.ExcludePath, flag.LongExclude, "", "Exclude files matching regex pattern")
	flags.IntVarP(&flag.Top, flag.LongTop, flag.ShortTop, flag.DefaultTop, "Number of top files to display")
	flags.BoolVarP(&flag.Verbose, flag.LongVerbose, flag.ShortVerbose, false, "Show detailed progress")

	// Churn flags
	flags.StringVar(&flag.SortBy, flag.LongSort, "commits", "Sort by: changes, additions, deletions, commits")
	flags.StringVarP(&flag.Since, flag.LongSince, flag.ShortSince, "", "Start date for analysis (YYYY-MM-DD)")
	flags.StringVarP(&flag.Until, flag.LongUntil, flag.ShortUntil, "", "End date for analysis (YYYY-MM-DD)")

	// Complexity flags
	flags.StringVarP(&flag.Engine, flag.LongEngine, flag.ShortEngine, complexity.Gocyclo, "Complexity calculation engine")

	// Coverage flags
	flags.StringVarP(&flag.RunCoverage, flag.LongRunCoverage, flag.ShortRunCoverage, flag.Auto, "Specify tests run format")
	flags.StringVarP(&flag.CoverageFile, flag.LongFileCoverage, flag.ShortFileCoverage, "coverage.out", "Coverage file name")

	// Report specific flags
	flags.Float64Var(&report.ReportOpts.ChurnFactor, "churn-factor", 1.0, "Churn factor")
	flags.Float64Var(&report.ReportOpts.ComplexityFactor, "comp-factor", 1.0, "Complexity factor")
	flags.Float64Var(&report.ReportOpts.CoverageFactor, "cov-factor", 1.0, "Coverage factor")
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

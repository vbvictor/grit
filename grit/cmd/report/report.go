package report

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/complexity"
	"github.com/vbvictor/grit/pkg/coverage"
	"github.com/vbvictor/grit/pkg/git"
	"github.com/vbvictor/grit/pkg/report"
)

type factorError struct {
	factor string
	value  float64
}

func (e *factorError) Error() string {
	return fmt.Sprintf("%s factor is lower that 0: %f", e.factor, e.value)
}

const (
	DefaultTop = 10
)

var (
	excludeRegex string
	top          int
	since        string
	until        string
	repoPath     string
)

var churnOpts = &git.ChurnOptions{
	SortBy:      git.Changes,
	Top:         0,
	Extensions:  nil,
	Since:       time.Time{},
	Until:       time.Time{},
	Path:        "",
	ExcludePath: "",
}

var complexityOpts = complexity.Options{
	Engine:      complexity.Gocyclo,
	ExcludePath: "",
	Top:         0, //nolint:mnd // default value
}

var coverageOpts = coverage.Options{
	SortBy:           coverage.Worst,
	Top:              0, //nolint:mnd // default value
	ExcludePath:      "",
	RunCoverage:      flag.Auto,
	CoverageFilename: "coverage.out",
}

var reportOpts = report.Options{
	Top:              DefaultTop,
	ExcludePath:      "",
	ChurnFactor:      1.0,
	ComplexityFactor: 1.0,
	CoverageFactor:   1.0,
}

var ReportCmd = &cobra.Command{
	Use:           "report [flags] <repository>",
	Short:         "Compare code complexity and churn metrics",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	PreRunE: func(_ *cobra.Command, _ []string) error {
		return validateFactors(&reportOpts)
	},
	RunE: func(_ *cobra.Command, args []string) error {
		var err error
		repoPath, err = filepath.Abs(args[0])
		if err != nil {
			return errors.Join(&flag.AbsRepoPathError{Path: args[0]}, err)
		}

		flag.LogIfVerbose("Processing directory: %s\n", repoPath)

		flag.LogIfVerbose("Analyzing churn data...\n")
		if err := git.PopulateOpts(churnOpts, []string{"go"}, since, until, repoPath); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}

		churns, err := git.ReadGitChurn(repoPath, churnOpts)
		if err != nil {
			return fmt.Errorf("error getting churn metrics: %w", err)
		}
		flag.LogIfVerbose("Got %d churn files\n", len(churns))

		flag.LogIfVerbose("Analyzing complexity data...\n")
		complexityStats, err := complexity.RunComplexity(repoPath, complexityOpts)
		if err != nil {
			return fmt.Errorf("error running complexity analysis: %w", err)
		}

		flag.LogIfVerbose("Got %d complexity files\n", len(complexityStats))

		flag.LogIfVerbose("Analyzing coverage data...\n")
		covData, err := coverage.GetCoverageData(repoPath, coverageOpts)
		if err != nil {
			return fmt.Errorf("failed to get coverage data: %w", err)
		}
		flag.LogIfVerbose("Got %d coverage files\n", len(covData))

		/*


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
		*/

		return nil
	},
}

func init() {
	flags := ReportCmd.PersistentFlags()

	// Common flags
	flags.StringVar(&excludeRegex, flag.LongExclude, "", "Exclude files matching regex pattern")
	flags.IntVarP(&top, flag.LongTop, flag.ShortTop, flag.DefaultTop, "Number of top files to display")
	flags.BoolVarP(&flag.Verbose, flag.LongVerbose, flag.ShortVerbose, false, "Show detailed progress")

	// Churn flags
	flags.StringVarP(&since, flag.LongSince, flag.ShortSince, "", "Start date for analysis in format 'YYYY-MM-DD'")
	flags.StringVarP(&until, flag.LongUntil, flag.ShortUntil, "", "End date for analysis in format 'YYYY-MM-DD'")

	// Complexity flags
	flags.StringVarP(&complexityOpts.Engine, flag.LongEngine, flag.ShortEngine, complexity.Gocyclo,
		"Complexity calculation engine")

	// Coverage flags
	flags.StringVarP(&coverageOpts.RunCoverage, flag.LongRunCoverage, flag.ShortRunCoverage, flag.Auto, "tests run format")
	flags.StringVarP(&coverageOpts.CoverageFilename, flag.LongFileCoverage, flag.ShortFileCoverage, "coverage.out",
		"Coverage file name")

	// Report specific flags
	flags.Float64Var(&reportOpts.ChurnFactor, "churn-factor", 1.0, "Churn factor")
	flags.Float64Var(&reportOpts.ComplexityFactor, "comp-factor", 1.0, "Complexity factor")
	flags.Float64Var(&reportOpts.CoverageFactor, "cov-factor", 1.0, "Coverage factor")
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

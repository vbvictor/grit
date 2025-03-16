package report

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/complexity"
	"github.com/vbvictor/grit/pkg/coverage"
	"github.com/vbvictor/grit/pkg/git"
	"github.com/vbvictor/grit/pkg/report"
)

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
	SortBy:       git.Changes,
	Top:          0,
	Extensions:   nil,
	Since:        time.Time{},
	Until:        time.Time{},
	Path:         "",
	ExcludeRegex: nil,
}

var complexityOpts = &complexity.Options{
	Engine:       complexity.Gocyclo,
	ExcludeRegex: nil,
	Top:          0, //nolint:mnd // default value
}

var coverageOpts = &coverage.Options{
	SortBy:           coverage.Worst,
	Top:              0, //nolint:mnd // default value
	ExcludeRegex:     nil,
	RunCoverage:      flag.Auto,
	CoverageFilename: "coverage.out",
}

var reportOpts = report.Options{
	Top:              DefaultTop,
	ExcludePath:      "",
	ChurnFactor:      1.0,
	ComplexityFactor: 1.0,
	CoverageFactor:   1.0,
	PerfectCoverage:  100.0, //nolint:mnd // default value
}

var ReportCmd = &cobra.Command{
	Use:           "report [flags] <repository>",
	Short:         "Creates maintainability report",
	Long:          "Creates maintainability report based on churn, complexity and coverage",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		var err error
		repoPath, err = filepath.Abs(args[0])
		if err != nil {
			return errors.Join(&flag.AbsRepoPathError{Path: args[0]}, err)
		}

		flag.LogIfVerbose("Processing directory: %s\n", repoPath)

		flag.LogIfVerbose("Analyzing churn data...\n")
		if err := git.PopulateOpts(churnOpts, []string{"go"}, since, until, repoPath, excludeRegex); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}

		churns, err := git.ReadGitChurn(repoPath, churnOpts)
		if err != nil {
			return fmt.Errorf("error getting churn metrics: %w", err)
		}
		flag.LogIfVerbose("Got %d churn files\n", len(churns))

		flag.LogIfVerbose("Analyzing complexity data...\n")
		if err := complexity.PopulateOpts(complexityOpts, excludeRegex); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}
		complexityStats, err := complexity.RunComplexity(repoPath, complexityOpts)
		if err != nil {
			return fmt.Errorf("error running complexity analysis: %w", err)
		}

		flag.LogIfVerbose("Got %d complexity files\n", len(complexityStats))

		flag.LogIfVerbose("Analyzing coverage data...\n")
		if err := coverage.PopulateOpts(coverageOpts, excludeRegex); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}
		covData, err := coverage.GetCoverageData(repoPath, coverageOpts)
		if err != nil {
			return fmt.Errorf("failed to get coverage data: %w", err)
		}
		flag.LogIfVerbose("Got %d coverage files\n", len(covData))

		fileScores := report.CombineMetrics(churns, complexityStats, covData)
		fileScores = report.SortByScore(report.CalculateScores(fileScores, reportOpts))
		flag.LogIfVerbose("Got %d file scores\n", len(fileScores))

		return report.PrintStats(fileScores, os.Stdout, reportOpts)
	},
}

func init() {
	flags := ReportCmd.PersistentFlags()

	// Common flags
	flag.ExcludeRegexFlag(flags, &excludeRegex)
	flag.TopFlag(flags, &top)
	flag.VerboseFlag(flags, &flag.Verbose)

	// Churn flags
	flag.SinceFlag(flags, &since)
	flag.UntilFlag(flags, &until)

	// Complexity flags
	flag.ComplexityEngineFlag(flags, &complexityOpts.Engine)

	// Coverage flags
	flag.RunCoverageFlag(flags, &coverageOpts.RunCoverage)
	flag.CoverageFilenameFlag(flags, &coverageOpts.CoverageFilename)

	// Report specific flags
	flag.PerfectCoverageFlag(flags, &reportOpts.PerfectCoverage)
}

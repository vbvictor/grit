package report

import (
	"fmt"
	"io"
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

var (
	excludeRegex string
	top          int
	since        string
	until        string
	outputFormat string
)

var churnOpts = &git.ChurnOptions{
	SortBy:       git.Commits,
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
	Top:              flag.DefaultTop,
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
		path := filepath.ToSlash(filepath.Clean(args[0]))
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("repository does not exist: %w", err)
		}

		flag.LogIfVerbose("Processing directory: %s\n", path)

		flag.LogIfVerbose("Analyzing churn data...\n")
		if err := git.PopulateOpts(churnOpts, []string{"go"}, since, until, path, excludeRegex); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}

		churns, err := git.ReadGitChurn(path, churnOpts)
		if err != nil {
			return fmt.Errorf("error getting churn metrics: %w", err)
		}
		flag.LogIfVerbose("Got %d churn files\n", len(churns))

		flag.LogIfVerbose("Analyzing complexity data...\n")
		if err := complexity.PopulateOpts(complexityOpts, excludeRegex); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}
		complexityStats, err := complexity.RunComplexity(path, complexityOpts)
		if err != nil {
			return fmt.Errorf("error running complexity analysis: %w", err)
		}

		flag.LogIfVerbose("Got %d complexity files\n", len(complexityStats))

		flag.LogIfVerbose("Analyzing coverage data...\n")
		if err := coverage.PopulateOpts(coverageOpts, excludeRegex); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}
		covData, err := coverage.GetCoverageData(path, coverageOpts)
		if err != nil {
			return fmt.Errorf("failed to get coverage data: %w", err)
		}
		flag.LogIfVerbose("Got %d coverage files\n", len(covData))

		fileScores := report.CombineMetrics(churns, complexityStats, covData)
		fileScores = report.SortAndLimit(report.CalculateScores(fileScores, reportOpts), top)
		flag.LogIfVerbose("Got %d file scores\n", len(fileScores))

		return printReport(fileScores, os.Stdout, &reportOpts, outputFormat)
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
	flag.ChurnTypeFlag(flags, &churnOpts.SortBy, git.Commits)

	// Complexity flags
	flag.ComplexityEngineFlag(flags, &complexityOpts.Engine)

	// Coverage flags
	flag.RunCoverageFlag(flags, &coverageOpts.RunCoverage)
	flag.CoverageFilenameFlag(flags, &coverageOpts.CoverageFilename)

	// Report specific flags
	flag.PerfectCoverageFlag(flags, &reportOpts.PerfectCoverage)
	flag.OutputFormatFlag(flags, &outputFormat)
}

func printReport(results []*report.FileScore, out io.Writer, opts *report.Options, format string) error {
	switch format {
	case flag.CSV:
		report.PrintCSV(results, out, opts)
	case flag.Tabular:
		report.PrintTabular(results, out, opts)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}

	return nil
}

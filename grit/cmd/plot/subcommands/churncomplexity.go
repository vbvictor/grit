package plot

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/complexity"
	"github.com/vbvictor/grit/pkg/git"
	"github.com/vbvictor/grit/pkg/plot"
)

var (
	outputFile   string
	since        string
	until        string
	churnType    git.ChurnType
	excludeRegex string
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

var ChurnComplexityCmd = &cobra.Command{
	Use:   "churn-vs-complexity [flags] <repository>",
	Short: "Creates churn vs complexity graph",
	Long: `
Creates churn vs complexity graph for a given repository.
Open generated file '.html' in a browser to view the graph.`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(_ *cobra.Command, _ []string) error {
		return plot.ValidateRiskThresholds()
	},
	RunE: func(_ *cobra.Command, args []string) error {
		path := filepath.Clean(args[0])
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("repository does not exist: %w", err)
		}

		flag.LogIfVerbose("Processing directory: %s\n", path)

		if err := git.PopulateOpts(churnOpts, []string{"go"}, since, until, path, excludeRegex); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}

		flag.LogIfVerbose("Analyzing churn data...\n")

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

		plotEntries := plot.PreparePlotData(complexityStats, churns, churnType)

		if err := plot.CreateScatterChart(plotEntries, plot.NewNoopMapper(), outputFile); err != nil {
			return fmt.Errorf("error creating chart: %w", err)
		}

		fmt.Printf("Chart generated: %s\n", outputFile)

		return nil
	},
}

func init() {
	flags := ChurnComplexityCmd.PersistentFlags()

	// Common flags
	flag.VerboseFlag(flags, &flag.Verbose)
	flag.OutputFlag(flags, &outputFile, "complexity_churn.html")
	flag.ExcludeRegexFlag(flags, &excludeRegex)
	flag.ChurnTypeFlag(flags, &churnType, git.Commits)

	// Churn flags
	flag.SinceFlag(flags, &since)
	flag.UntilFlag(flags, &until)

	// Complexity flags
	flag.EngineFlag(flags, &complexityOpts.Engine, complexity.Gocyclo)

	ChurnComplexityCmd.Flag(flag.LongUntil).DefValue = flag.DefaultUntil
	ChurnComplexityCmd.Flag(flag.LongSince).DefValue = flag.DefaultSince

	ChurnComplexityCmd.Flag(flag.LongUntil).DefValue = flag.DefaultUntil
	ChurnComplexityCmd.Flag(flag.LongSince).DefValue = flag.DefaultSince
}

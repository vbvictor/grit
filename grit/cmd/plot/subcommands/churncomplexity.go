package plot

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/complexity"
	"github.com/vbvictor/grit/pkg/git"
	"github.com/vbvictor/grit/pkg/plot"
)

var (
	outputFile string
	since      string
	until      string
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
		repoPath, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("error getting absolute path: %w", err)
		}

		flag.LogIfVerbose("Processing directory: %s\n", repoPath)

		if err := git.PopulateOpts(churnOpts, []string{"go"}, since, until, repoPath); err != nil {
			return fmt.Errorf("failed to create options: %w", err)
		}

		flag.LogIfVerbose("Analyzing churn data...\n")

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

		plotEntries := plot.PreparePlotData(complexityStats, churns)

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
	flags.BoolVarP(&flag.Verbose, flag.LongVerbose, flag.ShortVerbose, false, "Show detailed progress")
	flags.StringVarP(&outputFile, "output", "o", "complexity_churn.html", "Output graph file name")

	flags.StringVarP(&plot.Plot, "plot-type", "t", git.Commits,
		fmt.Sprintf("Specify churn type: [%s, %s, %s, %s]", git.Changes, git.Additions, git.Deletions, git.Commits))

	// Churn flags
	flags.StringVarP(&since, flag.LongSince, flag.ShortSince, "", "Start date for churn analysis (YYYY-MM-DD)")
	flags.StringVarP(&until, flag.LongUntil, flag.ShortUntil, "", "End date for churn analysis (YYYY-MM-DD)")

	// Complexity flags
	flags.StringVarP(&complexityOpts.Engine, flag.LongEngine, flag.ShortEngine, complexity.Gocyclo,
		fmt.Sprintf("Specify complexity calculation engine: [%s, %s]", complexity.Gocyclo, complexity.Gocognit))

	ChurnComplexityCmd.Flag(flag.LongUntil).DefValue = flag.DefaultUntil
	ChurnComplexityCmd.Flag(flag.LongSince).DefValue = flag.DefaultSince
}

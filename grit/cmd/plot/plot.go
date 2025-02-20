package plot

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/plot"
)

var (
	outputFile string
	PlotCmd    = &cobra.Command{
		Use:   "plot [flags] <repository>",
		Short: "Compare code complexity and churn metrics",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return plot.ValidateRiskThresholds()
		},
		RunE: func(_ *cobra.Command, args []string) error {
			repoPath, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("error getting absolute path: %w", err)
			}

			if flag.Verbose {
				fmt.Printf("Processing repository: %s\n", repoPath)
			}

			//TODO(v.baranov) refactor
			/*
				churns, err := git.MostGitChurnFiles(repoPath)
				if err != nil {
					return fmt.Errorf("error reading churn data: %w", err)
				}


					if err := complexity.CheckLizardExecutable(); err != nil {

						return fmt.Errorf("failed to find lizard Executable: %w", err)
					}

					fileStat, err := complexity.RunLizardCmd(repoPath, complexity.Opts)
					if err != nil {
						return fmt.Errorf("error running lizard command: %w", err)
					}

					fileStat = complexity.ApplyFilters(fileStat,
						complexity.MinComplexityFilter{MinComplexity: complexity.MinComplexityDefault}.Filter)
					plotEntries := plot.PreparePlotData(fileStat, churns)

					if err := plot.CreateScatterChart(plotEntries, plot.NewRisksMapper(), outputFile); err != nil {
						return fmt.Errorf("error creating chart: %w", err)
					}

					if flag.Verbose {
						fmt.Printf("Chart generated: %s\n", outputFile)
					}
			*/

			return nil
		},
	}
)

func init() {
	// flags := PlotCmd.PersistentFlags()
	/*
		flags.StringVarP(&outputFile, "output", "o", "complexity_churn.html", "Output file path")
		flags.BoolVarP(&flag.Verbose, "verbose", "v", false, "Enable verbose output")
		flags.StringVarP(&plot.Plot, "plot-type", "t", "changes", "Specify OY plot type")
		flags.IntVar(&git.ChurnOpts.Top, "top", git.DefaultTop, "Number of top files to display")
		flags.StringVar(&git.ChurnOpts.ExcludePath, "exclude", "", "Exclude files matching regex pattern")
		flags.StringSliceVar(&process.Extensions, "ext", nil,
			"Only include files with given extensions in comma-separated list. For example go,h,c")
		flags.Var(&git.ChurnOpts.Since, "since", "Start date for analysis (YYYY-MM-DD)")
		flags.Var(&git.ChurnOpts.Until, "until", "End date for analysis (YYYY-MM-DD)")
	*/
	// flags.IntVar(&complexity.Opts.Threads, "t", 1, "Number of threads to run")
}

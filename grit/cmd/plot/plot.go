package plot

import (
	"github.com/spf13/cobra"
	plot "github.com/vbvictor/grit/grit/cmd/plot/subcommands"
)

var PlotCmd = &cobra.Command{ //nolint:exhaustruct // no need to set all fields
	Use:   "plot",
	Short: `Creates graphs for code metrics`,
	Long: `
Creates visual graphs for code metrics.
Graphs are rendered using echarts library, browser is required to view them.`,
}

func init() {
	PlotCmd.AddCommand(plot.ChurnComplexityCmd)
}

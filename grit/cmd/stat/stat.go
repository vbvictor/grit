package stat

import (
	"github.com/spf13/cobra"
	stat "github.com/vbvictor/grit/grit/cmd/stat/subcommands"
)

var StatCmd = &cobra.Command{
	Use:   "stat",
	Short: "Calculate code metrics",
}

func init() {
	StatCmd.AddCommand(stat.ChurnCmd)
	StatCmd.AddCommand(stat.ComplexityCmd)
	StatCmd.AddCommand(stat.CoverageCmd)
}

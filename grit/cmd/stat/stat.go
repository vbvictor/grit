package stat

import (
	"github.com/spf13/cobra"
	stat "github.com/vbvictor/grit/grit/cmd/stat/subcommands"
)

var StatCmd = &cobra.Command{
	Use:   "stat",
	Short: "Some short stat description.",
	Long:  `Some long stat description.`,
}

func init() {
	StatCmd.AddCommand(stat.ChurnCmd)
	StatCmd.AddCommand(stat.ComplexityCmd)
	StatCmd.AddCommand(stat.CoverageCmd)
}

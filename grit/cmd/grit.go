package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/grit/cmd/stat"
)

var gritCmd = &cobra.Command{
	Use: "grit",
	Long: `GRIT is an all-in-one cli tool that helps developers understand
their codebase maintainability index through calculated metrics.`,
}

func Execute() {
	if err := gritCmd.Execute(); err != nil {
		// Errors are stored in pairs: first one is pretty-printed, second one is raw error from go-code.
		if uw, ok := err.(interface{ Unwrap() []error }); ok {
			errs := uw.Unwrap()

			gritCmd.PrintErrln(errs[0].Error())

			if flag.Verbose {
				gritCmd.PrintErrln(errs[1].Error())
			}
		} else {
			gritCmd.PrintErrln(err.Error())
		}

		os.Exit(1)
	}
}

func init() {
	gritCmd.Root().CompletionOptions.DisableDefaultCmd = true

	gritCmd.AddCommand(stat.StatCmd)
}

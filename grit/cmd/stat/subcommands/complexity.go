package stat

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vbvictor/grit/grit/cmd/flag"
	"github.com/vbvictor/grit/pkg/complexity"
)

var ComplexityCmd = &cobra.Command{
	Use:   "complexity <path>",
	Short: "Get complexity metrics of a repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		repoPath, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("error getting absolute path: %w", err)
		}

		if flag.Verbose {
			fmt.Printf("Processing repository: %s\n", repoPath)
		}

		fileStat, err := complexity.RunComplexity(repoPath, complexity.Opts)
		if err != nil {
			return fmt.Errorf("error running complexity analysis: %w", err)
		}

		complexity.PrintTabular(fileStat, os.Stdout)

		return nil
	},
}

func init() {
	flags := ComplexityCmd.PersistentFlags()

	flags.StringVar(&complexity.Opts.Engine, "engine", complexity.Gocyclo,
		"Complexity calculation engine to use: gocyclo or gocognit")
}

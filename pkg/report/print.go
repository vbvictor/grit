package report

import (
	"fmt"
	"io"

	"github.com/bndr/gotabulate"
)

func PrintStats(results []*FileScore, out io.Writer, opts Options) error {
	return printTabular(results, out, opts)
}

func printTabular(results []*FileScore, out io.Writer, opts Options) error {
	fmt.Fprintf(out, "\nCode health analysis results (top %d):\n", opts.Top)

	data := make([][]any, len(results))
	for i, result := range results {
		data[i] = []any{
			result.File,
			fmt.Sprintf("%.2f", result.Score),
			fmt.Sprintf("%.2f", result.Churn),
			fmt.Sprintf("%.2f", result.Complexity),
			fmt.Sprintf("%.2f%%", result.Coverage),
		}
	}

	table := gotabulate.Create(data)
	table.SetHeaders([]string{"FILEPATH", "SCORE", "CHURN", "COMPLEXITY", "COVERAGE"})
	table.SetAlign("left")

	if _, err := io.WriteString(out, table.Render("grid")); err != nil {
		return fmt.Errorf("failed to write grid: %w", err)
	}

	return nil
}

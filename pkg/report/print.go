package report

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/bndr/gotabulate"
)

func PrintTabular(results []*FileScore, out io.Writer, opts *Options) error {
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

func PrintCSV(results []*FileScore, out io.Writer, _ *Options) error {
	writer := csv.NewWriter(out)
	defer writer.Flush()

	// Write headers
	if err := writer.Write([]string{"FILEPATH", "SCORE", "CHURN", "COMPLEXITY", "COVERAGE"}); err != nil {
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// Write data
	for _, result := range results {
		record := []string{
			result.File,
			fmt.Sprintf("%.2f", result.Score),
			fmt.Sprintf("%.2f", result.Churn),
			fmt.Sprintf("%.2f", result.Complexity),
			fmt.Sprintf("%.2f", result.Coverage),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

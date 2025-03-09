package coverage

import (
	"fmt"
	"io"

	"github.com/bndr/gotabulate"
)

func PrintTabular(results []*FileCoverage, out io.Writer) {
	fmt.Fprintf(out, "\nCode coverage analysis results:\n")

	data := make([][]any, len(results))
	for i, result := range results {
		data[i] = []any{
			result.File,
			fmt.Sprintf("%.2f%%", result.Coverage),
			result.Statements,
			result.Covered,
		}
	}

	table := gotabulate.Create(data)
	table.SetHeaders([]string{"FILEPATH", "COVERAGE", "STATEMENTS", "COVERED"})
	table.SetAlign("left")

	_, _ = io.WriteString(out, table.Render("grid"))
}

func PrintCSV(results []*FileCoverage, out io.Writer) {
	_, _ = fmt.Fprintln(out, "filepath,coverage,statements,covered")

	for _, result := range results {
		_, _ = fmt.Fprintf(out, "%s,%.2f,%d,%d\n",
			result.File,
			result.Coverage,
			result.Statements,
			result.Covered,
		)
	}
}

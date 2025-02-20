package complexity

import (
	"io"

	"github.com/bndr/gotabulate"
)

func PrintTabular(results []*FileStat, out io.Writer) {
	_, _ = io.WriteString(out, "\nCode complexity analysis results:\n")

	data := make([][]interface{}, len(results))
	for i, result := range results {
		data[i] = []interface{}{result.Path, result.AvgComplexity}
	}

	table := gotabulate.Create(data)
	table.SetHeaders([]string{"FILEPATH", "COMPLEXITY"})
	table.SetAlign("left")

	_, _ = io.WriteString(out, table.Render("grid"))
}

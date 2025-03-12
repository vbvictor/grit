package complexity

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/bndr/gotabulate"
	"github.com/vbvictor/grit/grit/cmd/flag"
)

func PrintStats(results []*FileStat, out io.Writer, opts *Options) error {
	switch opts.OutputFormat {
	case flag.CSV:
		PrintCSV(results, out)
	case flag.Tabular:
		PrintTabular(results, out)
	default:
		return fmt.Errorf("unsupported output format: %s", opts.OutputFormat)
	}

	return nil
}

func PrintTabular(results []*FileStat, out io.Writer) {
	_, _ = io.WriteString(out, "\nCode complexity analysis results:\n")

	data := make([][]any, len(results))
	for i, result := range results {
		data[i] = []any{result.Path, result.AvgComplexity}
	}

	table := gotabulate.Create(data)
	table.SetHeaders([]string{"FILEPATH", "COMPLEXITY"})
	table.SetAlign("left")

	_, _ = io.WriteString(out, table.Render("grid"))
}

func PrintCSV(results []*FileStat, out io.Writer) {
	writer := csv.NewWriter(out)
	defer writer.Flush()

	_ = writer.Write([]string{"FILEPATH", "COMPLEXITY"})

	for _, result := range results {
		record := []string{
			result.Path,
			strconv.FormatFloat(result.AvgComplexity, 'f', 2, 64),
		}
		_ = writer.Write(record)
	}
}

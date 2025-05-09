package git

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/bndr/gotabulate"
)

func PrintTable(results []*ChurnChunk, out io.Writer, opts *ChurnOptions) {
	fmt.Fprintf(out, "\nTop %d most modified files by %s:\n", opts.Top, opts.SortBy)

	data := make([][]any, len(results))

	for i, result := range results {
		data[i] = []any{result.Churn, result.Added, result.Removed, result.Commits, result.File}
	}

	table := gotabulate.Create(data)
	table.SetHeaders([]string{"CHANGES", "ADDED", "DELETED", "COMMITS", "FILEPATH"})
	table.SetAlign("left")

	_, _ = io.WriteString(out, table.Render("grid"))
}

func PrintCSV(results []*ChurnChunk, out io.Writer, _ *ChurnOptions) {
	writer := csv.NewWriter(out)
	defer writer.Flush()

	_ = writer.Write([]string{"FILEPATH", "CHANGES", "ADDED", "DELETED", "COMMITS"})

	for _, result := range results {
		record := []string{
			result.File,
			strconv.Itoa(result.Churn),
			strconv.Itoa(result.Added),
			strconv.Itoa(result.Removed),
			strconv.Itoa(result.Commits),
		}
		_ = writer.Write(record)
	}
}

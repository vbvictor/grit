package git

import (
	"fmt"
	"io"

	"github.com/bndr/gotabulate"
)

func PrintStats(results []*ChurnChunk, out io.Writer, opts ChurnOptions) error {
	printTable(results, out, opts)

	return nil
}

func printTable(results []*ChurnChunk, out io.Writer, opts ChurnOptions) {
	fmt.Fprintf(out, "\nTop %d most modified files by %s:\n", opts.Top, opts.SortBy)

	data := make([][]interface{}, len(results))

	for i, result := range results {
		data[i] = []interface{}{result.Churn, result.Added, result.Removed, result.Commits, result.File}
	}

	table := gotabulate.Create(data)
	table.SetHeaders([]string{"CHANGES", "ADDED", "DELETED", "COMMITS", "FILEPATH"})
	table.SetAlign("left")

	_, _ = io.WriteString(out, table.Render("grid"))
}

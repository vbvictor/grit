package git

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintTable(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*ChurnChunk
		opts     *ChurnOptions
		expected []string
	}{
		{
			name: "single file churn",
			input: []*ChurnChunk{
				{
					File:    "main.go",
					Churn:   100,
					Added:   80,
					Removed: 20,
					Commits: 5,
				},
			},
			opts: &ChurnOptions{
				Top:    1,
				SortBy: "changes",
			},
			expected: []string{
				"Top 1 most modified files by changes",
				"CHANGES", "ADDED", "DELETED", "COMMITS", "FILEPATH",
				"100", "80", "20", "5", "main.go",
			},
		},
		{
			name: "multiple files churn",
			input: []*ChurnChunk{
				{
					File:    "path/to/foo.go",
					Churn:   50,
					Added:   30,
					Removed: 20,
					Commits: 3,
				},
				{
					File:    "bar.go",
					Churn:   150,
					Added:   100,
					Removed: 50,
					Commits: 10,
				},
			},
			opts: &ChurnOptions{
				Top:    2,
				SortBy: "commits",
			},
			expected: []string{
				"Top 2 most modified files by commits",
				"CHANGES", "ADDED", "DELETED", "COMMITS", "FILEPATH",
				"50", "30", "20", "3", "path/to/foo.go",
				"150", "100", "50", "10", "bar.go",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			printTable(tc.input, &buf, tc.opts)

			output := buf.String()
			for _, exp := range tc.expected {
				if !strings.Contains(output, exp) {
					t.Errorf("Expected output to contain %q, but it didn't.\nGot: %s", exp, output)
				}
			}
		})
	}
}

func TestPrintCSV(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*ChurnChunk
		opts     *ChurnOptions
		expected [][]string
	}{
		{
			name: "single file churn",
			input: []*ChurnChunk{
				{
					File:    "main.go",
					Churn:   100,
					Added:   80,
					Removed: 20,
					Commits: 5,
				},
			},
			opts: &ChurnOptions{
				Top:    1,
				SortBy: "changes",
			},
			expected: [][]string{
				{"CHANGES", "ADDED", "DELETED", "COMMITS", "FILEPATH"},
				{"100", "80", "20", "5", "main.go"},
			},
		},
		{
			name: "multiple files churn",
			input: []*ChurnChunk{
				{
					File:    "path/to/foo.go",
					Churn:   50,
					Added:   30,
					Removed: 20,
					Commits: 3,
				},
				{
					File:    "bar.go",
					Churn:   150,
					Added:   100,
					Removed: 50,
					Commits: 10,
				},
			},
			opts: &ChurnOptions{
				Top:    2,
				SortBy: "commits",
			},
			expected: [][]string{
				{"CHANGES", "ADDED", "DELETED", "COMMITS", "FILEPATH"},
				{"50", "30", "20", "3", "path/to/foo.go"},
				{"150", "100", "50", "10", "bar.go"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			printCSV(tc.input, &buf, &ChurnOptions{})

			reader := csv.NewReader(bytes.NewReader(buf.Bytes()))
			output, err := reader.ReadAll()
			require.NoError(t, err, "Failed to parse CSV output")

			assert.Equal(t, len(tc.expected), len(output), "Row count mismatch")

			for i, expectedRow := range tc.expected {
				assert.ElementsMatch(t, expectedRow, output[i])
			}
		})
	}
}

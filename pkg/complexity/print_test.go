package complexity

import (
	"bytes"
	"encoding/csv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintTabular(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*FileStat
		expected []string
	}{
		{
			name: "single file complexity",
			input: []*FileStat{
				{
					Path:          "main.go",
					AvgComplexity: 4,
				},
			},
			expected: []string{
				"main.go",
				"4",
			},
		},
		{
			name: "multiple files complexity",
			input: []*FileStat{
				{
					Path:          "path/to/foo.go",
					AvgComplexity: 2,
				},
				{
					Path:          "bar.go",
					AvgComplexity: 6,
				},
			},
			expected: []string{
				"path/to/foo.go",
				"2",
				"bar.go",
				"6",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			PrintTabular(tc.input, &buf)

			output := buf.String()
			for _, exp := range tc.expected {
				assert.Contains(t, output, exp, "Expected output to contain %q", exp)
			}
		})
	}
}

func TestPrintCSV(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*FileStat
		expected [][]string
	}{
		{
			name: "single file complexity",
			input: []*FileStat{
				{
					Path:          "main.go",
					AvgComplexity: 4,
				},
			},
			expected: [][]string{
				{"FILEPATH", "COMPLEXITY"},
				{"main.go", "4.00"},
			},
		},
		{
			name: "multiple files complexity",
			input: []*FileStat{
				{
					Path:          "path/to/foo.go",
					AvgComplexity: 2,
				},
				{
					Path:          "bar.go",
					AvgComplexity: 6,
				},
			},
			expected: [][]string{
				{"FILEPATH", "COMPLEXITY"},
				{"path/to/foo.go", "2.00"},
				{"bar.go", "6.00"},
			},
		},
		{
			name:  "empty input",
			input: []*FileStat{},
			expected: [][]string{
				{"FILEPATH", "COMPLEXITY"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			PrintCSV(tc.input, &buf)

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

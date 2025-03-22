package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrintTabular(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*FileScore
		expected []string
	}{
		{
			name: "single file score",
			input: []*FileScore{
				{
					File:       "main.go",
					Coverage:   75.5,
					Complexity: 4.2,
					Churn:      100,
					Score:      42.5,
				},
			},
			expected: []string{
				"main.go",
				"75.50%",
				"4.20",
				"100.00",
				"42.50",
			},
		},
		{
			name: "multiple files score",
			input: []*FileScore{
				{
					File:       "path/to/foo.go",
					Coverage:   90.0,
					Complexity: 2.5,
					Churn:      50,
					Score:      20.5,
				},
				{
					File:       "bar.go",
					Coverage:   60.5,
					Complexity: 6.0,
					Churn:      150,
					Score:      85.2,
				},
			},
			expected: []string{
				"path/to/foo.go",
				"foo",
				"90.00%",
				"2.50",
				"50.00",
				"20.50",
				"bar.go",
				"60.50%",
				"6.00",
				"150.00",
				"85.20",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := PrintTabular(tc.input, &buf, &Options{})
			require.NoError(t, err)

			output := buf.String()
			for _, exp := range tc.expected {
				require.Contains(t, output, exp,
					"Expected output to contain %q, but it didn't.\nGot: %s", exp, output)
			}
		})
	}
}

func TestPrintCSV(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*FileScore
		expected []string
	}{
		{
			name: "single file score",
			input: []*FileScore{
				{
					File:       "main.go",
					Coverage:   75.5,
					Complexity: 4.2,
					Churn:      100,
					Score:      42.5,
				},
			},
			expected: []string{
				"FILEPATH,SCORE,CHURN,COMPLEXITY,COVERAGE",
				"main.go,42.50,100.00,4.20,75.50",
			},
		},
		{
			name: "multiple files score",
			input: []*FileScore{
				{
					File:       "path/to/foo.go",
					Coverage:   90.0,
					Complexity: 2.5,
					Churn:      50,
					Score:      20.5,
				},
				{
					File:       "bar.go",
					Coverage:   60.5,
					Complexity: 6.0,
					Churn:      150,
					Score:      85.2,
				},
			},
			expected: []string{
				"FILEPATH,SCORE,CHURN,COMPLEXITY,COVERAGE",
				"path/to/foo.go,20.50,50.00,2.50,90.00",
				"bar.go,85.20,150.00,6.00,60.50",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := PrintCSV(tc.input, &buf, &Options{})
			require.NoError(t, err)

			output := buf.String()
			lines := strings.Split(output, "\n")

			require.Equal(t, len(tc.expected), len(lines)-1) // -1 for trailing newline
			for i, exp := range tc.expected {
				require.Equal(t, exp, lines[i])
			}
		})
	}
}

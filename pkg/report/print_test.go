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

			err := printTabular(tc.input, &buf, Options{})
			require.NoError(t, err)

			output := buf.String()
			for _, exp := range tc.expected {
				if !strings.Contains(output, exp) {
					t.Errorf("Expected output to contain %q, but it didn't.\nGot: %s", exp, output)
				}
			}
		})
	}
}

package coverage

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintTabular(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*FileCoverage
		expected []string
	}{
		{
			name: "single file coverage",
			input: []*FileCoverage{
				{
					File:       "main.go",
					Coverage:   75.5,
					Statements: 100,
					Covered:    75,
				},
			},
			expected: []string{
				"main.go",
				"75.50%",
				"100",
				"75",
			},
		},
		{
			name: "multiple files coverage",
			input: []*FileCoverage{
				{
					File:       "path/to/foo.go",
					Coverage:   90.0,
					Statements: 50,
					Covered:    45,
				},
				{
					File:       "bar.go",
					Coverage:   60.5,
					Statements: 200,
					Covered:    121,
				},
			},
			expected: []string{
				"path/to/foo.go",
				"90.00%",
				"50",
				"45",
				"bar.go",
				"60.50%",
				"200",
				"121",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			PrintTabular(tc.input, &buf)

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
		input    []*FileCoverage
		expected string
	}{
		{
			name: "single file coverage",
			input: []*FileCoverage{
				{
					File:       "main.go",
					Coverage:   75.5,
					Statements: 100,
					Covered:    75,
				},
			},
			expected: "filepath,coverage,statements,covered\nmain.go,75.50,100,75\n",
		},
		{
			name: "multiple files coverage",
			input: []*FileCoverage{
				{
					File:       "path/to/foo.go",
					Coverage:   90.0,
					Statements: 50,
					Covered:    45,
				},
				{
					File:       "bar.go",
					Coverage:   60.5,
					Statements: 200,
					Covered:    121,
				},
			},
			expected: "filepath,coverage,statements,covered\npath/to/foo.go,90.00,50,45\nbar.go,60.50,200,121\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			PrintCSV(tc.input, &buf)

			if got := buf.String(); got != tc.expected {
				t.Errorf("PrintCSV() = %q, want %q", got, tc.expected)
			}
		})
	}
}

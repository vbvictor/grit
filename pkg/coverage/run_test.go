package coverage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadCoverage(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []*FileCoverage
		wantErr error
	}{
		{
			name: "Valid set mode coverage",
			content: `mode: set
example.com/pkg/file1.go:10.20,30.2 3 1
example.com/pkg/file1.go:32.20,35.2 2 0
example.com/pkg/file2.go:5.20,8.2 2 1`,
			want: []*FileCoverage{
				{
					File:       "example.com/pkg/file2.go",
					Coverage:   100.0,
					Statements: 2,
					Covered:    2,
				},
				{
					File:       "example.com/pkg/file1.go",
					Coverage:   60.0,
					Statements: 5,
					Covered:    3,
				},
			},
		},
		{
			name: "Unsupported count mode",
			content: `mode: count
example.com/pkg/file1.go:10.20,30.2 3 1`,
			wantErr: errUnsupportedMode,
		},
		{
			name: "Unsupported atomic mode",
			content: `mode: atomic
example.com/pkg/file1.go:10.20,30.2 3 1`,
			wantErr: errUnsupportedMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp(t.TempDir(), "coverage.*.out")
			require.NoError(t, err, "Failed to create temp file")
			defer os.Remove(tmpfile.Name())

			_, err = tmpfile.WriteString(tt.content)
			require.NoError(t, err, "Failed to write to temp file")

			err = tmpfile.Close()
			require.NoError(t, err, "Failed to close temp file")

			got, err := ReadCoverage(tmpfile.Name(), Options{Top: 10, SortBy: Worst, ExcludePath: ""})

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)

				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tt.want), len(got), "Results length mismatch")

			for i, want := range tt.want {
				assert.Contains(t, got, want, "Result mismatch at index %d", i)
			}
		})
	}
}

func TestSortByCoverage(t *testing.T) {
	tests := []struct {
		name     string
		input    []FileCoverage
		asc      bool
		expected []FileCoverage
	}{
		{
			name: "Sort ascending",
			input: []FileCoverage{
				{File: "file1.go", Coverage: 75.0},
				{File: "file2.go", Coverage: 100.0},
				{File: "file3.go", Coverage: 50.0},
			},
			asc: true,
			expected: []FileCoverage{
				{File: "file3.go", Coverage: 50.0},
				{File: "file1.go", Coverage: 75.0},
				{File: "file2.go", Coverage: 100.0},
			},
		},
		{
			name: "Sort descending",
			input: []FileCoverage{
				{File: "file1.go", Coverage: 75.0},
				{File: "file2.go", Coverage: 100.0},
				{File: "file3.go", Coverage: 50.0},
			},
			asc: false,
			expected: []FileCoverage{
				{File: "file2.go", Coverage: 100.0},
				{File: "file1.go", Coverage: 75.0},
				{File: "file3.go", Coverage: 50.0},
			},
		},
		{
			name:     "Empty slice",
			input:    []FileCoverage{},
			asc:      true,
			expected: []FileCoverage{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of input to verify it's not modified
			original := make([]FileCoverage, len(tt.input))
			copy(original, tt.input)

			result := sortByCoverage(tt.input, tt.asc)

			assert.Equal(t, tt.expected, result, "Incorrect sorting result")
			assert.Equal(t, original, tt.input, "Original slice was modified")
		})
	}
}

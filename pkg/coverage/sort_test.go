package coverage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

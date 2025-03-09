package report

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateScores(t *testing.T) {
	tests := []struct {
		name     string
		input    []*FileScore
		opts     Options
		expected []*FileScore
	}{
		{
			name: "calculate with churn factor zero",
			input: []*FileScore{
				{File: "test.go", Complexity: 10, Coverage: 80, Churn: 5},
			},
			opts: Options{
				ChurnFactor:      0,
				ComplexityFactor: 2.0,
				CoverageFactor:   1.0,
			},
			expected: []*FileScore{
				{File: "test.go", Complexity: 10, Coverage: 80, Churn: 5, Score: (10 * 2.0) / (80 + 1.0)},
			},
		},
		{
			name: "calculate with complexity factor zero",
			input: []*FileScore{
				{File: "test.go", Complexity: 10, Coverage: 80, Churn: 5},
			},
			opts: Options{
				ChurnFactor:      2.0,
				ComplexityFactor: 0,
				CoverageFactor:   1.0,
			},
			expected: []*FileScore{
				{File: "test.go", Complexity: 10, Coverage: 80, Churn: 5, Score: (5 * 2.0) / (80 + 1.0)},
			},
		},
		{
			name: "calculate with coverage factor zero",
			input: []*FileScore{
				{File: "test.go", Complexity: 10, Coverage: 80, Churn: 5},
			},
			opts: Options{
				ChurnFactor:      2.0,
				ComplexityFactor: 3.0,
				CoverageFactor:   0,
			},
			expected: []*FileScore{
				{File: "test.go", Complexity: 10, Coverage: 80, Churn: 5, Score: (5 * 2.0) + (10 * 3.0)},
			},
		},
		{
			name: "calculate with all factors non-zero",
			input: []*FileScore{
				{File: "test.go", Complexity: 10, Coverage: 80, Churn: 5},
			},
			opts: Options{
				ChurnFactor:      2.0,
				ComplexityFactor: 3.0,
				CoverageFactor:   1.0,
			},
			expected: []*FileScore{
				{File: "test.go", Complexity: 10, Coverage: 80, Churn: 5, Score: (5 * 2.0) + (10*3.0)/(80+1.0)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateScores(tt.input, tt.opts)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestSortByScore(t *testing.T) {
	tests := []struct {
		name     string
		input    []*FileScore
		expected []*FileScore
	}{
		{
			name:     "sort empty slice",
			input:    []*FileScore{},
			expected: []*FileScore{},
		},
		{
			name: "sort already sorted slice",
			input: []*FileScore{
				{File: "high.go", Score: 100},
				{File: "medium.go", Score: 50},
				{File: "low.go", Score: 10},
			},
			expected: []*FileScore{
				{File: "high.go", Score: 100},
				{File: "medium.go", Score: 50},
				{File: "low.go", Score: 10},
			},
		},
		{
			name: "sort reverse ordered slice",
			input: []*FileScore{
				{File: "low.go", Score: 10},
				{File: "medium.go", Score: 50},
				{File: "high.go", Score: 100},
			},
			expected: []*FileScore{
				{File: "high.go", Score: 100},
				{File: "medium.go", Score: 50},
				{File: "low.go", Score: 10},
			},
		},
		{
			name: "sort equal scores",
			input: []*FileScore{
				{File: "a.go", Score: 50},
				{File: "b.go", Score: 50},
				{File: "c.go", Score: 50},
			},
			expected: []*FileScore{
				{File: "a.go", Score: 50},
				{File: "b.go", Score: 50},
				{File: "c.go", Score: 50},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortByScore(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

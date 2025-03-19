package report

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vbvictor/grit/pkg/complexity"
	"github.com/vbvictor/grit/pkg/coverage"
	"github.com/vbvictor/grit/pkg/git"
)

func TestCalculateScores(t *testing.T) {
	tests := []struct {
		name     string
		input    []*FileScore
		opts     Options
		expected []*FileScore
	}{
		{
			name: "calculate with zero churn",
			input: []*FileScore{
				{File: "test.go", Complexity: 10, Coverage: 80, Churn: 0},
			},
			opts: Options{
				PerfectCoverage: 100,
			},
			expected: []*FileScore{
				{
					File:            "test.go",
					Complexity:      10,
					Coverage:        80,
					Churn:           0,
					ChurnComplexity: 10,              // When Churn is 0, ChurnComplexity = Complexity
					Score:           10 * (100 - 80), // Score = ChurnComplexity * (PerfectCoverage - Coverage)
				},
			},
		},
		{
			name: "calculate with zero complexity",
			input: []*FileScore{
				{File: "test.go", Complexity: 0, Coverage: 80, Churn: 5},
			},
			opts: Options{
				PerfectCoverage: 100,
			},
			expected: []*FileScore{
				{
					File:            "test.go",
					Complexity:      0,
					Coverage:        80,
					Churn:           5,
					ChurnComplexity: 5,              // When Complexity is 0, ChurnComplexity = Churn
					Score:           5 * (100 - 80), // Score = ChurnComplexity * (PerfectCoverage - Coverage)
				},
			},
		},
		{
			name: "calculate with non-zero churn and complexity",
			input: []*FileScore{
				{File: "test.go", Complexity: 10, Coverage: 80, Churn: 5},
			},
			opts: Options{
				PerfectCoverage: 100,
			},
			expected: []*FileScore{
				{
					File:            "test.go",
					Complexity:      10,
					Coverage:        80,
					Churn:           5,
					ChurnComplexity: 50,              // ChurnComplexity = Churn * Complexity
					Score:           50 * (100 - 80), // Score = ChurnComplexity * (PerfectCoverage - Coverage)
				},
			},
		},
		{
			name: "calculate with coverage equal to perfect coverage",
			input: []*FileScore{
				{File: "test.go", Complexity: 10, Coverage: 100, Churn: 5},
			},
			opts: Options{
				PerfectCoverage: 100,
			},
			expected: []*FileScore{
				{
					File:            "test.go",
					Complexity:      10,
					Coverage:        100,
					Churn:           5,
					ChurnComplexity: 50, // ChurnComplexity = Churn * Complexity
					Score:           50, // When Coverage >= PerfectCoverage, Score = ChurnComplexity
				},
			},
		},
		{
			name: "calculate with coverage greater than perfect coverage",
			input: []*FileScore{
				{File: "test.go", Complexity: 10, Coverage: 120, Churn: 5},
			},
			opts: Options{
				PerfectCoverage: 100,
			},
			expected: []*FileScore{
				{
					File:            "test.go",
					Complexity:      10,
					Coverage:        120,
					Churn:           5,
					ChurnComplexity: 50, // ChurnComplexity = Churn * Complexity
					Score:           50, // When Coverage >= PerfectCoverage, Score = ChurnComplexity
				},
			},
		},
		{
			name: "calculate with different perfect coverage value",
			input: []*FileScore{
				{File: "test.go", Complexity: 10, Coverage: 70, Churn: 5},
			},
			opts: Options{
				PerfectCoverage: 80,
			},
			expected: []*FileScore{
				{
					File:            "test.go",
					Complexity:      10,
					Coverage:        70,
					Churn:           5,
					ChurnComplexity: 50,             // ChurnComplexity = Churn * Complexity
					Score:           50 * (80 - 70), // Score = ChurnComplexity * (PerfectCoverage - Coverage)
				},
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

func TestCombineMetrics(t *testing.T) {
	// Define test data
	churnData := []*git.ChurnChunk{
		{File: "file1.go", Churn: 100},
		{File: filepath.Join(".", "path", "to", "file3.go"), Churn: 200},
	}

	complexityData := []*complexity.FileStat{
		{Path: "file1.go", AvgComplexity: 10.0},
		{Path: filepath.Join("path", "to", "file3.go"), AvgComplexity: 5.0},
		{Path: "file4.go", AvgComplexity: 15.0},
	}

	coverageData := []*coverage.FileCoverage{
		{File: "file1.go", Coverage: 80.0},
		{File: "file2.go", Coverage: 90.0},
		{File: filepath.Join("path", "to", "file3.go"), Coverage: 70.0},
		{File: "file5.go", Coverage: 60.0},
	}

	// Expected results (sorted by file name for consistent comparison)
	expected := []*FileScore{
		{File: "file1.go", Churn: 100, Complexity: 10.0, Coverage: 80.0, ChurnComplexity: 1000.0},
		{File: filepath.Join("path", "to", "file3.go"), Churn: 200, Complexity: 5.0, Coverage: 70.0, ChurnComplexity: 1000.0},
	}

	result := CombineMetrics(churnData, complexityData, coverageData)

	assert.Equal(t, len(expected), len(result), "Result length mismatch")
	assert.ElementsMatch(t, result, expected)
}

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name            string
		input           *FileScore
		perfectCoverage float64
		expected        *FileScore
	}{
		{
			name: "churn is zero",
			input: &FileScore{
				File:       "test1.go",
				Churn:      0,
				Complexity: 10.0,
				Coverage:   75.0,
			},
			perfectCoverage: 100.0,
			expected: &FileScore{
				File:            "test1.go",
				Churn:           0,
				Complexity:      10.0,
				Coverage:        75.0,
				ChurnComplexity: 10.0,                  // When Churn is 0, ChurnComplexity equals Complexity
				Score:           10.0 * (100.0 - 75.0), // Score = ChurnComplexity * (perfectCoverage - Coverage)
			},
		},
		{
			name: "complexity is zero",
			input: &FileScore{
				File:       "test2.go",
				Churn:      20.0,
				Complexity: 0,
				Coverage:   80.0,
			},
			perfectCoverage: 100.0,
			expected: &FileScore{
				File:            "test2.go",
				Churn:           20.0,
				Complexity:      0,
				Coverage:        80.0,
				ChurnComplexity: 20.0,                  // When Complexity is 0, ChurnComplexity equals Churn
				Score:           20.0 * (100.0 - 80.0), // Score = ChurnComplexity * (perfectCoverage - Coverage)
			},
		},
		{
			name: "both churn and complexity non-zero",
			input: &FileScore{
				File:       "test3.go",
				Churn:      5.0,
				Complexity: 8.0,
				Coverage:   60.0,
			},
			perfectCoverage: 100.0,
			expected: &FileScore{
				File:            "test3.go",
				Churn:           5.0,
				Complexity:      8.0,
				Coverage:        60.0,
				ChurnComplexity: 5.0 * 8.0,                    // ChurnComplexity = Churn * Complexity
				Score:           (5.0 * 8.0) * (100.0 - 60.0), // Score = ChurnComplexity * (perfectCoverage - Coverage)
			},
		},
		{
			name: "coverage equals perfect coverage",
			input: &FileScore{
				File:       "test4.go",
				Churn:      15.0,
				Complexity: 5.0,
				Coverage:   100.0,
			},
			perfectCoverage: 100.0,
			expected: &FileScore{
				File:            "test4.go",
				Churn:           15.0,
				Complexity:      5.0,
				Coverage:        100.0,
				ChurnComplexity: 15.0 * 5.0, // ChurnComplexity = Churn * Complexity
				Score:           15.0 * 5.0, // When Coverage >= perfectCoverage, Score = ChurnComplexity
			},
		},
		{
			name: "coverage exceeds perfect coverage",
			input: &FileScore{
				File:       "test5.go",
				Churn:      10.0,
				Complexity: 6.0,
				Coverage:   110.0,
			},
			perfectCoverage: 100.0,
			expected: &FileScore{
				File:            "test5.go",
				Churn:           10.0,
				Complexity:      6.0,
				Coverage:        110.0,
				ChurnComplexity: 10.0 * 6.0, // ChurnComplexity = Churn * Complexity
				Score:           10.0 * 6.0, // When Coverage > perfectCoverage, Score = ChurnComplexity
			},
		},
		{
			name: "different perfect coverage value",
			input: &FileScore{
				File:       "test6.go",
				Churn:      7.0,
				Complexity: 9.0,
				Coverage:   70.0,
			},
			perfectCoverage: 80.0,
			expected: &FileScore{
				File:            "test6.go",
				Churn:           7.0,
				Complexity:      9.0,
				Coverage:        70.0,
				ChurnComplexity: 7.0 * 9.0,                   // ChurnComplexity = Churn * Complexity
				Score:           (7.0 * 9.0) * (80.0 - 70.0), // Score = ChurnComplexity * (perfectCoverage - Coverage)
			},
		},
		{
			name: "all values zero",
			input: &FileScore{
				File:       "test7.go",
				Churn:      0,
				Complexity: 0,
				Coverage:   0,
			},
			perfectCoverage: 100.0,
			expected: &FileScore{
				File:            "test7.go",
				Churn:           0,
				Complexity:      0,
				Coverage:        0,
				ChurnComplexity: 0,         // When Churn is 0, ChurnComplexity equals Complexity (which is also 0)
				Score:           0 * 100.0, // Score = ChurnComplexity * (perfectCoverage - Coverage)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculateScore(tt.input, tt.perfectCoverage)

			// Compare the actual result with the expected result
			assert.InDelta(t, tt.expected.ChurnComplexity, tt.input.ChurnComplexity, 0.001, "ChurnComplexity mismatch")
			assert.InDelta(t, tt.expected.Score, tt.input.Score, 0.001, "Score mismatch")
		})
	}
}

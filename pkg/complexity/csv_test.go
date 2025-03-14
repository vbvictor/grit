package complexity

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadComplexityFromCSV(t *testing.T) {
	tests := []struct {
		name    string
		csv     string
		want    []*FunctionStat
		wantErr bool
	}{
		{
			name: "valid csv with all fields",
			csv: `filename.go,Calculate,120,15,42,pkg1;pkg2
anothefile.go,Process,80,8,123,main`,
			want: []*FunctionStat{
				{
					File:       "filename.go",
					Name:       "Calculate",
					Length:     120,
					Complexity: 15,
					Line:       42,
					Package:    []string{"pkg1", "pkg2"},
				},
				{
					File:       "anothefile.go",
					Name:       "Process",
					Length:     80,
					Complexity: 8,
					Line:       123,
					Package:    []string{"main"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid csv with header row",
			csv:  `filename.go,Calculate,120,15,42,pkg1;pkg2`,
			want: []*FunctionStat{
				{
					File:       "filename.go",
					Name:       "Calculate",
					Length:     120,
					Complexity: 15,
					Line:       42,
					Package:    []string{"pkg1", "pkg2"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid csv with minimum required fields",
			csv: `filename.go,Calculate,120,15
anothefile.go,Process,80,8`,
			want: []*FunctionStat{
				{
					File:       "filename.go",
					Name:       "Calculate",
					Length:     120,
					Complexity: 15,
				},
				{
					File:       "anothefile.go",
					Name:       "Process",
					Length:     80,
					Complexity: 8,
				},
			},
			wantErr: false,
		},
		{
			name: "valid csv with mixed field counts",
			csv: `filename.go,Calculate,120,15,42
anothefile.go,Process,80,8,,main`,
			want: []*FunctionStat{
				{
					File:       "filename.go",
					Name:       "Calculate",
					Length:     120,
					Complexity: 15,
					Line:       42,
				},
				{
					File:       "anothefile.go",
					Name:       "Process",
					Length:     80,
					Complexity: 8,
					Package:    []string{"main"},
				},
			},
			wantErr: false,
		},
		{
			name:    "empty csv",
			csv:     "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "insufficient columns",
			csv:     "filename.go,Calculate,120",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid length value",
			csv:     "filename.go,Calculate,invalid,15",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid complexity value",
			csv:     "filename.go,Calculate,120,invalid",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid line value",
			csv:     "filename.go,Calculate,120,15,invalid",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.csv)
			got, err := readComplexityFromCSV(reader)

			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRunCSV(t *testing.T) {
	csvContent := `file1.go,func1,50,5,10,pkg1
file1.go,func2,70,10,20,pkg1
file2.go,func3,30,3,15,pkg2
file3.go,func4,100,15,25,pkg3;pkg4`

	csvPath := filepath.Join(t.TempDir(), "complexity.csv")
	err := os.WriteFile(csvPath, []byte(csvContent), 0o600)

	require.NoError(t, err)

	opts := &Options{
		Engine: CSV,
		Top:    10,
	}

	results, err := RunCSV(csvPath, opts)
	require.NoError(t, err)
	assert.Len(t, results, 3)

	expectedFiles := map[string]struct {
		FunctionCount int
		AvgComplexity float64
	}{
		"file1.go": {FunctionCount: 2, AvgComplexity: 7.5}, // (5+10)/2 = 7.5
		"file2.go": {FunctionCount: 1, AvgComplexity: 3.0},
		"file3.go": {FunctionCount: 1, AvgComplexity: 15.0},
	}

	for _, file := range results {
		expected, exists := expectedFiles[file.Path]
		assert.True(t, exists, "Unexpected file in results: %s", file.Path)

		if exists {
			assert.Len(t, file.Functions, expected.FunctionCount, "Incorrect function count for %s", file.Path)
			assert.InEpsilon(t, expected.AvgComplexity, file.AvgComplexity, 0.001,
				"Incorrect average complexity for %s", file.Path)
		}
	}

	// Test with exclude regex
	excludeOpts := &Options{
		Engine:       CSV,
		ExcludeRegex: regexp.MustCompile(`file[12]\.go`),
	}

	filteredResults, err := RunCSV(csvPath, excludeOpts)
	require.NoError(t, err)
	assert.Len(t, filteredResults, 1)
	assert.Equal(t, "file3.go", filteredResults[0].Path)
}

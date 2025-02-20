package plot

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupByFile(t *testing.T) {
	tests := []struct {
		name    string
		entries []ScatterEntry
		want    []groupedEntry
	}{
		{
			name:    "empty entries returns empty result",
			entries: []ScatterEntry{},
			want:    []groupedEntry{},
		},
		{
			name: "single file entry returns single group",
			entries: []ScatterEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file1.go",
				},
			},
			want: []groupedEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					Files:       []string{"file1.go"},
				},
			},
		},
		{
			name: "multiple files with same metrics group together",
			entries: []ScatterEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file1.go",
				},
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file2.go",
				},
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file3.go",
				},
			},
			want: []groupedEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					Files:       []string{"file1.go", "file2.go", "file3.go"},
				},
			},
		},
		{
			name: "multiple files multiple metrics group",
			entries: []ScatterEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file1.go",
				},
				{
					ScatterData: ScatterData{Complexity: 9.0, Churn: 3},
					File:        "file2.go",
				},
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file3.go",
				},
				{
					ScatterData: ScatterData{Complexity: 7.0, Churn: 3},
					File:        "file4.go",
				},
				{
					ScatterData: ScatterData{Complexity: 9.0, Churn: 5},
					File:        "file5.go",
				},
			},
			want: []groupedEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					Files:       []string{"file1.go", "file3.go"},
				},
				{
					ScatterData: ScatterData{Complexity: 9.0, Churn: 3},
					Files:       []string{"file2.go"},
				},
				{
					ScatterData: ScatterData{Complexity: 7.0, Churn: 3},
					Files:       []string{"file4.go"},
				},
				{
					ScatterData: ScatterData{Complexity: 9.0, Churn: 5},
					Files:       []string{"file5.go"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := groupByFile(tt.entries)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

type stubMapper struct{}

func (m *stubMapper) Map(data ScatterData) Category {
	if data.Complexity >= 10.0 {
		return "critical"
	}

	if data.Complexity >= 5.0 {
		return "warning"
	}

	return "normal"
}

func (m *stubMapper) Style(_ Category) opts.ItemStyle {
	return opts.ItemStyle{}
}

func TestFormDataSeries(t *testing.T) {
	tests := []struct {
		name    string
		entries []ScatterEntry
		want    ScatterSeries
	}{
		{
			name:    "empty entries returns empty series",
			entries: []ScatterEntry{},
			want:    ScatterSeries{},
		},
		{
			name: "entries are grouped by same metrics and mapped to categories",
			entries: []ScatterEntry{
				{
					ScatterData: ScatterData{Complexity: 12.0, Churn: 5},
					File:        "critical1.go",
				},
				{
					ScatterData: ScatterData{Complexity: 12.0, Churn: 5},
					File:        "critical2.go",
				},
				{
					ScatterData: ScatterData{Complexity: 7.0, Churn: 3},
					File:        "warning.go",
				},
				{
					ScatterData: ScatterData{Complexity: 3.0, Churn: 1},
					File:        "normal1.go",
				},
				{
					ScatterData: ScatterData{Complexity: 1.0, Churn: 1},
					File:        "normal2.go",
				},
				{
					ScatterData: ScatterData{Complexity: 2.0, Churn: 1},
					File:        "normal3.go",
				},
			},
			want: ScatterSeries{
				"critical": []opts.ScatterData{
					{
						Value:      []interface{}{12.0, 5, "critical1.go<br/>critical2.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
				},
				"warning": []opts.ScatterData{
					{
						Value:      []interface{}{7.0, 3, "warning.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
				},
				"normal": []opts.ScatterData{
					{
						Value:      []interface{}{3.0, 1, "normal1.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
					{
						Value:      []interface{}{1.0, 1, "normal2.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
					{
						Value:      []interface{}{2.0, 1, "normal3.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
				},
			},
		},
	}

	mapper := &stubMapper{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formDataSeries(tt.entries, mapper)

			for category, series := range tt.want {
				assert.ElementsMatch(t, series, got[category])
			}
		})
	}
}

func readCSVChartEntries(filepath string) ([]ScatterEntry, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv file: %w", err)
	}

	entries := make([]ScatterEntry, 0, len(records))

	for _, record := range records {
		currComp, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse complexity: %w", err)
		}

		churn, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse churn: %w", err)
		}

		entry := ScatterEntry{
			File:        record[0],
			ScatterData: ScatterData{Complexity: currComp, Churn: int(churn)},
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func createTestChart(t *testing.T, entries []ScatterEntry, outputPath string) {
	t.Helper()

	err := CreateScatterChart(entries, NewRisksMapper(), outputPath)
	require.NoError(t, err)

	_, err = os.Stat(outputPath)
	assert.NoError(t, err)
}

var (
	CSVDataDir = "../../test/data/"
	OutputDir  = "../../test/charts/"
)

func TestCreateScatterCharts(t *testing.T) {
	err := os.MkdirAll(OutputDir, 0o755)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		csvFile string
		outFile string
	}{
		{
			name:    "200 different entries",
			csvFile: "plot_200.csv",
			outFile: "scatter-200.html",
		},
		{
			name:    "2000 different entries",
			csvFile: "plot_2000.csv",
			outFile: "scatter-2000.html",
		},
		{
			name:    "10 same entries",
			csvFile: "plot_10-same.csv",
			outFile: "scatter-10-same.html",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entries, err := readCSVChartEntries(CSVDataDir + tc.csvFile)
			require.NoError(t, err)

			createTestChart(t, entries, OutputDir+tc.outFile)
		})
	}
}

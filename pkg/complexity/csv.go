package complexity

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func RunCSV(filepath string, opts *Options) ([]*FileStat, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file at %s: %w", filepath, err)
	}
	defer file.Close()

	functionStats, err := readComplexityFromCSV(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse complexity data from CSV: %w", err)
	}

	fileMap := make(map[string][]FunctionStat)
	for _, stat := range functionStats {
		fileMap[stat.File] = append(fileMap[stat.File], *stat)
	}

	result := make([]*FileStat, 0, len(fileMap))

	for file, functions := range fileMap {
		if opts.ExcludeRegex != nil && opts.ExcludeRegex.MatchString(file) {
			continue
		}

		result = append(result, &FileStat{
			Path:      file,
			Functions: functions,
		})
	}

	// Calculate average complexity for each file
	AvgComplexity(result)

	return result, nil
}

const minimalColumns = 4

func readComplexityFromCSV(r io.Reader) ([]*FunctionStat, error) { //nolint:cyclop // complexity is not a problem here
	csvReader := csv.NewReader(r)
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields per record

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV data: %w", err)
	}

	if len(records) == 0 {
		return nil, errors.New("CSV data is empty")
	}

	result := make([]*FunctionStat, 0, len(records))

	for pos, record := range records {
		// Minimum required fields: filename, function, length, complexity
		if len(record) < minimalColumns {
			return nil, fmt.Errorf("row %d: insufficient columns, expected at least 4", pos+1)
		}

		stat := &FunctionStat{
			File: record[0],
			Name: record[1],
		}

		length, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid length value '%s': %w", pos+1, record[2], err)
		}

		stat.Length = length

		complexity, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid complexity value '%s': %w", pos+1, record[3], err)
		}

		stat.Complexity = complexity

		// Optional: line number
		if len(record) > 4 && record[4] != "" {
			line, err := strconv.Atoi(record[4])
			if err != nil {
				return nil, fmt.Errorf("row %d: invalid line value '%s': %w", pos+1, record[4], err)
			}

			stat.Line = line
		}

		// Optional: packages
		if len(record) > 5 && record[5] != "" {
			stat.Package = strings.Split(record[5], ";")
		}

		result = append(result, stat)
	}

	return result, nil
}

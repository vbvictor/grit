package coverage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/vbvictor/grit/grit/cmd/flag"
	"golang.org/x/tools/cover"
)

var errUnsupportedMode = errors.New("unsupported coverage mode")

const (
	percentMultiplier = 100.0
	defaultTop        = 10
)

type SortType = string

var (
	Worst SortType = "worst"
	Best  SortType = "best"
)

type FileCoverage struct {
	File       string
	Coverage   float64
	Statements int
	Covered    int
}

type Options struct {
	SortBy      SortType
	Top         int
	ExcludePath string
}

func ReadCoverage(coverageFile string, opts Options) ([]*FileCoverage, error) {
	profiles, err := cover.ParseProfiles(coverageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse profiles data: %w", err)
	}

	results := make([]*FileCoverage, 0)

	for _, profile := range profiles {
		total := 0
		covered := 0

		if profile.Mode != "set" {
			return nil, fmt.Errorf("%w: %s", errUnsupportedMode, profile.Mode)
		}

		for _, block := range profile.Blocks {
			total += block.NumStmt

			if block.Count > 0 {
				covered += block.NumStmt
			}
		}

		coverage := 0.0
		if total > 0 {
			coverage = float64(covered) * percentMultiplier / float64(total)
		}

		results = append(results, &FileCoverage{
			File:       profile.FileName,
			Coverage:   coverage,
			Statements: total,
			Covered:    covered,
		})
	}

	return results, nil
}

func SortAndLimit(result []*FileCoverage, sortBy SortType, limit int) []*FileCoverage {
	less := func() func(i, j int) bool {
		switch sortBy {
		case Worst:
			return func(i, j int) bool { return result[i].Coverage < result[j].Coverage }
		case Best:
			return func(i, j int) bool { return result[i].Coverage > result[j].Coverage }
		default:
			return nil
		}
	}()

	sort.Slice(result, less)

	if limit >= 0 && len(result) > limit {
		return result[:limit]
	}

	return result
}

func RunCoverage(repoPath, coverageFile string) error {
	coverageArg := filepath.Join(repoPath, "...")
	cmd := exec.Command("go", "test", coverageArg, "-coverprofile="+coverageFile)

	flag.LogIfVerbose("Running command: %s\n", cmd.String())

	var stderr bytes.Buffer

	if flag.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	} else {
		cmd.Stderr = &stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run tests: %w\nstderr: %s", err, stderr.String())
	}

	return nil
}

// sortByCoverage sorts FileCoverage slice by Coverage field
// Ascending order if asc is true, descending if false
func sortByCoverage(files []FileCoverage, asc bool) []FileCoverage {
	sorted := make([]FileCoverage, len(files))
	copy(sorted, files)

	sort.Slice(sorted, func(i, j int) bool {
		if asc {
			return sorted[i].Coverage < sorted[j].Coverage
		}

		return sorted[i].Coverage > sorted[j].Coverage
	})

	return sorted
}

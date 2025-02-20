package coverage

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"

	"golang.org/x/tools/cover"
)

var errUnsupportedMode = errors.New("unsupported coverage mode")

const (
	percentMultiplier = 100.0
	defaultTop        = 10
)

type SortType = string

var (
	CoverageAsc  SortType = "asc"
	CoverageDesc SortType = "desc"
)

type Options struct {
	SortBy      SortType
	Top         int
	ExcludePath string
}

var CoverageOpts = Options{
	SortBy:      CoverageAsc,
	Top:         defaultTop,
	ExcludePath: "",
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

	sortAndLimit(results, opts.SortBy, opts.Top)

	return results, nil
}

func sortAndLimit(result []*FileCoverage, sortBy SortType, limit int) {
	less := func() func(i, j int) bool {
		switch sortBy {
		case CoverageAsc:
			return func(i, j int) bool { return result[i].Coverage > result[j].Coverage }
		case CoverageDesc:
			return func(i, j int) bool { return result[i].Coverage > result[j].Coverage }
		default:
			return nil
		}
	}()

	sort.Slice(result, less)

	if limit >= 0 && len(result) > limit {
		result = result[:limit]
	}
}

func RunCoverage(repoPath, coverageFile string) error {
	coverageArg := filepath.Join(repoPath, "...")
	cmd := exec.Command("go", "test", coverageArg, "-coverprofile="+coverageFile)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run tests: %w\nstderr: %s", err, stderr.String())
	}

	return nil
}

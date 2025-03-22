package coverage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

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
	SortBy           SortType
	Top              int
	ExcludeRegex     *regexp.Regexp
	RunCoverage      string
	CoverageFilename string
	OutputFormat     string
}

func PopulateOpts(opts *Options, excludeRegex string) error {
	if excludeRegex != "" {
		var err error

		opts.ExcludeRegex, err = regexp.Compile(excludeRegex)
		if err != nil {
			return fmt.Errorf("invalid exclude pattern: %w", err)
		}
	}

	return nil
}

func GetCoverageData(repoPath string, coverageOpts *Options) ([]*FileCoverage, error) {
	coveragePath := filepath.Join(repoPath, coverageOpts.CoverageFilename)

	_, err := os.Stat(coveragePath)
	if os.IsNotExist(err) {
		flag.LogIfVerbose("Coverage file %s not found\n", coveragePath)

		if coverageOpts.RunCoverage != flag.Never {
			flag.LogIfVerbose("Running test suite\n\n")

			if err = RunCoverage(repoPath, coverageOpts.CoverageFilename); err != nil {
				return nil, errors.Join(flag.ErrRunCoverage, err)
			}

			flag.LogIfVerbose("Coverage file %s created\n", coveragePath)
		} else {
			flag.LogIfVerbose("Go test did not run since flag run='Never'\n")

			return nil, errors.Join(flag.ErrCoverageNotFound, err)
		}
	} else if coverageOpts.RunCoverage == flag.Always {
		flag.LogIfVerbose("Removing previous test coverage file %s\n", coveragePath)
		os.Remove(coveragePath)
		flag.LogIfVerbose("Running test suite\n\n")

		if err = RunCoverage(repoPath, coverageOpts.CoverageFilename); err != nil {
			return nil, errors.Join(flag.ErrRunCoverage, err)
		}

		flag.LogIfVerbose("Coverage file %s created\n", coveragePath)
	}

	covData, err := ReadCoverage(repoPath, coverageOpts.CoverageFilename, coverageOpts)
	if err != nil {
		return nil, errors.Join(flag.ErrReadCoverage, err)
	}

	return covData, nil
}

func ReadCoverage(path, file string, opts *Options) ([]*FileCoverage, error) {
	profiles, err := cover.ParseProfiles(filepath.Join(path, file))
	if err != nil {
		return nil, fmt.Errorf("failed to parse profiles data: %w", err)
	}

	results := make([]*FileCoverage, 0)

	for _, profile := range profiles {
		if opts.ExcludeRegex != nil && opts.ExcludeRegex.MatchString(profile.FileName) {
			continue
		}

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
			File:       extractRelativePath(profile.FileName),
			Coverage:   coverage,
			Statements: total,
			Covered:    covered,
		})
	}

	return results, nil
}

const minPathPaths = 3

func extractRelativePath(fullPath string) string {
	// Convert to consistent path format
	fullPath = filepath.FromSlash(fullPath)

	// Split the path by separator
	parts := strings.Split(fullPath, string(os.PathSeparator))

	// If we have at least 3 components (typically github.com/username/module/...)
	if len(parts) > minPathPaths {
		// Skip the first two components (domain and username)
		return filepath.Join(parts[minPathPaths:]...)
	}

	// Fallback to original path if it doesn't have enough components
	return fullPath
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

	if limit > 0 && len(result) > limit {
		return result[:limit]
	}

	return result
}

func RunCoverage(repoPath, coverageFile string) error {
	cmd := exec.Command( //nolint:gosec // go test is allowed command
		"go",
		"test",
		"./...",
		"-v",
		"-coverprofile="+coverageFile)
	cmd.Dir = repoPath

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

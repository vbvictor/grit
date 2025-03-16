package complexity

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
)

type Engine = string

const (
	Gocyclo  = "gocyclo"
	Gocognit = "gocognit"
	CSV      = "csv-file"
)

type FileStat struct {
	Path          string
	Functions     []FunctionStat
	AvgComplexity float64
}

type FunctionStat struct {
	File       string
	Package    []string
	Name       string
	Line       int
	Length     int
	Complexity int
}

type Options struct {
	Engine       string
	ExcludeRegex *regexp.Regexp
	Top          int
	OutputFormat string
}

var ErrUnsupportedEngine = errors.New("unsupported complexity engine")

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

func RunComplexity(repoPath string, opts *Options) ([]*FileStat, error) {
	switch opts.Engine {
	case Gocyclo:
		return RunGocyclo(repoPath, opts)
	case Gocognit:
		return RunGocognit(repoPath, opts)
	case CSV:
		return RunCSV(filepath.Join(repoPath, "complexity.csv"), opts)
	default:
		return nil, ErrUnsupportedEngine
	}
}

func SortAndLimit(fileStat []*FileStat, opts Options) []*FileStat {
	slices.SortFunc(fileStat, func(a, b *FileStat) int {
		if b.AvgComplexity > a.AvgComplexity {
			return 1
		}

		if b.AvgComplexity < a.AvgComplexity {
			return -1
		}

		return 0
	})

	if opts.Top > 0 && opts.Top < len(fileStat) {
		return fileStat[:opts.Top]
	}

	return fileStat
}

func AvgComplexity(files []*FileStat) {
	for _, file := range files {
		if len(file.Functions) == 0 {
			continue
		}

		fileComplexity := 0.0
		for _, fn := range file.Functions {
			fileComplexity += float64(fn.Complexity)
		}

		complexity := fileComplexity / float64(len(file.Functions))
		file.AvgComplexity = complexity
	}
}

type FilesFilterFunc func(files []FileStat) []FileStat

func ApplyFilters(files []FileStat, filters ...FilesFilterFunc) []FileStat {
	result := files

	for _, filter := range filters {
		result = filter(result)
	}

	return result
}

type MinComplexityFilter struct {
	MinComplexity int
}

const (
	MinComplexityDefault = 5
)

func (f MinComplexityFilter) Filter(files []FileStat) []FileStat {
	result := make([]FileStat, 0, len(files))

	for _, file := range files {
		filteredFuncs := make([]FunctionStat, 0)

		for _, fn := range file.Functions {
			if fn.Complexity >= f.MinComplexity {
				filteredFuncs = append(filteredFuncs, fn)
			}
		}

		if len(filteredFuncs) > 0 {
			result = append(result, FileStat{
				Path:      file.Path,
				Functions: filteredFuncs,
			})
		}
	}

	return result
}

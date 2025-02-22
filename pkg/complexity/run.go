package complexity

import (
	"fmt"
	"slices"
)

type Engine = string

const (
	Gocyclo  = "gocyclo"
	Gocognit = "gocognit"
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
	Engine      string
	ExcludePath string
	Top         int
}

func RunComplexity(repoPath string, opts Options) ([]*FileStat, error) {
	switch opts.Engine {
	case Gocyclo:
		return RunGocyclo(repoPath, opts)
	case Gocognit:
		return RunGocognit(repoPath, opts)
	default:
		return nil, fmt.Errorf("unsupported complexity engine: %s", opts.Engine)
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

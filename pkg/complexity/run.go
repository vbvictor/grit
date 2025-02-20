package complexity

import (
	"fmt"
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

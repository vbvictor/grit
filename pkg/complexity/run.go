package complexity

import (
	"fmt"
)

type Engine = string

const (
	Gocyclo  = "gocyclo"
	Gocognit = "gocognit"
)

type Options struct {
	Engine  string
	Exclude string
}

var Opts = Options{
	Engine:  Gocyclo,
	Exclude: "",
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

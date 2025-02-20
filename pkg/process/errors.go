package process

import "errors"

type ErrAbsRepoPath struct {
	Path string
}

func (e *ErrAbsRepoPath) Error() string {
	return "failed to get absolute path from " + e.Path
}

var (
	ErrCoverageNotFound = errors.New("failed to find file with code coverage")
	ErrReadCoverage     = errors.New("failed to read coverage file")
	ErrRunCoverage      = errors.New("failed to run coverage")
)

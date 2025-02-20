package complexity

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/fzipp/gocyclo"
)

func RunGocyclo(repoPath string, opts Options) ([]*FileStat, error) {
	var excludeRegex *regexp.Regexp

	if opts.ExcludePath != "" {
		var err error

		excludeRegex, err = regexp.Compile(opts.ExcludePath)
		if err != nil {
			return nil, fmt.Errorf("invalid exclude pattern: %w", err)
		}
	}

	paths := []string{repoPath}
	stats := gocyclo.Analyze(paths, excludeRegex)

	result := make([]*FileStat, 0)
	fileMap := make(map[string][]FunctionStat)

	for _, stat := range stats {
		relPath, err := filepath.Rel(repoPath, stat.Pos.Filename)
		if err != nil {
			return nil, err
		}

		funcStat := FunctionStat{
			File:       relPath,
			Package:    []string{stat.PkgName},
			Name:       stat.FuncName,
			Line:       stat.Pos.Line,
			Complexity: stat.Complexity,
		}

		fileMap[relPath] = append(fileMap[relPath], funcStat)
	}

	for filePath, functions := range fileMap {
		result = append(result, &FileStat{
			Path:      filePath,
			Functions: functions,
		})
	}

	AvgComplexity(result)

	return result, nil
}

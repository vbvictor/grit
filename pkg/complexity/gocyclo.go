package complexity

import (
	"fmt"
	"path/filepath"

	"github.com/fzipp/gocyclo"
)

func RunGocyclo(repoPath string, opts *Options) ([]*FileStat, error) {
	paths := []string{repoPath}
	stats := gocyclo.Analyze(paths, opts.ExcludeRegex)

	result := make([]*FileStat, 0)
	fileMap := make(map[string][]FunctionStat)

	for _, stat := range stats {
		relPath, err := filepath.Rel(repoPath, stat.Pos.Filename)
		if err != nil {
			return nil, fmt.Errorf("failed to get relative path: %w", err)
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
			Path:      filepath.ToSlash(filepath.Clean(filePath)),
			Functions: functions,
		})
	}

	AvgComplexity(result)

	return result, nil
}

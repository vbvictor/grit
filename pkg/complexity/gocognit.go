package complexity

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/uudashr/gocognit"
)

func RunGocognit(repoPath string, opts *Options) ([]*FileStat, error) { //nolint:funlen,cyclop // TODO: refactor later
	fileMap := make(map[string][]FunctionStat)

	err := filepath.Walk(repoPath, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path: %w", err)
		}

		if opts.ExcludeRegex != nil && opts.ExcludeRegex.MatchString(path) {
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fileSet := token.NewFileSet()

		file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		stats := gocognit.ComplexityStats(file, fileSet, nil)
		functions := make([]FunctionStat, 0, len(stats))

		for _, stat := range stats {
			relPath, err := filepath.Rel(repoPath, stat.Pos.Filename)
			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}

			functions = append(functions, FunctionStat{
				File:       relPath,
				Package:    []string{stat.PkgName},
				Name:       stat.FuncName,
				Line:       stat.Pos.Line,
				Complexity: stat.Complexity,
			})
		}

		if len(functions) > 0 {
			relPath, err := filepath.Rel(repoPath, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path for file: %w", err)
			}

			fileMap[relPath] = functions
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk repository: %w", err)
	}

	result := make([]*FileStat, 0)

	for filePath, functions := range fileMap {
		result = append(result, &FileStat{
			Path:      filePath,
			Functions: functions,
		})
	}

	AvgComplexity(result)

	return result, nil
}

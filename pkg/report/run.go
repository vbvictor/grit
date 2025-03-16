package report

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/vbvictor/grit/pkg/complexity"
	"github.com/vbvictor/grit/pkg/coverage"
	"github.com/vbvictor/grit/pkg/git"
	"golang.org/x/exp/maps"
)

type FileScore struct {
	File            string
	Coverage        float64
	Churn           float64
	Complexity      float64
	ChurnComplexity float64
	Score           float64
}

type Options struct {
	ChurnFactor      float64
	ComplexityFactor float64
	CoverageFactor   float64
	PerfectCoverage  float64
	Top              int
	ExcludePath      string
}

func CalculateScores(data []*FileScore, opts Options) []*FileScore {
	for _, file := range data {
		calculateScore(file, opts.PerfectCoverage)
	}

	return data
}

func calculateScore(fileScore *FileScore, perfectCoverage float64) {
	switch {
	case fileScore.Churn == 0:
		fileScore.ChurnComplexity = fileScore.Complexity
	case fileScore.Complexity == 0:
		fileScore.ChurnComplexity = fileScore.Churn
	default:
		fileScore.ChurnComplexity = fileScore.Churn * fileScore.Complexity
	}

	if fileScore.Coverage >= perfectCoverage {
		fileScore.Score = fileScore.ChurnComplexity
	} else {
		fileScore.Score = fileScore.ChurnComplexity * (perfectCoverage - fileScore.Coverage)
	}
}

func SortByScore(files []*FileScore) []*FileScore {
	slices.SortFunc(files, func(lhs *FileScore, rhs *FileScore) int {
		if lhs.Score < rhs.Score {
			return 1
		}

		if lhs.Score > rhs.Score {
			return -1
		}

		return 0
	})

	return files
}

func CombineMetrics(
	churnData []*git.ChurnChunk,
	complexityData []*complexity.FileStat,
	coverageData []*coverage.FileCoverage,
) []*FileScore {
	fileMap := make(map[string]*FileScore)

	for _, chunk := range churnData {
		normalizedPath := normalizePath(chunk.File)
		score, exists := fileMap[normalizedPath]

		if !exists {
			score = &FileScore{
				File: chunk.File,
			}
			fileMap[normalizedPath] = score
		}

		score.Churn = float64(chunk.Churn)
	}

	for _, stat := range complexityData {
		normalizedPath := normalizePath(stat.Path)
		if _, exists := fileMap[normalizedPath]; exists {
			fileMap[normalizedPath].Complexity = stat.AvgComplexity
		}
	}

	for _, cov := range coverageData {
		normalizedPath := normalizePath(cov.File)
		if _, exists := fileMap[normalizedPath]; exists {
			fileMap[normalizedPath].Coverage = cov.Coverage
		}
	}

	for _, score := range fileMap {
		score.ChurnComplexity = score.Churn * score.Complexity
	}

	return maps.Values(fileMap)
}

func normalizePath(path string) string {
	normalized := filepath.ToSlash(path)
	normalized = strings.TrimPrefix(normalized, "./")
	normalized = strings.TrimPrefix(normalized, "/")

	return normalized
}

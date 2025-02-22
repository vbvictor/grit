package report

import "slices"

type FileScore struct {
	File       string
	Coverage   float64
	Complexity float64
	Churn      float64
	Score      float64
}

const defaultTop = 10

type Options struct {
	ChurnFactor      float64
	ComplexityFactor float64
	CoverageFactor   float64
	Top              int
	ExcludePath      string
}

var ReportOpts = Options{
	Top:              defaultTop,
	ExcludePath:      "",
	ChurnFactor:      1.0,
	ComplexityFactor: 1.0,
	CoverageFactor:   1.0,
}

type scoreFunc func(*FileScore)

func CalculateScores(data []*FileScore, opts Options) []*FileScore {
	var calculateScore scoreFunc

	switch {
	case opts.ChurnFactor == 0:
		calculateScore = func(fs *FileScore) {
			fs.Score = (fs.Complexity * opts.ComplexityFactor) / (fs.Coverage + opts.CoverageFactor)
		}
	case opts.ComplexityFactor == 0:
		calculateScore = func(fs *FileScore) {
			fs.Score = (fs.Churn * opts.ChurnFactor) / (fs.Coverage + opts.CoverageFactor)
		}
	case opts.CoverageFactor == 0:
		calculateScore = func(fs *FileScore) {
			fs.Score = (fs.Churn * opts.ChurnFactor) + (fs.Complexity * opts.ComplexityFactor)
		}
	default:
		calculateScore = func(fs *FileScore) {
			fs.Score = (fs.Churn * opts.ChurnFactor) + (fs.Complexity*opts.ComplexityFactor)/(fs.Coverage+opts.CoverageFactor)
		}
	}

	for _, file := range data {
		calculateScore(file)
	}

	return data
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

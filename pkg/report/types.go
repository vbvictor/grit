package report

type FileScore struct {
	File       string
	Package    string
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

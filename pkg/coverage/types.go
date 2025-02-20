package coverage

type FileCoverage struct {
	File       string
	Coverage   float64
	Statements int
	Covered    int
}

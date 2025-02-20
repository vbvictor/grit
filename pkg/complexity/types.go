package complexity

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

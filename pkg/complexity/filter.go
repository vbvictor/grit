package complexity

type FilesFilterFunc func(files []FileStat) []FileStat

func ApplyFilters(files []FileStat, filters ...FilesFilterFunc) []FileStat {
	result := files

	for _, filter := range filters {
		result = filter(result)
	}

	return result
}

type MinComplexityFilter struct {
	MinComplexity int
}

const (
	MinComplexityDefault = 5
)

func (f MinComplexityFilter) Filter(files []FileStat) []FileStat {
	result := make([]FileStat, 0, len(files))

	for _, file := range files {
		filteredFuncs := make([]FunctionStat, 0)

		for _, fn := range file.Functions {
			if fn.Complexity >= f.MinComplexity {
				filteredFuncs = append(filteredFuncs, fn)
			}
		}

		if len(filteredFuncs) > 0 {
			result = append(result, FileStat{
				Path:      file.Path,
				Functions: filteredFuncs,
			})
		}
	}

	return result
}

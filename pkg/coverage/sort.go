package coverage

import "sort"

// sortByCoverage sorts FileCoverage slice by Coverage field
// Ascending order if asc is true, descending if false
func sortByCoverage(files []FileCoverage, asc bool) []FileCoverage {
	sorted := make([]FileCoverage, len(files))
	copy(sorted, files)

	sort.Slice(sorted, func(i, j int) bool {
		if asc {
			return sorted[i].Coverage < sorted[j].Coverage
		}

		return sorted[i].Coverage > sorted[j].Coverage
	})

	return sorted
}

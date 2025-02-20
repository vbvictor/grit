package complexity

// AvgComplexity Calculates complexity bases on functions in file: sum(funcComplexity) / funcCount.
func AvgComplexity(files []*FileStat) {
	// fileComplexities := make([]FileComplexity, 0, len(files))

	for _, file := range files {
		if len(file.Functions) == 0 {
			continue
		}

		fileComplexity := 0.0
		for _, fn := range file.Functions {
			fileComplexity += float64(fn.Complexity)
		}

		complexity := fileComplexity / float64(len(file.Functions))
		file.AvgComplexity = complexity

		// fileComplexities = append(fileComplexities, FileComplexity{
		// 	File:       file.Path,
		// 	Complexity: complexity,
		// })
	}
}

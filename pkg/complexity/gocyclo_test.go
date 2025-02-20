package complexity

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunGocyclo(t *testing.T) {
	tests := []struct {
		name          string
		directory     string
		expectedFiles map[string]ExpectedFile
	}{
		{
			name:          "Empty directory",
			directory:     "empty",
			expectedFiles: map[string]ExpectedFile{},
		},
		{
			name:      "Nested directory structure",
			directory: "nested",
			expectedFiles: map[string]ExpectedFile{
				"main.go": {
					Functions: map[string]TestFunction{
						"BaseFunction":    {1, 3, "nested", "main.go"},
						"SimpleCondition": {2, 7, "nested", "main.go"},
					},
					AvgComplexity: 1.5,
				},
				filepath.Join("level1", "file1.go"): {
					Functions: map[string]TestFunction{
						"NestedIf":          {3, 3, "level1", filepath.Join("level1", "file1.go")},
						"LoopWithCondition": {4, 13, "level1", filepath.Join("level1", "file1.go")},
					},
					AvgComplexity: 3.5,
				},
				filepath.Join("level1", "level2", "file1.go"): {
					Functions: map[string]TestFunction{
						"Func1": {4, 3, "level2", filepath.Join("level1", "level2", "file1.go")},
						"Func2": {4, 19, "level2", filepath.Join("level1", "level2", "file1.go")},
					},
					AvgComplexity: 4.0,
				},
				filepath.Join("level1", "level2", "file2.go"): {
					Functions: map[string]TestFunction{
						"Func3": {4, 3, "level2", filepath.Join("level1", "level2", "file2.go")},
						"Func4": {5, 15, "level2", filepath.Join("level1", "level2", "file2.go")},
					},
					AvgComplexity: 4.5,
				},
				filepath.Join("level1", "level2", "file3.go"): {
					Functions: map[string]TestFunction{
						"NestedLoopsWithConditions": {5, 3, "level2", filepath.Join("level1", "level2", "file3.go")},
						"SwitchWithLoops":           {6, 17, "level2", filepath.Join("level1", "level2", "file3.go")},
					},
					AvgComplexity: 5.5,
				},
				filepath.Join("level1", "level2", "morelevel2", "file1.go"): {
					Functions: map[string]TestFunction{
						"Func5": {6, 3, "level2", filepath.Join("level1", "level2", "morelevel2", "file1.go")},
						"Func6": {4, 19, "level2", filepath.Join("level1", "level2", "morelevel2", "file1.go")},
					},
					AvgComplexity: 5.0,
				},
				filepath.Join("level1", "level2", "morelevel2", "file2.go"): {
					Functions: map[string]TestFunction{
						"Func7": {7, 3, "level2", filepath.Join("level1", "level2", "morelevel2", "file2.go")},
						"Func8": {7, 24, "level2", filepath.Join("level1", "level2", "morelevel2", "file2.go")},
					},
					AvgComplexity: 7.0,
				},
				filepath.Join("level1", "level2", "morelevel2", "file3.go"): {
					Functions: map[string]TestFunction{
						"ComplexNestedStructure": {8, 3, "level2", filepath.Join("level1", "level2", "morelevel2", "file3.go")},
						"MultipleControlFlows":   {8, 24, "level2", filepath.Join("level1", "level2", "morelevel2", "file3.go")},
					},
					AvgComplexity: 8.0,
				},
			},
		},
		{
			name:      "Mixed complexity functions",
			directory: "mixed",
			expectedFiles: map[string]ExpectedFile{
				"main.go": {
					Functions: map[string]TestFunction{
						"simpleFunction":  {1, 3, "mixed", "main.go"},
						"complexFunction": {8, 7, "mixed", "main.go"},
					},
					AvgComplexity: 4.5,
				},
			},
		},
		{
			name:      "Special cases",
			directory: "special",
			expectedFiles: map[string]ExpectedFile{
				"my main.go": {
					Functions: map[string]TestFunction{
						"simpleFunction":  {1, 3, "special", "my main.go"},
						"complexFunction": {8, 7, "special", "my main.go"},
					},
					AvgComplexity: 4.5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPath := filepath.Join(".", "..", "..", "test", "complexity", "gocode", tt.directory)
			result, err := RunGocyclo(testPath, Options{})

			require.NoError(t, err)
			assert.Len(t, result, len(tt.expectedFiles))

			if len(tt.expectedFiles) == 0 {
				return
			}

			for _, file := range result {
				expected, exists := tt.expectedFiles[file.Path]
				assert.True(t, exists, "File %s should exist", file.Path)

				if exists {
					// Check average complexity
					assert.InDelta(t, expected.AvgComplexity, file.AvgComplexity, 0.01,
						"File %s should have average complexity %f", file.Path, expected.AvgComplexity)

					// Check functions
					for _, fn := range file.Functions {
						expectedFn, fnExists := expected.Functions[fn.Name]
						assert.True(t, fnExists, "Function %s should exist", fn.Name)

						if fnExists {
							assert.Equal(t, expectedFn.complexity, fn.Complexity)
							assert.Equal(t, expectedFn.line, fn.Line)
							assert.Equal(t, []string{expectedFn.pkg}, fn.Package)
							assert.Equal(t, expectedFn.file, fn.File)
						}
					}
				}
			}
		})
	}
}

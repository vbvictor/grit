package test

import (
	"testing"
)

// To run these integration tests, you need to have following:
// 1) Grit installed in your system.
// 		You can install it via go install.
// 		You can download it from https://github.com/vbvictor/grit/releases and place it in your PATH.

func TestGritBasicFunctionality(t *testing.T) {
	targetDir := `C:\Users\Victor\repos\grit`

	tests := []GritTest{
		{
			Name:      "Basic run with default settings",
			RunDir:    targetDir,
			Args:      []string{},
			Validator: NoopValidator,
		},
		{
			Name:      "Run with CSV output format",
			RunDir:    targetDir,
			Args:      []string{"--format", "csv"},
			Validator: NoopValidator,
		},
		{
			Name:      "Run with specific engine",
			RunDir:    targetDir,
			Args:      []string{"--engine", "gocyclo"},
			Validator: NoopValidator,
		},
		{
			Name:      "Run with plot output",
			RunDir:    targetDir,
			Args:      []string{"--plot", "complexity_chart.html"},
			Validator: NoopValidator,
		},
		{
			Name:      "Run with output file",
			RunDir:    targetDir,
			Args:      []string{"--output", "analysis_output.txt"},
			Validator: NoopValidator,
		},
		{
			Name:      "Run with help flag",
			RunDir:    targetDir,
			Args:      []string{"--help"},
			Validator: ContainsValidator,
		},
	}

	RunGritTests(t, tests)
}

func TestGritAdvancedUsage(t *testing.T) {
	// Target directory to analyze
	targetDir := `C:\Users\Victor\repos\fasthttp`

	tests := []GritTest{
		{
			Name:      "Run with file pattern",
			RunDir:    targetDir,
			Args:      []string{"--include", "*.go"},
			Validator: NoopValidator,
		},
		{
			Name:      "Run with exclusion pattern",
			RunDir:    targetDir,
			Args:      []string{"--exclude", "vendor/"},
			Validator: NoopValidator,
		},
		{
			Name:      "Run with threshold",
			RunDir:    targetDir,
			Args:      []string{"--threshold", "10"},
			Validator: NoopValidator,
		},
		{
			Name:      "Run with verbose output",
			RunDir:    targetDir,
			Args:      []string{"--verbose"},
			Validator: NoopValidator,
		},
		{
			Name:      "Run with combination of options",
			RunDir:    targetDir,
			Args:      []string{"--format", "csv", "--output", "output.csv", "--exclude", "vendor/"},
			Validator: NoopValidator,
		},
	}

	RunGritTests(t, tests)
}

func TestGritErrorCases(t *testing.T) {
	// Target directory to analyze
	targetDir := `C:\Users\Victor\repos\fasthttp`

	// Test with non-existent directory
	nonExistentDir := `C:\Users\Victor\repos\non_existent_dir`

	tests := []GritTest{
		{
			Name:        "Run with invalid format option",
			RunDir:      targetDir,
			Args:        []string{"--format", "invalid_format"},
			ExpectError: true,
			Validator:   ContainsValidator,
		},
		{
			Name:        "Run with non-existent directory",
			RunDir:      nonExistentDir,
			Args:        []string{},
			ExpectError: true,
			Validator:   NoopValidator,
		},
		{
			Name:        "Run with invalid engine",
			RunDir:      targetDir,
			Args:        []string{"--engine", "invalid_engine"},
			ExpectError: true,
			Validator:   ContainsValidator,
		},
	}

	RunGritTests(t, tests)
}

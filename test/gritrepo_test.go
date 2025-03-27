package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// To run these integration tests, you need to have following:
// 1) Grit installed in your system.
// 		You can install it via go install.
// 		You can download it from https://github.com/vbvictor/grit/releases and place it in your PATH.

func createTempRepo(t *testing.T) (string, string) { //nolint
	t.Helper()

	tempDir := t.TempDir()
	gritDir := filepath.Join(tempDir, "grit")
	err := os.MkdirAll(gritDir, 0o755)
	require.NoError(t, err)

	Unbundle(t, filepath.Join("..", "testdata", "bundles", "grit-test.bundle"), gritDir)

	return tempDir, gritDir
}

func TestGritBasicFunctionality(t *testing.T) {
	_, gritDir := createTempRepo(t)

	tests := []GritTest{
		{
			Name:        "Run no params",
			RunDir:      gritDir,
			Args:        []string{},
			Validator:   NewContainsValidator("GRIT is an all-in-one cli tool"),
			ExpectError: false,
		},
		{
			Name:        "Run general help",
			RunDir:      gritDir,
			Args:        []string{"--help"},
			Validator:   NewContainsValidator(`Use "grit [command] --help" for more information about a command.`),
			ExpectError: false,
		},
		{
			Name:        "Run plot help",
			RunDir:      gritDir,
			Args:        []string{"plot", "--help"},
			Validator:   NewContainsValidator(`Creates visual graphs for code metrics.`),
			ExpectError: false,
		},
		{
			Name:        "Run report help",
			RunDir:      gritDir,
			Args:        []string{"report", "--help"},
			Validator:   NewContainsValidator(`Creates maintainability report based on churn, complexity and coverage`),
			ExpectError: false,
		},
		{
			Name:   "Run stat help",
			RunDir: gritDir,
			Args:   []string{"stat", "--help"},
			Validator: NewContainsValidator(`Calculate code metrics`,
				`Finds files with the most changes in git repository`,
				`Finds the most complex files`,
				`Finds files with the least unit-test coverage`),
			ExpectError: false,
		},
		{
			Name:        "Run stat churn help",
			RunDir:      gritDir,
			Args:        []string{"stat", "churn", "--help"},
			Validator:   NewContainsValidator(`Finds files with the most changes in git repository`),
			ExpectError: false,
		},
		{
			Name:        "Run stat complexity help",
			RunDir:      gritDir,
			Args:        []string{"stat", "complexity", "--help"},
			Validator:   NewContainsValidator(`Finds the most complex files`),
			ExpectError: false,
		},
		{
			Name:        "Run stat coverage help",
			RunDir:      gritDir,
			Args:        []string{"stat", "coverage", "--help"},
			Validator:   NewContainsValidator(`Finds files with the least unit-test coverage`),
			ExpectError: false,
		},
	}

	RunGritTests(t, tests)
}

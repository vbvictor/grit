package test

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// GritTest represents a test case for the grit command.
type GritTest struct {
	Name        string          // Name of the test
	RunDir      string          // Directory to run grit on
	Args        []string        // Arguments to pass to grit
	Validator   OutputValidator // Function to validate output
	ExpectError bool            // Whether we expect grit to exit with non-zero status
}

type OutputValidator func(t *testing.T, stdout, stderr string) bool

func NoopValidator(_ *testing.T, _, _ string) bool {
	return true
}

func NewContainsValidator(expectedOutputs ...string) OutputValidator {
	return func(t *testing.T, stdout, _ string) bool {
		t.Helper()

		for _, output := range expectedOutputs {
			if strings.Contains(stdout, output) {
				return true
			}
		}

		return false
	}
}

// RunGritTest runs a single grit test.
func RunGritTest(t *testing.T, test GritTest) {
	t.Helper()
	t.Run(test.Name, func(t *testing.T) {
		gritPath, err := findGritExecutable()
		if err != nil {
			t.Fatalf("Failed to find grit executable: %v", err)
		}

		cmd := exec.Command(gritPath, test.Args...)
		cmd.Dir = test.RunDir

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err = cmd.Run()

		if err != nil && !test.ExpectError {
			t.Errorf("grit command failed unexpectedly: %v", err)
			t.Logf("Stdout: %s", stdout.String())
			t.Logf("Stderr: %s", stderr.String())

			return
		} else if err == nil && test.ExpectError {
			t.Errorf("grit command succeeded but was expected to fail")

			return
		}

		if test.Validator != nil {
			if !test.Validator(t, stdout.String(), stderr.String()) {
				t.Fail()
			}
		}
	})
}

// findGritExecutable tries to find the grit executable in common locations.
func findGritExecutable() (string, error) {
	// Check if grit is in PATH
	if path, err := exec.LookPath("grit"); err == nil {
		return path, nil
	}

	return "", errors.New("could not find grit executable")
}

// RunGritTests runs multiple grit tests.
func RunGritTests(t *testing.T, tests []GritTest) {
	t.Helper()

	for _, test := range tests {
		RunGritTest(t, test)
	}
}

func Unbundle(t *testing.T, src, dst string) {
	t.Helper()

	cmd := exec.Command("git", "clone", src, dst)
	require.NoError(t, cmd.Run())
}

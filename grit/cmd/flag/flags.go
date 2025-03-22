package flag

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/vbvictor/grit/pkg/complexity"
	"github.com/vbvictor/grit/pkg/git"
)

// Global flags used across all commands

type OutputType = string

var (
	// Common flags.
	Verbose    bool
	Extensions []string

	// Plot command flags.
	Plot       string
	OutputFile string

	// Git/Churn related flags.
	Top         int
	ExcludePath string
	SortBy      string
	Since       string
	Until       string

	// Complexity flags.
	Threads int
	Engine  string

	// Report flags.
	ChurnFactor      float64
	ComplexityFactor float64
	CoverageFactor   float64

	// Output format flags.
	Tabular                OutputType = "tabular"
	CSV                    OutputType = "csv"
	AvailableOutputFormats            = []OutputType{Tabular, CSV}

	// Coverage Run formats.
	Always = "always"
	Never  = "never"
	Auto   = "auto"

	OutputFormat OutputType
)

// Default values.
const (
	DefaultOutputFormat = "tabular"
	DefaultEngine       = "gocyclo"
	DefaultThreads      = 1
	DefaultTop          = 10
)

const (
	LongSort         = "sort"
	LongTop          = "top"
	LongVerbose      = "verbose"
	LongExclude      = "exclude"
	LongExtensions   = "ext"
	LongSince        = "since"
	LongUntil        = "until"
	LongFormat       = "format"
	LongEngine       = "complexity-engine"
	LongRunCoverage  = "run-tests"
	LongFileCoverage = "coverage-file"

	// Flag shortcuts.
	ShortTop          = "t"
	ShortVerbose      = "v"
	ShortExt          = "e"
	ShortSince        = "s"
	ShortUntil        = "u"
	ShortFormat       = "f"
	ShortEngine       = "g"
	ShortRunCoverage  = "r"
	ShortFileCoverage = "c"

	DefaultUntil = "current date"
	DefaultSince = "one year ago"
)

func LogIfVerbose(format string, args ...any) {
	if Verbose {
		fmt.Printf(format, args...)
	}
}

type AbsRepoPathError struct {
	Path string
}

func (e *AbsRepoPathError) Error() string {
	return "failed to get absolute path from " + e.Path
}

var (
	ErrCoverageNotFound = errors.New("failed to find file with code coverage")
	ErrReadCoverage     = errors.New("failed to read coverage file")
	ErrRunCoverage      = errors.New("failed to run coverage")
)

func VerboseFlag(f *pflag.FlagSet, verbose *bool) {
	f.BoolVarP(verbose, LongVerbose, ShortVerbose, false, "Show detailed progress")
}

func OutputFlag(f *pflag.FlagSet, output *string, defaultValue string) {
	f.StringVarP(output, "output", "o", defaultValue, "Output graph file name")
}

func ExcludeRegexFlag(f *pflag.FlagSet, excludeRegex *string) {
	f.StringVar(excludeRegex, LongExclude, "", "Exclude files matching regex pattern")
}

func ChurnTypeFlag(f *pflag.FlagSet, churnType *string, defaultValue string) {
	f.StringVar(churnType, "churn-type", defaultValue,
		fmt.Sprintf("Specify churn type: [%s, %s]", git.Changes, git.Commits))
}

func SinceFlag(f *pflag.FlagSet, since *string) {
	f.StringVarP(since, LongSince, ShortSince, "", "Start date for analysis in format 'YYYY-MM-DD'")
}

func UntilFlag(f *pflag.FlagSet, until *string) {
	f.StringVarP(until, LongUntil, ShortUntil, "", "End date for analysis in format 'YYYY-MM-DD'")
}

func EngineFlag(f *pflag.FlagSet, engine *string, defaultValue string) {
	f.StringVarP(engine, LongEngine, ShortEngine, defaultValue,
		fmt.Sprintf("Specify complexity calculation engine: [%s, %s]", complexity.Gocyclo, complexity.Gocognit))
}

func SortFlag(f *pflag.FlagSet, sortBy *string, defaultValue string, description string) { //nolint: gocritic // unified
	f.StringVar(sortBy, LongSort, defaultValue, description)
}

func RunCoverageFlag(f *pflag.FlagSet, runCoverage *string) {
	f.StringVarP(runCoverage, LongRunCoverage, ShortRunCoverage, Auto,
		`Specify tests run format:
  'Auto' will run unit tests if coverage file is not found
  'Always' will run unit tests on every invoke 'coverage' command
  'Never' will never run unit tests and always look for present test-coverage file
`)
}

func CoverageFilenameFlag(f *pflag.FlagSet, filename *string) {
	f.StringVarP(filename, LongFileCoverage, ShortFileCoverage, "coverage.out",
		"Name of code coverage file to read or create")
}

func TopFlag(f *pflag.FlagSet, top *int) {
	f.IntVarP(top, LongTop, ShortTop, DefaultTop, "Number of top files to display")
}

func OutputFormatFlag(f *pflag.FlagSet, format *string) {
	f.StringVarP(format, LongFormat, ShortFormat, Tabular,
		fmt.Sprintf("Specify output format: [%s, %s]", Tabular, CSV))
}

func ExtensionsFlag(f *pflag.FlagSet, extensions *[]string) {
	f.StringSliceVarP(extensions, LongExtensions, ShortExt, nil,
		"Only include files with given extensions in comma-separated list, e.g. 'go,h,c'")
}

func ComplexityEngineFlag(f *pflag.FlagSet, engine *string) {
	f.StringVarP(engine, LongEngine, ShortEngine, complexity.Gocyclo,
		fmt.Sprintf(`Specify complexity calculation engine: [%s, %s, %s].
When CSV engine is specified, GRIT will try to read function complexity data from CSV file
'complexity.csv' located in <path>. The file should have following fields:
"filename,function,complexity,line-count (optional),packages (optional)"
`, complexity.Gocyclo, complexity.Gocognit, complexity.CSV))
}

func PerfectCoverageFlag(f *pflag.FlagSet, perfectCoverage *float64) {
	f.Float64Var(perfectCoverage, "perfect-coverage", 100.0, //nolint:mnd // default value
		"Specify code coverage penalty threshold")
}

package flag

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
	RunCoverage      bool
	CoverageFile     string

	// Output format flags.
	JSON                   OutputType = "json"
	Tabular                OutputType = "tabular"
	AvailableOutputFormats            = []OutputType{JSON, Tabular}

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
	LongEngine       = "engine"
	LongRunCoverage  = "run"
	LongFileCoverage = "file"

	// Flag shortcuts.
	ShortTop          = "t"
	ShortVerbose      = "v"
	ShortExt          = "e"
	ShortSince        = "s"
	ShortUntil        = "u"
	ShortFormat       = "f"
	ShortEngine       = "g"
	ShortRunCoverage  = "r"
	ShortFileCoverage = "f"

	DefaultUntil = "current date"
	DefaultSince = "one year ago"
)

func GetExtMap(extensions []string) map[string]struct{} {
	extMap := make(map[string]struct{})

	for _, ext := range extensions {
		extMap[ext] = struct{}{}
	}

	return extMap
}

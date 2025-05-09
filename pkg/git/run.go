package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"golang.org/x/exp/maps"
)

// ChurnType represents the type of sorting to be performed on the results of git log.
type ChurnType = string

var (
	Changes   ChurnType = "changes"
	Additions ChurnType = "additions"
	Deletions ChurnType = "deletions"
	Commits   ChurnType = "commits"
)

const (
	DefaultTop = 10
	HashLength = 40
)

var _ pflag.Value = (*Date)(nil)

type Date struct {
	time.Time
}

func (d *Date) Type() string {
	return "Date"
}

func (d *Date) String() string {
	return d.Format(time.DateOnly)
}

func (d *Date) Set(value string) error {
	parsedTime, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	*d = Date{parsedTime}

	return nil
}

type ChurnOptions struct {
	SortBy       ChurnType
	Top          int
	Path         string
	ExcludeRegex *regexp.Regexp
	Extensions   map[string]struct{}
	Since        time.Time
	Until        time.Time
	OutputFormat string
}

type ChurnChunk struct {
	File    string `json:"path"`
	Churn   int    `json:"changes"`
	Added   int    `json:"additions"`
	Removed int    `json:"deletions"`
	Commits int    `json:"commits"`
}

func getExtMap(extensions []string) map[string]struct{} {
	extMap := make(map[string]struct{})

	for _, ext := range extensions {
		extMap[ext] = struct{}{}
	}

	return extMap
}

func PopulateOpts(opts *ChurnOptions, extensionList []string,
	since string,
	until string,
	repoPath string,
	excludeRegex string,
) error {
	opts.Path = repoPath

	if err := setSinceOpt(opts, since); err != nil {
		return fmt.Errorf("error setting since option: %w", err)
	}

	if err := setUntilOpt(opts, until); err != nil {
		return fmt.Errorf("error setting until option: %w", err)
	}

	if extensionList != nil {
		opts.Extensions = getExtMap(extensionList)
	}

	if excludeRegex != "" {
		var err error

		opts.ExcludeRegex, err = regexp.Compile(excludeRegex)
		if err != nil {
			return fmt.Errorf("invalid exclude pattern: %w", err)
		}
	}

	return nil
}

func setSinceOpt(opts *ChurnOptions, since string) error {
	if since != "" {
		var err error

		opts.Since, err = time.Parse(time.DateOnly, since)
		if err != nil {
			return fmt.Errorf("error parsing since date: %w", err)
		}
	} else {
		opts.Since = time.Now().AddDate(-1, 0, 0)
	}

	return nil
}

func setUntilOpt(opts *ChurnOptions, until string) error {
	if until != "" {
		var err error

		opts.Until, err = time.Parse(time.DateOnly, until)
		if err != nil {
			return fmt.Errorf("error parsing until date: %w", err)
		}
	} else {
		opts.Until = time.Now()
	}

	return nil
}

func ReadGitChurn(repoPath string, opts *ChurnOptions) ([]*ChurnChunk, error) {
	cmd := buildGitCommand(opts)

	output, err := executeGitCommand(repoPath, cmd)
	if err != nil {
		return nil, err
	}

	fileStats := make(map[string]*ChurnChunk)
	lines := strings.Split(string(output), "\n")

	processLines(lines, fileStats, opts)

	return maps.Values(fileStats), nil
}

func buildGitCommand(opts *ChurnOptions) []string {
	cmd := []string{"git", "log", "--pretty=format:%H", "--numstat"}

	if !opts.Since.IsZero() {
		cmd = append(cmd, "--since="+opts.Since.Format(time.DateOnly))
	}

	if !opts.Until.IsZero() {
		cmd = append(cmd, "--until="+opts.Until.Format(time.DateOnly))
	}

	cmd = append(cmd, "--", ".")

	return cmd
}

func executeGitCommand(path string, cmd []string) ([]byte, error) {
	gitCmd := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec // This command is built via buildGitCommand
	gitCmd.Dir = path

	output, err := gitCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute git command: %w", err)
	}

	return output, nil
}

func processLines(lines []string, fileStats map[string]*ChurnChunk, opts *ChurnOptions) {
	currentCommit := ""
	modifiedInCommit := make(map[string]bool)

	for _, line := range lines {
		if line == "" {
			continue
		}

		if len(line) == HashLength {
			processCommit(currentCommit, modifiedInCommit, fileStats, opts)
			currentCommit = line
			modifiedInCommit = make(map[string]bool)
		} else {
			processFileLine(line, fileStats, modifiedInCommit, opts)
		}
	}
}

func processCommit(currentCommit string, modifiedInCommit map[string]bool, fileStats map[string]*ChurnChunk,
	opts *ChurnOptions,
) {
	if currentCommit != "" && len(modifiedInCommit) > 0 {
		for file := range modifiedInCommit {
			path := file

			if shouldSkipFile(path, opts) {
				continue
			}

			fileStats[path].Commits++
		}
	}
}

func processFileLine(line string, fileStats map[string]*ChurnChunk, modifiedInCommit map[string]bool,
	opts *ChurnOptions,
) {
	parts := strings.Fields(line)
	if len(parts) == 3 && isNumeric(parts[0]) && isNumeric(parts[1]) {
		additions, _ := strconv.Atoi(parts[0])
		deletions, _ := strconv.Atoi(parts[1])

		path := localizeClean(parts[2])

		if shouldSkipFile(path, opts) {
			return
		}

		updateFileStats(fileStats, path, additions, deletions)

		modifiedInCommit[path] = true
	}
}

func updateFileStats(fileStats map[string]*ChurnChunk, path string, additions, deletions int) {
	if _, exists := fileStats[path]; !exists {
		fileStats[path] = &ChurnChunk{File: path}
	}

	fileStats[path].Added += additions
	fileStats[path].Removed += deletions
	fileStats[path].Churn += additions + deletions
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)

	return err == nil
}

func shouldSkipFile(file string, opts *ChurnOptions) bool {
	if opts.ExcludeRegex != nil && opts.ExcludeRegex.MatchString(file) {
		return true
	}

	if opts.Extensions != nil {
		fileExt := filepath.Ext(file)

		if fileExt == "" {
			return false
		}

		if _, exists := opts.Extensions[fileExt[1:]]; !exists {
			return true
		}
	}

	return false
}

func SortAndLimit(result []*ChurnChunk, sortBy ChurnType, limit int) []*ChurnChunk {
	less := func() func(i, j int) bool {
		switch sortBy {
		case Changes:
			return func(i, j int) bool { return result[i].Churn > result[j].Churn }
		case Additions:
			return func(i, j int) bool { return result[i].Added > result[j].Added }
		case Deletions:
			return func(i, j int) bool { return result[i].Removed > result[j].Removed }
		case Commits:
			return func(i, j int) bool { return result[i].Commits > result[j].Commits }
		default:
			return nil
		}
	}()

	sort.Slice(result, less)

	// Limit the number of results
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result
}

func localizeClean(path string) string {
	localized, err := filepath.Localize(path)
	if err != nil {
		panic("failed top localize path " + path)
	}

	return filepath.Clean(localized)
}

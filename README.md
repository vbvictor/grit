# Go Refactoring Insight Tool

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

GRIT is a cli tool that helps developers understand their codebase
maintainability index through calculated key metrics: code churn, code complexity and test
coverage. Use calculated maintainability index to make decisions about
refactoring and testing priorities.

## Table of contents

- [Go Refactoring Insight Tool](#go-refactoring-insight-tool)
  - [Table of contents](#table-of-contents)
  - [What GRIT Measures](#what-grit-measures)
  - [Getting Started](#getting-started)
  - [Usage](#usage)
    - [Help command](#help-command)
    - [Churn command](#churn-command)
    - [Complexity command](#complexity-command)
    - [Coverage command](#coverage-command)
  - [Roadmap](#roadmap)
  - [Contributing](#contributing)
  - [License](#license)

## What GRIT Measures

<!-- - **Maintainability Score**: Combines metrics to rate maintainability index-->
<!-- - **Visual Analytics**: Generates churn vs complexity graphs -->
- **Code Churn**: Tracks how frequently files change over time.
- **Code Complexity**: Calculates cyclomatic complexity metric per file.
- **Test Coverage**: Analyzes test coverage percentage per file.

All of these metrics are useful when making decisions about:

- Best candidates for refactoring efforts
- What most complex files to cover with unit tests first

<!-- These metrics if measured regularly can address appearing maintainability issues in a large codebase. -->
  
## Getting Started

Install via tool via `go install` command:

```bash
go install github.com/vbvictor/grit@latest
```

Or Download the latest binary release from [Github Releases](https://github.com/vbvictor/grit/releases) page.

## Usage

### Help command

Run `grit -h` to check out all available commands and general help for `grit`:

```sh
All-in-one tool for getting refactoring statistics.

Usage:
  grit [command]

Available Commands:
  help        Help about any command
  stat        Get various statistics about code in the repository

Flags:
  -h, --help   help for grit
```

### Churn command

### Complexity command

### Coverage command

## Roadmap

- Features
  - [x] collect churn, complexity, coverage metrics
  - [ ] create maintainability report
  - [ ] render churn vs complexity graph
  - [ ] render churn vs complexity vs coverage graph
  - [ ] add more output formats
  - [ ] support custom files with metrics to supporting other languages
- Improvements
  - [ ] enhance readme with more examples and metric descriptions
  - [ ] use different library for formatting tabular output

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.

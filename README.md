# Go Refactoring Insight Tool

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)]()

GRIT helps developers understand their codebase maintainability through key metrics including code churn, complexity, and test coverage.
WUse these insights to make data-driven decisions about refactoring and testing priorities.

## What It Measures

- **Code Churn**: Tracks how frequently files change over time
- **Code Complexity**: Calculates cyclomatic complexity metrics
- **Test Coverage**: Analyzes test coverage percentage per file
- **Maintainability Score**: Combines metrics to rate maintainability
- **Visual Analytics**: Generates churn vs complexity graphs

## Key Features

- **Metrics**
  - Churn by file
  - Complexity by file
  - Code coverage by file
- Reports
  - Generate maintainability report
  - Generate churn vs complexity graph

## Getting Started

Install via `go`:

```bash
go install  github.com/vbvictor/grit@latest
```

## Usage

### Help commands

Run `grit -h` to check out all available commands:

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

## Examples

```bash
grit stat complexity
```

## Roadmap

- Features
  - [x] collect churn, complexity, coverage metrics
  - [ ] create maintainability report
  - [ ] render churn vs complexity graph
  - [ ] render churn vs complexity vs coverage graph
  - [ ] add more output formats
  - [ ] support custom user-files with metrics for supporting languages
- Improvements
  - [ ] use different library for formatting tabular output

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.

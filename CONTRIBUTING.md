# Contributing to goscaf

First off, thank you for considering contributing to `goscaf`! It's people like you who make it such a great tool.

## How Can I Contribute?

### Reporting Bugs
- Use the [Bug Report template](.github/ISSUE_TEMPLATE/bug_report.md).
- Describe the steps to reproduce the issue.
- Include your Go version and OS.

### Suggesting Enhancements
- Use the [Feature Request template](.github/ISSUE_TEMPLATE/feature_request.md).
- Explain why this enhancement would be useful to most users.

### Pull Requests
1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests!
3. Ensure the test suite passes (`make test`).
4. Run the linter (`make lint`).
5. Run the smoke test (`make smoke-test`).
6. Update documentation if necessary.

## Development Setup

1. Clone your fork.
2. Install dependencies: `go mod download`.
3. Install tools: `make install-tools` (if available) or manual install of `golangci-lint`.
4. Build the project: `make build`.

## Project Structure
- `main.go`: Entry point — delegates immediately to `cmd.Execute()`.
- `cmd/`: CLI commands (Cobra-based).
    - `cmd/root.go`: Root command and banner.
    - `cmd/init.go`: `goscaf init` — collects config and drives the generator.
    - `cmd/add.go`: `goscaf add` — adds a new service scaffold to an existing project.
- `internal/`: Core logic — not importable by external packages.
    - `internal/config/`: `ProjectConfig` struct and typed constants for framework, logger, and database choices.
    - `internal/userconfig/`: Loads and merges `~/.goscaf.yaml` (global) and `./.goscaf.yaml` (local) defaults.
    - `internal/prompt/`: Interactive survey prompts; accepts a `UserConfig` to pre-fill defaults.
    - `internal/generator/`: Orchestrates file writes for `goscaf init` and `goscaf add`.
    - `internal/templates/`: Go `text/template` strings for every generated file.

## Style Guide
- We follow standard Go idioms.
- Run `go fmt ./...` before committing.
- Clear, descriptive commit messages are appreciated.

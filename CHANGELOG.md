# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added

- **Output Format Module** (`pkg/output/`)
  - Printer interface for unified output formatting
  - TablePrinter for table format output
  - JSONPrinter for JSON format output
  - SimplePrinter for simple text output
  - TemplatePrinter for Go template support

- **Format Options**
  - `--format` option for `issue list` command (json/simple/table)
  - `--time-format` option for relative/absolute time display

- **Time Formatting**
  - Relative time display (e.g., "2 hours ago")
  - Support for both relative and absolute time formats

- **Template System**
  - Go template support for custom output formatting
  - Template functions: upper, lower, trunc, json

### Fixed

- Dynamic width calculation for issue/PR number alignment
  - Fixed alignment issue when issue numbers exceed 6 digits
  - Affects `gc issue list`, `gc pr list`, `gc release list`

### Changed

- Refactored issue list output to use Printer interface
- Improved output consistency across commands

## [0.3.9] - 2026-03-XX

### Added
- Initial release features
- Basic issue, PR, repo, release commands
- Authentication via token
- JSON output support

[Unreleased]: https://gitcode.com/gitcode-cli/cli/compare/v0.3.9...HEAD
[0.3.9]: https://gitcode.com/gitcode-cli/cli/releases/tag/v0.3.9
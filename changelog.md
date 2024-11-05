# Changelog

All notable changes to the `delock` project will be documented in this file.

## [0.3.0] - 2024-11-05

### Fixed

- Resolved goroutine accumulation issue in `Mutex` and `RWMutex`.
- Updated `Lock` and `RLock` methods to use `context.WithTimeout` for efficient timeout management, ensuring that goroutines are terminated properly when the timeout expires.
- Simplified locking mechanism to avoid redundant channels and intermediate goroutines, significantly improving resource management and stability.

## [0.2.0] - 2023-12-06

### Improved

- Deadlock report format: Reports are now more concise, grouping lock occurrences by type and simplifying the output to focus on key code lines for quicker analysis.

## [0.1.0] - 2023-12-05

### Added

- Initial release of `delock`.
- Features include deadlock detection for `Mutex` and `RWMutex` with stack trace logging and customizable timeout.
- `README.md` with detailed usage instructions and examples.
- `CHANGELOG.md` for tracking changes and version history.
- `LICENSE` file with MIT license information.
- `example-go16.sh` script that tests library works on go version 1.16.
- Additional examples to demonstrate various deadlock scenarios and the usage of `delock` in different contexts.

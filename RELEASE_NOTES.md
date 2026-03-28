# gscex v1.1.0 Release Notes

## New Features

- **File Filter**: Press Ctrl+F to filter results by filename substring (e.g., "zm_tomb")
- **Loading Indicator**: Spinner animation shows during search with query details
- **Input Protection**: 'q' and '?' keys blocked while typing in input fields

## Bug Fixes

- **Search Completeness**: Files sorted alphabetically for consistent results, limit increased to 10,000
- **Escape Key**: Now unfocuses input fields in search mode
- **Results Navigation**: 'q' key properly works in results/preview modes

## Technical Changes

- MaxResults increased from 100 to 10,000 for both TUI and CLI
- File search now uses sorted iteration for predictable results
- Added spinner animation with 100ms tick rate
- Enhanced keyboard shortcut handling

## Binaries

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

SHA256 checksums included in `checksums.txt`
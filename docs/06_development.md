# Development Guide

**Description:** Guide for building, modifying, and contributing to the gscex-cli tool.

**Goal:** Enable developers to understand the codebase, build from source, and extend functionality.

**Applicable for:** Developers who want to contribute or customize the tool.

---

## Index

| File | Topic | Purpose |
|------|-------|---------|
| [INDEX.md](INDEX.md) | Main Index | Project overview and navigation |
| [01_architecture.md](01_architecture.md) | Architecture | How the tool works |
| [02_installation.md](02_installation.md) | Installation | Setup instructions |
| [03_usage.md](03_usage.md) | Usage Guide | Command reference |
| [04_stock_scripts.md](04_stock_scripts.md) | Stock Scripts | Stock script structure |
| [05_search_features.md](05_search_features.md) | Search Features | Search capabilities |

---

## Project Structure

```
gscex-cli/
├── cmd/
│   └── gscex/
│       └── main.go           # Entry point
├── pkg/
│   ├── index/
│   │   ├── indexer.go        # GSC parsing and indexing
│   │   └── storage.go        # Index storage/retrieval
│   ├── search/
│   │   ├── engine.go         # Search implementation
│   │   └── results.go        # Result formatting
│   ├── config/
│   │   └── config.go         # Configuration management
│   └── git/
│       └── client.go         # Git operations
├── go.mod
├── go.sum
└── README.md
```

---

## Building from Source

### Prerequisites

- Go 1.21+
- Git

### Development Setup

```bash
# Clone repository
git clone https://github.com/yourusername/gscex-cli.git
cd gscex-cli

# Download dependencies
go mod download

# Build for current platform
go build -o gscex cmd/gscex/main.go

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o gscex-linux cmd/gscex/main.go
GOOS=windows GOARCH=amd64 go build -o gscex.exe cmd/gscex/main.go
GOOS=darwin GOARCH=amd64 go build -o gscex-macos cmd/gscex/main.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/index/
go test ./pkg/search/
```

---

## Architecture Details

### Indexer Package

Responsible for parsing GSC files and building the searchable index.

**Key types:**
```go
type Function struct {
    Name      string
    File      string
    Line      int
    Signature string
    Context   []string
}

type Index struct {
    Functions map[string]Function
    Methods   map[string][]Function
    Dvars     map[string][]Location
    Includes  map[string][]string
}
```

**Parsing regex patterns:**
```go
// Function definition: function_name(params) {
var funcDefRegex = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*\([^)]*\)\s*\{`)

// Method call: entity.method(params)
var methodCallRegex = regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)\.([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`)

// Include: #include path;
var includeRegex = regexp.MustCompile(`^#include\s+([^;]+);`)
```

### Search Package

Implements search algorithms over the index.

```go
type Engine struct {
    index *Index
}

func (e *Engine) SearchFunction(name string) (*Function, error)
func (e *Engine) SearchMethod(entity, name string) ([]Function, error)
func (e *Engine) SearchText(pattern string, opts SearchOptions) ([]Result, error)
```

### CLI Package

Command-line interface using cobra or urfave/cli.

```go
var rootCmd = &cobra.Command{
    Use:   "gscex",
    Short: "GSC Reference CLI for T6",
}

func init() {
    rootCmd.AddCommand(initCmd)
    rootCmd.AddCommand(searchCmd)
    rootCmd.AddCommand(funcCmd)
    rootCmd.AddCommand(methodCmd)
    rootCmd.AddCommand(filesCmd)
    rootCmd.AddCommand(updateCmd)
}
```

---

## Adding New Features

### Example: Adding a New Command

1. Create command file in `cmd/gscex/commands/`:

```go
// cmd/gscex/commands/stats.go
package commands

import (
    "fmt"
    "github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
    Use:   "stats",
    Short: "Show index statistics",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}
```

2. Register in `main.go`:

```go
rootCmd.AddCommand(commands.StatsCmd)
```

3. Write tests:

```go
// pkg/commands/stats_test.go
func TestStatsCommand(t *testing.T) {
    // Test implementation
}
```

---

## Contributing Guidelines

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature-name`
3. **Write** tests for new functionality
4. **Ensure** all tests pass: `go test ./...`
5. **Commit** with clear messages
6. **Push** to your fork
7. **Submit** a pull request

### Code Style

- Follow Go conventions
- Use `gofmt` for formatting
- Add comments for exported functions
- Write tests for all new code
- Minimize lines of code (per aikb principles)

---

## Release Process

```bash
# Tag version
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Build binaries
./scripts/build.sh

# Upload to GitHub releases
```

---

## Debugging

```bash
# Verbose mode
gscex -v search "pattern"

# Debug build
go build -gcflags="all=-N -l" -o gscex-debug cmd/gscex/main.go

# Profile performance
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

---

## Configuration

Config file location: `~/.gscex/config.json`

```json
{
  "scripts_repo": "https://github.com/plutoniummod/t6-scripts",
  "scripts_branch": "main",
  "cache_dir": "~/.gscex",
  "auto_update": false,
  "max_results": 20,
  "context_lines": 3,
  "fuzzy_threshold": 0.8
}
```

---

*Back to: [INDEX.md](INDEX.md)*

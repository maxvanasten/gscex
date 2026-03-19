# GSC Explorer (gscex)

A powerful CLI and TUI tool for exploring Call of Duty: Black Ops 2 (Plutonium T6) GSC scripts. Search through 1,469 stock script files containing 10,950+ functions instantly.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)
![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-lightgrey.svg)

## Features

- 🔍 **Instant Search** - Find functions, methods, and text across all stock scripts
- 🖥️ **Interactive TUI** - Browse with an intuitive terminal interface (vim-style keys)
- 🎯 **Fuzzy Matching** - Search "giveweapon" finds "give_weapon" automatically
- 📁 **File Browser** - Navigate the complete script hierarchy
- 📄 **Code Preview** - View files with line highlighting and context
- ⚡ **Async Loading** - Results stream in real-time
- 🔧 **Multiple Modes** - Text search, function search, method search, file search
- 📦 **Cross Platform** - Works on Linux, macOS, and Windows

## Quick Start

### Installation

#### Download Pre-built Binary

Download the latest release for your platform from the [Releases](https://github.com/yourusername/gscex/releases) page.

#### Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/gscex.git
cd gscex

# Build
go build -o gscex ./cmd/gscex/

# Or install to $GOPATH/bin
go install ./cmd/gscex/
```

### Usage

```bash
# Initialize (download and index stock scripts)
gscex init

# Launch interactive TUI browser (recommended)
gscex tui

# Or use CLI commands:
# Search for text
gscex search "magic bullet"

# Find function definition
gscex func giveweapon

# Search method calls
gscex method player give_weapon

# List files
gscex files "_zm"

# Show tool info
gscex info
```

## Interactive TUI Mode

The TUI provides the best experience for exploring GSC scripts:

```bash
gscex tui
```

**Keyboard Shortcuts:**

| Key | Action |
|-----|--------|
| `?` | Toggle help |
| `Tab` | Switch search mode (text/func/method/files) |
| `Enter` | Execute search / Open result |
| `Ctrl+N` | Next result / Scroll down |
| `Ctrl+P` | Previous result / Scroll up |
| `Ctrl+F` | Focus search box |
| `Ctrl+B` | Back to results |
| `↑/↓` or `j/k` | Scroll in preview |
| `PgUp/PgDn` | Page scroll |
| `q` or `Esc` | Go back / Quit |

**Features:**
- Live fuzzy search through all functions
- Results sorted by file and line number
- Current line highlighting in preview
- Async loading (up to 100 results)
- Smooth scrolling

## Search Modes

### Text Search
Find any text across all 1,469 script files.

```bash
gscex search "player_damage"
gscex search "callback" --json
```

### Function Search
Find function definitions with fuzzy matching.

```bash
# Exact match
gscex func give_weapon

# Fuzzy match (underscore-insensitive)
gscex func giveweapon  # Finds give_weapon
```

### Method Search
Find method calls on specific entities.

```bash
# Format: "entity method"
gscex method player give_weapon
gscex method zombie spawn
```

### File Search
Browse and search filenames.

```bash
gscex files "_zm"        # All zombie-related files
gscex files "gametype"   # Gametype scripts
```

## CLI Options

All commands support these flags:

```bash
# JSON output for scripting
--json, -j

# Limit results
--max, -n 50

# Show context lines
--context, -c 3

# Files only (no content)
--files-only, -f
```

## Project Structure

```
gscex/
├── cmd/gscex/          # Main application
│   ├── main.go         # CLI commands
│   └── tui.go          # Interactive TUI
├── pkg/
│   ├── config/         # Configuration management
│   ├── git/            # Git operations
│   ├── index/          # GSC parsing and indexing
│   └── search/         # Search engine
├── docs/               # Documentation
└── scripts/            # Build scripts
```

## Configuration

Config stored in `~/.gscex/config.json`:

```json
{
  "scripts_repo": "https://github.com/plutoniummod/t6-scripts",
  "scripts_branch": "main",
  "cache_dir": "~/.gscex",
  "max_results": 20,
  "context_lines": 3
}
```

## Development

```bash
# Run tests
go test ./...

# Build for current platform
go build -o gscex ./cmd/gscex/

# Build all platforms
./scripts/build-release.sh v1.0.0
```

## Documentation

- [Architecture](docs/01_architecture.md) - How gscex works
- [Installation](docs/02_installation.md) - Detailed setup
- [Usage Guide](docs/03_usage.md) - Complete command reference
- [Stock Scripts](docs/04_stock_scripts.md) - Understanding T6 scripts
- [Search Features](docs/05_search_features.md) - Advanced search
- [Development](docs/06_development.md) - Contributing guide

## Stock Scripts

This tool indexes the official Plutonium T6 stock scripts:
- **1,469 files** across MP and ZM modes
- **10,950+ functions** indexed
- **10,246+ method calls** tracked
- All core utilities, gametypes, and zombie scripts included

## License

MIT License - See [LICENSE](LICENSE) for details.

## Acknowledgments

- [Plutonium](https://plutonium.pw/) - For the T6 modding platform
- [t6-scripts](https://github.com/plutoniummod/t6-scripts) - Stock script repository
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework

---

**Note:** This tool is for educational and modding purposes. Not affiliated with Activision or Treyarch.

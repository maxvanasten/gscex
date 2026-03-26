# GSC Explorer (gscex)

A powerful CLI and TUI tool for exploring Call of Duty: Black Ops 1-2 (Plutonium T5/T6) GSC scripts. Search through thousands of stock script files containing 10,000+ functions instantly.

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

Download the raw binary for your platform from the [Releases](https://github.com/maxvanasten/gscex/releases) page. No installation needed - just download and run!

**Available binaries:**
- `gscex-v1.0.0-linux-amd64` - Linux x86_64
- `gscex-v1.0.0-linux-arm64` - Linux ARM64
- `gscex-v1.0.0-darwin-amd64` - macOS Intel
- `gscex-v1.0.0-darwin-arm64` - macOS Apple Silicon
- `gscex-v1.0.0-windows-amd64.exe` - Windows x64

**Linux/macOS:**
```bash
chmod +x gscex-v1.0.0-linux-amd64
./gscex-v1.0.0-linux-amd64 init
```

**Windows:**
Just double-click or run from PowerShell:
```powershell
gscex-v1.0.0-windows-amd64.exe init
```

#### Build from Source

```bash
# Clone the repository
git clone https://github.com/maxvanasten/gscex.git
cd gscex

# Build
go build -o gscex ./cmd/gscex/

# Or install to $GOPATH/bin
go install ./cmd/gscex/
```

### Usage

```bash
# Initialize (download and index stock scripts for both T5 and T6)
gscex init

# Or initialize only specific game
gscex init t5    # Only Black Ops 1
gscex init t6    # Only Black Ops 2

# Launch interactive TUI browser (recommended)
gscex tui

# Search across all games (default)
gscex search "magic bullet"

# Search only specific game
gscex --game t6 search "magic bullet"
gscex --game t5 search "give_weapon"

# Find function definition
gscex func giveweapon

# Search method calls
gscex method player give_weapon

# List files
gscex files "_zm"

# Show tool info
gscex info

# Update scripts (both games by default)
gscex update
gscex update t6  # Update only T6
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
- Live fuzzy search through all functions (T5 and T6)
- Results sorted by file and line number
- Current line highlighting in preview
- Async loading (up to 100 results)
- Smooth scrolling
- Multi-game support (automatically shows stats for all indexed games)

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

# Select specific game (t5 or t6)
--game t5
--game t6
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
  "games": {
    "t5": {
      "scripts_repo": "https://github.com/plutoniummod/t5-scripts",
      "scripts_branch": "main"
    },
    "t6": {
      "scripts_repo": "https://github.com/plutoniummod/t6-scripts",
      "scripts_branch": "main"
    }
  },
  "cache_dir": "~/.gscex",
  "max_results": 20,
  "context_lines": 3,
  "default_game": "t6"
}
```

**Game Selection:**
- Use `--game t5` or `--game t6` flag to search specific game
- Default searches all indexed games
- T5 = Call of Duty: Black Ops (2010)
- T6 = Call of Duty: Black Ops 2 (2012)

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

This tool indexes the official Plutonium stock scripts:

**T6 (Black Ops 2):**
- **1,469 files** across MP and ZM modes
- **10,950+ functions** indexed
- **10,246+ method calls** tracked
- All core utilities, gametypes, and zombie scripts included

**T5 (Black Ops 1):**
- Stock scripts for multiplayer and zombies
- Complete function and method indexing

Search works across all indexed games by default, or use `--game t5` or `--game t6` to filter.

## License

MIT License - See [LICENSE](LICENSE) for details.

---

**Note:** This tool is for educational and modding purposes. Not affiliated with Activision or Treyarch.

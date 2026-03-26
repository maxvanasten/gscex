# GSC Explorer (gscex)

A CLI and TUI tool for exploring Call of Duty: Black Ops 1-2 (Plutonium T5/T6) GSC scripts. Search through thousands of stock script files containing 10,000+ functions instantly.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.23+-00ADD8.svg)
![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-lightgrey.svg)

## Features

- **Interactive TUI** - Browse with a keyboard-driven terminal interface
- **Instant Search** - Find functions, methods, and text across all stock scripts
- **Context Display** - See 2 lines before/after matches with the matched line highlighted
- **Fuzzy Matching** - Search "giveweapon" finds "give_weapon" automatically
- **Multi-Game Support** - Search T5 (Black Ops 1) and T6 (Black Ops 2) simultaneously
- **Code Preview** - View files with line highlighting
- **JSON Output** - Export results for scripting with `--json`
- **Cross Platform** - Works on Linux, macOS, and Windows

## Quick Start

### Installation

#### Download Pre-built Binary

Download the latest binary from the [Releases](https://github.com/maxvanasten/gscex/releases) page. No installation required.

**Linux/macOS:**
```bash
chmod +x gscex-v1.0.0-linux-amd64
./gscex-v1.0.0-linux-amd64 init
```

**Windows:**
```powershell
# Run from PowerShell or Command Prompt
gscex-v1.0.0-windows-amd64.exe init
```

#### Build from Source

```bash
git clone https://github.com/maxvanasten/gscex.git
cd gscex
go build -o gscex ./cmd/gscex/
```

### Usage

```bash
# Initialize (download and index stock scripts)
gscex init              # Both T5 and T6
gscex init t5           # Only Black Ops 1
gscex init t6           # Only Black Ops 2

# Launch interactive TUI (recommended)
gscex tui

# Search commands
gscex search "magic bullet"
gscex func giveweapon
gscex method player give_weapon
gscex files "_zm"

# Search specific game
gscex --game t6 search "player_damage"

# Update scripts
gscex update
```

## Interactive TUI Mode

The TUI provides the most efficient way to explore GSC scripts:

```bash
gscex tui
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` | Switch search mode (text/func/method/files) |
| `Enter` | Execute search / Open file at line |
| `Ctrl+F` | Focus search input |
| `Ctrl+N` | Next result / Scroll down |
| `Ctrl+P` | Previous result / Scroll up |
| `Ctrl+B` | Back to results |
| `?` | Toggle help overlay |
| `Esc` or `q` | Go back / Quit |

### Search Modes

The TUI supports four search modes, cycled with `Tab`:

1. **Text Search** - Find any text across all script files
2. **Function Search** - Find function definitions and their usages
3. **Method Search** - Find method calls on specific entities (e.g., `player give_weapon`)
4. **File Search** - Browse and filter the complete file hierarchy

### TUI Features

- **Live search results** with context (2 lines before/after each match)
- **Highlighted matches** displayed with `>>` prefix
- **File preview** with line highlighting and keyboard navigation
- **Multi-game support** - automatically loads all indexed games
- **Stats display** showing indexed files, functions, and methods per game

## CLI Commands

### Text Search

Search for any text across all script files:

```bash
# Basic search
gscex search "player_damage"

# JSON output for scripting
gscex search "callback" --json

# Limit results
gscex search "magic" --max 50

# Adjust context lines
gscex search "weapon" --context 5

# Files only (no content)
gscex search "zombie" --files-only
```

### Function Search

Find function definitions with fuzzy matching:

```bash
# Exact match
gscex func give_weapon

# Fuzzy match (underscore-insensitive)
gscex func giveweapon  # Finds give_weapon

# JSON output with definition and all usages
gscex func spawn --json
```

### Method Search

Find method calls on specific entities:

```bash
# Format: gscex method [entity] [method]
gscex method player give_weapon
gscex method zombie spawn
gscex method self notify

# JSON output
gscex method player damage --json
```

### File Search

Browse and search filenames:

```bash
# List all files matching pattern
gscex files "_zm"        # Zombie-related files
gscex files "gametype"   # Gametype scripts
gscex files "main"       # Main scripts

# JSON output
gscex files "utility" --json
```

### Other Commands

```bash
# Show tool and index information
gscex info

# Update scripts (pull latest and rebuild index)
gscex update              # All games
gscex update t6           # Specific game

# Shell completion
gscex completion bash     # Bash
gscex completion zsh      # Zsh
gscex completion fish     # Fish
```

## Global Flags

All commands support these flags:

```bash
# Select specific game
--game t5          # Black Ops 1 only
--game t6          # Black Ops 2 only
--game all         # Both games (default)

# Use custom config
--config /path/to/config.json
```

Search commands additionally support:

```bash
# JSON output
-j, --json

# Limit results
-n, --max 50

# Show context lines
-c, --context 3

# Files only (no content)
-f, --files-only
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

**Settings:**
- `max_results` - Default result limit for searches (can be overridden with `--max`)
- `context_lines` - Lines of context around matches (default: 3, but TUI shows 2 before/after)
- `default_game` - Preferred game when one is not specified

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

## Development

```bash
# Run tests
go test ./...

# Build for current platform
go build -o gscex ./cmd/gscex/

# Build all platforms (requires GoReleaser or build script)
./scripts/build-release.sh v1.0.0
```

## Documentation

- [Architecture](docs/01_architecture.md) - How gscex works
- [Installation](docs/02_installation.md) - Detailed setup
- [Usage Guide](docs/03_usage.md) - Complete command reference
- [Stock Scripts](docs/04_stock_scripts.md) - Understanding T5/T6 scripts
- [Search Features](docs/05_search_features.md) - Advanced search
- [Development](docs/06_development.md) - Contributing guide

## Stock Scripts

This tool indexes the official Plutonium stock scripts:

**T6 (Black Ops 2):**
- 1,469+ files across MP and ZM modes
- 10,950+ functions indexed
- 10,246+ method calls tracked

**T5 (Black Ops 1):**
- Complete multiplayer and zombie scripts
- Full function and method indexing

Search works across all indexed games by default. Use `--game t5` or `--game t6` to filter.

## License

MIT License - See [LICENSE](LICENSE) for details.

---

**Note:** This tool is for educational and modding purposes. Not affiliated with Activision or Treyarch.

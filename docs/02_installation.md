# Installation

**Description:** Step-by-step guide to install and configure the gscex-cli tool.

**Goal:** Get the tool running on your system with minimal effort.

**Applicable for:** T6 modders ready to start using the reference tool.

---

## Index

| File | Topic | Purpose |
|------|-------|---------|
| [INDEX.md](INDEX.md) | Main Index | Project overview and navigation |
| [01_architecture.md](01_architecture.md) | Architecture | How the tool works |
| [03_usage.md](03_usage.md) | Usage Guide | Command reference |
| [04_stock_scripts.md](04_stock_scripts.md) | Stock Scripts | Stock script structure |
| [05_search_features.md](05_search_features.md) | Search Features | Search capabilities |
| [06_development.md](06_development.md) | Development | Building from source |

---

## Prerequisites

- **Go 1.21+** - [Download from golang.org](https://golang.org/dl/)
- **Git** - For cloning repositories
- **2-5MB disk space** - For the index and cached scripts

---

## Installation Methods

### Method 1: Pre-built Binary (Recommended)

Download the latest release:

```bash
# Linux/macOS
curl -L https://github.com/yourusername/gscex-cli/releases/latest/download/gscex-linux-amd64 -o gscex
chmod +x gscex
sudo mv gscex /usr/local/bin/

# Windows
# Download gscex-windows-amd64.exe from releases page
# Add to PATH or use directly
```

### Method 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/gscex-cli.git
cd gscex-cli

# Build
go build -o gscex cmd/gscex/main.go

# Install to $GOPATH/bin
go install ./cmd/gscex

# Or move to system PATH
sudo mv gscex /usr/local/bin/
```

### Method 3: Go Install

```bash
go install github.com/yourusername/gscex-cli/cmd/gscex@latest
```

---

## First Time Setup

After installation, run the init command to download stock scripts and build the index:

```bash
gscex init
```

This will:
1. Clone the plutoniummod/t6-scripts repository
2. Parse all GSC files and build an index
3. Store everything in `~/.gscex/`

Expected output:
```
📦 GSC Reference CLI Setup
━━━━━━━━━━━━━━━━━━━━━━━━━━
Downloading stock scripts... ✓
Parsing GSC files... ✓
Building index... ✓
Index size: 4.2MB
Total functions: 3,847
Total methods: 12,293
━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ Ready! Run 'gscex --help' to see available commands.
```

---

## Updating

To update the stock scripts and index:

```bash
gscex update
```

This refreshes the local copy of t6-scripts and rebuilds the index.

---

## Configuration

Config file: `~/.gscex/config.json`

Default configuration:
```json
{
  "scripts_repo": "https://github.com/plutoniummod/t6-scripts",
  "scripts_branch": "main",
  "cache_dir": "~/.gscex",
  "auto_update": false,
  "max_results": 20,
  "context_lines": 3
}
```

---

## Verification

Test your installation:

```bash
# Check version
gscex version

# Test search
gscex func init

# List commands
gscex --help
```

---

## Troubleshooting

**"command not found"**
- Ensure the binary is in your PATH
- Try: `export PATH=$PATH:/usr/local/bin`

**"failed to clone repository"**
- Check internet connection
- Verify Git is installed: `git --version`
- Check if ~/.gscex directory has write permissions

**"index not found"**
- Run: `gscex init`
- If still failing, clear cache: `rm -rf ~/.gscex` and re-run init

---

*Next: Read [03_usage.md](03_usage.md) to learn how to use the tool.*

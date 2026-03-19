# GSC Explorer (gscex)

**Description:** An interactive explorer and search tool for Plutonium T6 GSC scripts. Browse, search, and navigate stock scripts with an intuitive TUI or command-line interface.

**Goal:** Provide a fast, interactive, and comprehensive way to explore the T6 GSC codebase without manual file browsing.

**Applicable for:** T6/BO2 GSC modders who want to quickly explore, search, and understand stock functions and code patterns.

---

## Index

| File | Topic | Purpose |
|------|-------|---------|
| [01_architecture.md](01_architecture.md) | Architecture | How the tool works, components, and data flow |
| [02_installation.md](02_installation.md) | Installation | How to install and set up the tool |
| [03_usage.md](03_usage.md) | Usage Guide | Command reference and examples |
| [04_stock_scripts.md](04_stock_scripts.md) | Stock Scripts | How stock scripts are indexed and searched |
| [05_search_features.md](05_search_features.md) | Search Features | Advanced search capabilities |
| [06_development.md](06_development.md) | Development | Building and contributing to the tool |

---

## Quick Start

```bash
# Initialize the tool
gscex init

# Launch interactive TUI browser
gscex tui

# Or use CLI commands:
# Search for a function
gscex search "magic bullet"

# Find function definition
gscex func "player_damage"

# Show usage examples
gscex method player give_weapon

# List all files containing a pattern
gscex files "callback"
```

---

## Learning Path

1. **Start here:** Read [02_installation.md](02_installation.md) to get the tool running
2. **For interactive use:** Try the TUI with `gscex tui`
3. **Next:** Learn the CLI basics in [03_usage.md](03_usage.md)
4. **Then:** Understand search capabilities in [05_search_features.md](05_search_features.md)
5. **Reference:** Use [04_stock_scripts.md](04_stock_scripts.md) to understand the stock script structure
6. **Advanced:** Read [01_architecture.md](01_architecture.md) and [06_development.md](06_development.md) to customize

---

## Prerequisites

- Go 1.21+ installed
- Basic understanding of GSC syntax
- Familiarity with command-line interfaces
- Git installed (for downloading stock scripts)

---

*Project Status: In Progress - Initial Setup*

# Architecture

**Description:** Technical architecture and component design of the gscex-cli tool.

**Goal:** Understand how the tool indexes stock scripts, performs searches, and returns results.

**Applicable for:** Developers who want to understand or extend the tool.

---

## Index

| File | Topic | Purpose |
|------|-------|---------|
| [INDEX.md](INDEX.md) | Main Index | Project overview and navigation |
| [02_installation.md](02_installation.md) | Installation | Setup instructions |
| [03_usage.md](03_usage.md) | Usage Guide | Command reference |
| [04_stock_scripts.md](04_stock_scripts.md) | Stock Scripts | Stock script structure |
| [05_search_features.md](05_search_features.md) | Search Features | Search capabilities |
| [06_development.md](06_development.md) | Development | Building from source |

---

## Overview

The GSC Reference CLI is a Go-based command-line tool that provides instant search capabilities over Plutonium T5 and T6 stock GSC scripts.

```
┌─────────────────┐
│   User Input    │
│  gscex search │
│  "magic bullet" │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  CLI Parser     │
│  (cobra)        │
│  --game flag    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Search Engine  │
│  - Function     │
│  - Method       │
│  - Text         │
│  - Multi-game   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Index Storage  │
│  (JSON cache)   │
│  - index-t5.json│
│  - index-t6.json│
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   T5 Scripts    │
│  (plutoniummod  │
│  /t5-scripts)   │
├─────────────────┤
│   T6 Scripts    │
│  (plutoniummod  │
│  /t6-scripts)   │
└─────────────────┘
```

---

## Components

### 1. Indexer

Parses GSC files and builds a searchable index containing:
- Function definitions (with signatures)
- Method calls on entities
- DVARs and level variables
- Include statements
- Code context (surrounding lines)

**Index format (JSON):**
```json
{
  "functions": {
    "magic_bullet": {
      "file": "maps/mp/_utility.gsc",
      "line": 127,
      "signature": "magic_bullet(weapon, start, end, owner, damage)",
      "context": ["// Spawn bullet damage", "magic_bullet(...)"]
    }
  },
  "methods": {
    "player.give_weapon": [...],
    "player.take_all_weapons": [...]
  }
}
```

### 2. Search Engine

Three search modes:

**Function Search** (`gscex func <name>`)
- Exact match on function names
- Shows definition + all usages
- Displays signature and context

**Method Search** (`gscex method <entity> <name>`)
- Searches for entity method calls
- Example: `gscex method player give_weapon`

**Text Search** (`gscex search <pattern>`)
- Fuzzy/substring matching
- Searches across all code
- Shows file, line, and context

### 3. CLI Interface

Commands:
- `init` - Download and index stock scripts
- `search <pattern>` - Search for text across all scripts
- `func <name>` - Find function definition and usages
- `method <entity> <name>` - Find method calls
- `files <pattern>` - List files containing pattern
- `update` - Refresh stock scripts and rebuild index

### 4. Storage

**Local cache:**
- `~/.gscex/` - Config and index storage
- `~/.gscex/scripts-t5/` - Downloaded Black Ops 1 stock scripts
- `~/.gscex/scripts-t6/` - Downloaded Black Ops 2 stock scripts
- `~/.gscex/index-t5.json` - T5 search index
- `~/.gscex/index-t6.json` - T6 search index
- `~/.gscex/config.json` - Configuration (includes game definitions)

**Multi-game support:**
- Supports both T5 (Black Ops 1) and T6 (Black Ops 2)
- Separate indices and script directories per game
- `--game` flag filters results to specific game
- Default behavior searches all indexed games

---

## Data Flow

```
1. User runs command
        ↓
2. CLI parses arguments
        ↓
3. Load index from disk
        ↓
4. Execute search query
        ↓
5. Format and display results
```

---

## GSC Parsing Strategy

GSC is a C-like language with these patterns:

**Functions:**
```gsc
function_name(param1, param2)
{
    // body
}
```

**Methods (entity calls):**
```gsc
entity function_name(params)
```

**Includes:**
```gsc
#include maps/mp/_utility;
```

**DVARs:**
```gsc
getDvar("sv_maxclients")
setDvar("g_speed", 190)
```

The indexer uses regex patterns to extract:
- Function definitions: `/^([a-zA-Z_][a-zA-Z0-9_]*)\s*\([^)]*\)\s*\{/`
- Method calls: `/([a-zA-Z_][a-zA-Z0-9_]*)\.([a-zA-Z_][a-zA-Z0-9_]*)\s*\(/`
- Includes: `/^#include\s+([^;]+);/`

---

## Performance Considerations

- **Indexing:** One-time operation (~5-10 seconds for full repo)
- **Search:** O(1) to O(n) depending on query type
- **Memory:** Index stays in memory during session
- **Disk:** ~2-5MB index file

---

*Next: Read [02_installation.md](02_installation.md) to set up the tool.*

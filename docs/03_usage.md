# Usage Guide

**Description:** Complete reference for all gscex-cli commands with examples.

**Goal:** Learn to effectively search T6 stock scripts for functions, methods, and examples.

**Applicable for:** T6 modders using the reference tool day-to-day.

---

## Index

| File | Topic | Purpose |
|------|-------|---------|
| [INDEX.md](INDEX.md) | Main Index | Project overview and navigation |
| [01_architecture.md](01_architecture.md) | Architecture | How the tool works |
| [02_installation.md](02_installation.md) | Installation | Setup instructions |
| [04_stock_scripts.md](04_stock_scripts.md) | Stock Scripts | Stock script structure |
| [05_search_features.md](05_search_features.md) | Search Features | Search capabilities |
| [06_development.md](06_development.md) | Development | Building from source |

---

## Command Reference

### `init` - Initialize the Tool

Downloads stock scripts and builds the search index.

```bash
gscex init
```

**Options:**
- `--force` - Re-download even if already exists

---

### `tui` - Launch Interactive Browser

Opens a full-featured terminal UI for browsing GSC scripts interactively. This is the recommended way to explore the codebase.

```bash
gscex tui
```

**Features:**
- Interactive search with live results
- File browser with full file tree
- Code preview with line numbers
- **Current line highlighting** - selected result is automatically centered and highlighted
- Smooth scrolling with keyboard controls
- Multiple search modes (text, function, method, files)
- Keyboard-driven navigation
- In-app help (? key)

**Keyboard Shortcuts:**

| Key | Action |
|-----|--------|
| `?` | Toggle help |
| `Tab` | Switch search mode |
| `Enter` | Search / Open result |
| `Ctrl+N` | Next result (results mode) / Scroll down (preview mode) |
| `Ctrl+P` | Previous result (results mode) / Scroll up (preview mode) |
| `Ctrl+F` | Focus search |
| `Ctrl+B` | Back to results |
| `↑/↓` or `j/k` | Scroll in preview |
| `PgUp/PgDn` | Page scroll in preview (10 lines) |
| `q` or `Esc` | Go back / Quit |
| `Ctrl+C` | Quit |

**Navigation Flow:**
1. Start in search mode - type query and press Enter
2. Results appear - navigate with Ctrl+N/P or arrow keys
3. Press Enter on a result to open file preview
4. Preview opens centered on the matching line (highlighted)
5. Scroll with arrows, PgUp/PgDn, or Ctrl+N/P
6. Press q to go back to results
7. Press q again to go back to search

**Search Modes:**
- **text** - Search for any text in all files (up to 100 results)
- **func** - Search for function definitions with **async fuzzy matching** - searches through all 10,950+ functions and returns up to 100 matches (supports "giveweapon" → "give_weapon")
- **method** - Search for method calls (format: "entity method")
- **files** - Search for filenames

---

### `search` - Search for Text

Finds any text across all stock scripts. Results are sorted by filename and line number.

```bash
gscex search "magic bullet"
gscex search "callback"
gscex search "damage"
```

**Output format:**
```
3 results for "magic bullet"

1. maps/mp/_utility.gsc:127
     magic_bullet(weapon, start, end, owner, damage)
     
2. maps/mp/gametypes/_damage.gsc:89
     magic_bullet("none", self.origin, target.origin, self, damage)

3. maps/mp/zombies/_zm.gsc:2341
     magic_bullet(weapon, start, end, undefined, damage)
```

**Options:**
- `-c, --context` - Number of context lines (default: 3)
- `-n, --max` - Maximum results (default: 20)
- `-f, --files-only` - Show only file names
- `-j, --json` - Output as JSON

**JSON Output Example:**
```bash
gscex search "magic bullet" --json
```

```json
[
  {
    "file": "maps/mp/_utility.gsc",
    "line": 127,
    "content": "    magic_bullet(weapon, start, end, owner, damage)",
    "context": [...]
  }
]
```

---

### `func` - Find Function

Searches for function definitions and usages. Results are sorted by filename and line number.

```bash
gscex func init
gscex func give_weapon
gscex func player_damage
```

**Output format:**
```
Function: init

Definition:
  maps/mp/gametypes/_callbacksetup.gsc:45
  init()
      level.callbacksetup = true;
      level thread onPlayerConnect();

Usages (14 found):
  maps/mp/gametypes/_rank.gsc:23
  maps/mp/gametypes/_hud.gsc:15
  maps/mp/zombies/_zm.gsc:89
  ... (11 more)
```

**Note:** Function search supports fuzzy matching. You can search without underscores (e.g., `giveweapon` will match `give_weapon`), and it will find the first matching function.

**Options:**
- `-d, --definition-only` - Show only definition
- `-j, --json` - Output as JSON

**JSON Output Example:**
```bash
gscex func init --json
```

```json
{
  "function": {
    "name": "init",
    "file": "maps/mp/gametypes/_callbacksetup.gsc",
    "line": 45,
    "signature": "init()",
    "context": [...]
  },
  "usages": [...]
}
```

---

### `method` - Find Method Calls

Searches for method calls on specific entities. Results are sorted by filename and line number.

```bash
# Find all player method calls named "give_weapon"
gscex method player give_weapon

# Find zombie-related methods
gscex method zombie spawn

# Find hud methods
gscex method hud setText
```

**Output format:**
```
Method: player.give_weapon

Usages (27 found):
  maps/mp/zombies/_zm_weapons.gsc:156
     player give_weapon(weapon_name, ...)
  
  maps/mp/zombies/_zm_perks.gsc:234
     player give_weapon(perk_bottle, ...)
```

**Options:**
- `-j, --json` - Output as JSON

---

### `files` - List Files

Finds all files containing a pattern. Results are sorted alphabetically.

```bash
gscex files "damage"
gscex files "callback"
gscex files "zombie"
```

**Output format:**
```
8 files contain "damage"

maps/mp/gametypes/_damage.gsc
maps/mp/zombies/_zm_damage.gsc
maps/mp/zombies/_zm_laststand.gsc
...
```

**Options:**
- `-j, --json` - Output as JSON array

---

### `update` - Update Index

Refreshes stock scripts and rebuilds the index.

```bash
gscex update
```

---

### `info` - Show Tool Info

Displays tool version, index statistics, and config.

```bash
gscex info
gscex info --json
```

**Output:**
```
GSC Reference CLI v1.0.0

Index Stats:
  Total files: 892
  Total functions: 3,847
  Total methods: 12,293
  Index path: ~/.gscex/index.json

Config:
  Scripts repo: plutoniummod/t6-scripts
  Cache dir: ~/.gscex
```

**Options:**
- `-j, --json` - Output as JSON

---

## Common Workflows

### Find How to Use a Function

```bash
# Search for the function
gscex func magic_bullet

# See how it's used in practice
gscex search "magic_bullet"
```

### Learn Player Methods

```bash
# Find all player method calls
gscex search "player " -c 1

# Find specific method
gscex method player take_all_weapons
```

### Find Callbacks

```bash
# Find all callback registrations
gscex search "add_callback"

# Find specific callback
gscex search "callback_player_connect"
```

### Explore Weapon System

```bash
# Find weapon-related functions
gscex files "weapon"
gscex func give_weapon
gscex func take_all_weapons
```

### Export Results for Processing

```bash
# Export search results to JSON for processing in other tools
gscex search "weapon" --json > weapons.json

# Export function info
gscex func init --json > init_function.json

# Export file list
gscex files "_zm" --json > zombie_files.json
```

---

## Tips

1. **Use quotes** for multi-word searches: `gscex search "magic bullet"`
2. **Start broad, then narrow**: Use `search` first, then `func` or `method`
3. **Check context**: Use `-c 5` to see more surrounding code
4. **Find examples**: Search for function names to see real usage patterns
5. **JSON output**: Use `--json` for programmatic processing or piping to other tools
6. **Sorted results**: All output is sorted by filename and line number for easy reading

---

*Next: Read [04_stock_scripts.md](04_stock_scripts.md) to understand the stock script organization.*

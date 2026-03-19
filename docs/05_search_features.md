# Search Features

**Description:** Advanced search capabilities and techniques for the gscex-cli tool.

**Goal:** Master the search features to find code faster and more accurately.

**Applicable for:** T6 modders who want to maximize their productivity with the tool.

---

## Index

| File | Topic | Purpose |
|------|-------|---------|
| [INDEX.md](INDEX.md) | Main Index | Project overview and navigation |
| [01_architecture.md](01_architecture.md) | Architecture | How the tool works |
| [02_installation.md](02_installation.md) | Installation | Setup instructions |
| [03_usage.md](03_usage.md) | Usage Guide | Command reference |
| [04_stock_scripts.md](04_stock_scripts.md) | Stock Scripts | Stock script structure |
| [06_development.md](06_development.md) | Development | Building from source |

---

## Search Types

### 1. Exact Function Search

Finds function definitions by exact name match.

```bash
gscex func init
gscex func onPlayerConnect
```

**Behavior:**
- Case-insensitive
- Matches function definitions only
- Shows definition + all call sites

---

### 2. Method Entity Search

Finds method calls on specific entity types.

```bash
# Find player.give_weapon calls
gscex method player give_weapon

# Find all player methods
gscex search "player " | grep -E "\.\w+\("
```

**Entity Types:**
- `player` - Player entities
- `level` - Global level object
- `self` - Context-dependent entity
- `hud` - HUD elements
- `zombie` - Zombie entities

---

### 3. Fuzzy Text Search

Substring matching across all script content.

```bash
# Find anything containing "damage"
gscex search "damage"

# Find HUD-related code
gscex search "hud"

# Find specific patterns
gscex search "callback_player"
```

---

### 4. File Search

Finds which files contain a pattern.

```bash
# What files handle damage?
gscex files "damage"

# Find weapon-related files
gscex files "weapon"
```

---

## Advanced Techniques

### Combining Commands

```bash
# Find damage functions and see their definitions
for func in $(gscex search "function.*damage" -f | head -5); do
    gscex func "$func"
done
```

### Finding Related Code

```bash
# Find where a function is defined
gscex func magic_bullet

# Then see how it's used
gscex search "magic_bullet"

# Find the include path
gscex search "#include.*utility"
```

### Context Expansion

```bash
# See more context around matches
gscex search "give_weapon" -c 5

# See full function bodies
gscex func give_weapon -c 50
```

---

## Search Tips

### Tip 1: Start Broad

```bash
# Don't know exact function name?
gscex search "weapon.*give"

# Find patterns
gscex search "thread.*update"
```

### Tip 2: Use File Filter

```bash
# Search only zombies code
gscex search "perk" | grep "zombies"

# Search only gametypes
gscex search "callback" | grep "gametypes"
```

### Tip 3: Learn by Example

```bash
# Find usage patterns
gscex search "waittill" | head -20

# See real implementations
gscex search "foreach.*player"
```

### Tip 4: Find Callbacks

```bash
# All callbacks
gscex search "callback_"

# Specific callback
gscex func callback_player_connect
```

---

## Common Search Patterns

### Player-Related
```bash
gscex method player damage
gscex search "player.origin"
gscex search "player thread"
```

### Weapon-Related
```bash
gscex func give_weapon
gscex func take_all_weapons
gscex search "weapon.*ammo"
```

### HUD-Related
```bash
gscex search "hud ="
gscex search "hud setText"
gscex files "_hud"
```

### Zombies-Related
```bash
gscex files "_zm"
gscex func round_start
gscex search "zombie_spawn"
```

### DVAR-Related
```bash
gscex search "getDvar"
gscex search "setDvar"
gscex search "getDvarInt"
```

---

## Performance Tips

- **Use specific queries** - `func` is faster than `search` for finding functions
- **Limit results** - Use `-n 10` for quick previews
- **Cache is fast** - First search loads index, subsequent searches are instant

---

*Next: Read [06_development.md](06_development.md) to contribute or customize the tool.*

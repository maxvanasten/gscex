# Stock Scripts Reference

**Description:** Overview of the Plutonium T6 stock GSC scripts organization and structure.

**Goal:** Understand how stock scripts are organized to search them effectively.

**Applicable for:** T6 modders wanting to understand the script hierarchy.

---

## Index

| File | Topic | Purpose |
|------|-------|---------|
| [INDEX.md](INDEX.md) | Main Index | Project overview and navigation |
| [01_architecture.md](01_architecture.md) | Architecture | How the tool works |
| [02_installation.md](02_installation.md) | Installation | Setup instructions |
| [03_usage.md](03_usage.md) | Usage Guide | Command reference |
| [05_search_features.md](05_search_features.md) | Search Features | Search capabilities |
| [06_development.md](06_development.md) | Development | Building from source |

---

## Repository Structure

The T6 stock scripts are organized into two main modes:

```
t6-scripts/
├── MP/                          # Multiplayer scripts
│   ├── maps/mp/gametypes/      # Game modes
│   ├── maps/mp/_utility.gsc    # Core utilities
│   ├── maps/mp/_rank.gsc       # Ranking/XP
│   └── ...
│
└── ZM/                          # Zombies scripts
    ├── maps/mp/zombies/        # Core zombie systems
    ├── maps/mp/_zm.gsc         # Main zombie script
    ├── maps/mp/gametypes/      # Zombies gametypes
    └── ...
```

---

## Key Files by Category

### Core Utilities

| File | Purpose | Common Functions |
|------|---------|------------------|
| `maps/mp/_utility.gsc` | General utilities | `magic_bullet`, `waittill_notify_or_timeout` |
| `maps/mp/_rank.gsc` | Ranking/XP | `incRankXP`, `getRankInfo` |
| `maps/mp/_hud_util.gsc` | HUD utilities | `setPoint`, `setShader` |

### Multiplayer

| File | Purpose | Common Functions |
|------|---------|------------------|
| `maps/mp/gametypes/_callbacksetup.gsc` | Callbacks | `init`, `onPlayerConnect` |
| `maps/mp/gametypes/_damage.gsc` | Damage handling | `player_damage`, `damage_callback` |
| `maps/mp/gametypes/_weapons.gsc` | Weapon logic | `give_weapon`, `take_all_weapons` |

### Zombies

| File | Purpose | Common Functions |
|------|---------|------------------|
| `maps/mp/zombies/_zm.gsc` | Core zombie logic | `init`, `round_start` |
| `maps/mp/zombies/_zm_weapons.gsc` | Weapon box/PaP | `weapon_give`, `weapon_take` |
| `maps/mp/zombies/_zm_perks.gsc` | Perk machines | `give_perk`, `perk_set_drinking` |
| `maps/mp/zombies/_zm_powerups.gsc` | Power-ups | `powerup_setup`, `powerup_grab` |
| `maps/mp/zombies/_zm_score.gsc` | Points system | `add_to_player_score`, `minus_to_player_score` |
| `maps/mp/zombies/_zm_utility.gsc` | Zombie utilities | `wait_for_zombie_spawn`, `is_player_valid` |

---

## GSC Syntax Basics

### Function Definition
```gsc
function_name(param1, param2)
{
    // function body
    return value;
}
```

### Method Call (entity function)
```gsc
entity method_name(params)
```

### Includes
```gsc
#include maps/mp/_utility;
#include common_scripts/utility;
```

### Common Entity Types
- `player` - Player entity
- `level` - Global game state
- `self` - Current entity (context-dependent)
- `zombie` - Zombie entity
- `trigger` - Trigger entities

---

## Search Strategies

### Find Core Functions
```bash
gscex func init  # Start with init to see entry points
gscex func main  # Main game loop
```

### Find Callbacks
```bash
gscex search "add_callback"
gscex search "onPlayer"
```

### Find Utility Functions
```bash
gscex files "_utility"
gscex func waittill
```

### Find Weapon Code
```bash
gscex files "weapon"
gscex func give_weapon
```

---

*Next: Read [05_search_features.md](05_search_features.md) for advanced search techniques.*

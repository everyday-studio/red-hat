# SOMNIUM — Implementation Status

> **This document changes frequently.**
> Update this file every time a script or scene is created, completed, or changed in status.
> The agent must read this file at the start of every new task to avoid recreating existing work.

---

## Scripts

| File | Status | Notes |
|------|--------|-------|
| `scripts/game/game_state.gd` | ✅ Done | AutoLoad. Steam init, JWT/steam_id/player_name storage, session clear |
| `scripts/network/api_client.gd` | ✅ Done | AutoLoad. All HTTP wrapped. JWT auth header auto-injected |
| `scripts/network/steam_auth.gd` | ✅ Done | Full Steam ticket → POST /auth/steam → JWT flow |
| `scripts/network/matchmaking.gd` | 🔧 Stub | Signals defined. `join_queue()` and `_poll_status()` bodies empty — **blocked by backend** |
| `scripts/game/game_manager.gd` | 🔧 Stub | Constants/signals defined. Timer logic, gate open, match end all empty — **blocked by backend** |
| `scripts/game/player_controller.gd` | 🔧 Partial | Movement works. `_request_fire()` RPC body empty |
| `scripts/game/anima_system.gd` | 🔧 Partial | RPC drop-own-color guard implemented. `_apply_orb_color()` is TODO (needs AnimaOrb sprite node) |
| `scripts/game/gun_system.gd` | 🔧 Stub | File exists, logic not yet implemented |
| `scripts/game/task_system.gd` | 🔧 Stub | File exists, logic not yet implemented |
| `scripts/game/proximity_system.gd` | 🔧 Stub | Constants/signals defined. `_physics_process` and `get_nearby_players()` empty |
| `scripts/game/day_night_system.gd` | ❌ Not created | Controls day/night phase timer (2 min day / 1 min night), emits signals for phase change |
| `scripts/game/monster_controller.gd` | ❌ Not created | Individual monster AI: sleep/wake state, patrol, chase, aggro. Day = dormant unless woken |
| `scripts/game/ghost_system.gd` | ❌ Not created | Ghost roam, monster wake/distract, RED possession (exclusive, time-limited, cooldown) |
| `scripts/game/item_system.gd` | ❌ Not created | Item spawn, pickup, carry limit, use. Feeds into individual item effects |
| `scripts/game/portal_system.gd` | ❌ Not created | Escape portal: unlock on all tasks done, 10-sec channeling, non-SILVER instant death |
| `scripts/ui/character_select.gd` | ❌ Not created | Character + customization preset selection UI before match |
| `scripts/ui/lobby_controller.gd` | 🔧 Stub | `_on_find_match_pressed()` empty. No player list or ready UI yet |
| `scripts/ui/hud_controller.gd` | 🔧 Stub | File exists, logic not yet implemented |
| `scripts/ui/cctv_controller.gd` | 🔧 Stub | File exists, logic not yet implemented |

---

## Scenes

| File | Status | Notes |
|------|--------|-------|
| `scenes/ui/steam_auth_test.tscn` | ✅ Done | Dev-only test scene. Not part of the final game flow |
| `scenes/ui/lobby.tscn` | ❌ Not created | Pre-game lobby waiting room |
| `scenes/ui/hud.tscn` | ❌ Not created | In-game overlay (timer, task counter, chat) |
| `scenes/game/player.tscn` | ❌ Not created | Player character with sprite, Anima Orb, collision |
| `scenes/game/anima_orb.tscn` | ❌ Not created | Anima Orb sprite with color modulation |
| `scenes/game/world.tscn` | ❌ Not created | Map root scene (tilemap, spawn points, gate) |
| `scenes/result/result_screen.tscn` | ❌ Not created | End-game role reveal and winner announcement |
| `scenes/game/monster.tscn` | ❌ Not created | Monster sprite + CollisionShape + AI state machine |
| `scenes/game/item_pickup.tscn` | ❌ Not created | Interactable item on the map (sprite + Area2D pickup trigger) |

---

## What Can Be Built Without the Backend

The following can be fully implemented and tested **right now**, without any server:

- `scenes/game/player.tscn` — sprite, movement, collision shape
- `scenes/game/anima_orb.tscn` — orb color modulation sprite
- `scenes/game/world.tscn` — placeholder map, spawn points
- `scenes/ui/lobby.tscn` — layout, player list (local mock), ready button
- `scenes/ui/hud.tscn` — timer display, task counter layout, chat bubble layout
- `scenes/result/result_screen.tscn` — role reveal layout with placeholder data
- `scripts/game/proximity_system.gd` — pure local distance calculation, no server needed
- `scripts/ui/hud_controller.gd` — display logic only (reads from GameState)
- `scripts/ui/cctv_controller.gd` — camera switching logic (local)
- `scripts/game/gun_system.gd` — client-side visual only (muzzle flash, bullet sprite, SFX)
- `scripts/game/day_night_system.gd` — pure local timer + phase signal, no server needed
- `scripts/game/monster_controller.gd` — sleep/wake state + movement AI, fully local
- `scripts/game/ghost_system.gd` — ghost visuals, monster distract, CRIMSON possession UI (server decides outcome)
- `scripts/game/item_system.gd` — local pickup collision + inventory display
- `scripts/game/portal_system.gd` — channeling timer + visual, unlock logic driven by task_system
- `scenes/game/monster.tscn` — sprite, collision, AI state machine
- `scenes/game/item_pickup.tscn` — sprite + pickup Area2D
- `scripts/ui/character_select.gd` — local preset picker, saves to GameState

---

## What Is Blocked by the Backend

Do not implement these until the backend is ready:

- `scripts/network/matchmaking.gd` — needs `POST /matchmaking/join` and `GET /matchmaking/status`
- ENet peer connection to headless server — needs IP/port from matchmaking response
- `scripts/game/game_manager.gd` timer sync — needs server-driven timer events
- `scripts/game/task_system.gd` server sync — task state is server-authoritative
- `scripts/game/anima_system.gd` RPC calls — server sends orb colors after role assignment

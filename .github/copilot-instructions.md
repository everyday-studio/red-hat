# RED HAT — Copilot Agent Instructions

> **This document is the project constitution.**
> The agent must follow all rules defined here unconditionally.
> The agent may only modify this document when the user **explicitly instructs** a change.
> When a modification is made, the agent must immediately brief the user on exactly what was changed and why.

---

## 0. Agent Role & Boundaries

- The user is the **client-side developer** and works exclusively inside `apps/client/`.
- The agent must **never read, modify, or reference** files inside `apps/server/`. That directory is owned by the backend developer (a separate person).
- If the agent needs information about a server API (endpoint signature, request/response shape, etc.) and cannot determine it from this document, it must **stop immediately and ask the user** before proceeding.
- The agent must **never make assumptions alone**. If something is unclear, unknown, or missing — stop, ask the user, then continue.

---

## 1. Game Design Specification

> **This section is maintained in a separate file.**
> Before implementing any gameplay feature, read [`.github/GAME_SPEC.md`](.github/GAME_SPEC.md) for rules, UX flows, mechanics, and art direction.
> Update that file whenever game design decisions are made or changed.

---

## 2. Repository Structure

```
project-root/
├── .github/
│   └── copilot-instructions.md    ← this file
├── apps/
│   ├── client/                    ← Godot 4 project (USER'S TERRITORY — agent works here only)
│   │   ├── scenes/
│   │   │   ├── ui/
│   │   │   │   ├── lobby.tscn
│   │   │   │   └── hud.tscn
│   │   │   ├── game/
│   │   │   │   ├── player.tscn
│   │   │   │   ├── helmet.tscn
│   │   │   │   └── world.tscn
│   │   │   └── result/
│   │   │       └── result_screen.tscn
│   │   ├── scripts/
│   │   │   ├── network/
│   │   │   │   ├── api_client.gd      ← ALL HTTP calls go through here only
│   │   │   │   ├── matchmaking.gd
│   │   │   │   └── steam_auth.gd
│   │   │   ├── game/
│   │   │   │   ├── game_state.gd      ← AutoLoad singleton
│   │   │   │   ├── game_manager.gd
│   │   │   │   ├── player_controller.gd
│   │   │   │   ├── helmet_system.gd   ← SERVER AUTHORITY ONLY
│   │   │   │   ├── gun_system.gd
│   │   │   │   ├── task_system.gd
│   │   │   │   └── proximity_system.gd
│   │   │   └── ui/
│   │   │       ├── hud_controller.gd
│   │   │       ├── lobby_controller.gd
│   │   │       └── cctv_controller.gd
│   │   ├── assets/
│   │   │   ├── characters/
│   │   │   ├── helmets/
│   │   │   ├── maps/
│   │   │   └── sounds/
│   │   │       └── shared/
│   │   │           └── shot_or_explosion.wav  ← ONE FILE for both events (intentional)
│   │   └── addons/
│   │       └── godotsteam/
│   └── server/                    ← Go backend (DO NOT TOUCH — backend developer's territory)
└── .gitignore
```

---

## 3. Tech Stack

### Client (`apps/client/`) — Agent's working scope

| Layer | Technology |
|-------|-----------|
| Engine | Godot 4 (stable, Standard — NOT .NET) |
| Language | GDScript |
| IDE | VSCode + godot-tools extension |
| Multiplayer | Godot built-in ENet (UDP) via MultiplayerAPI |
| Steam Integration | GodotSteam plugin (Godot 4 compatible) |
| Voice (proximity) | Distance-based AudioStreamPlayer volume attenuation |
| Art | 2D pixel sprites, top-down |

### Server (`apps/server/`) — Do not touch

| Layer | Technology |
|-------|-----------|
| Language | Go |
| Database | PostgreSQL |
| Cache / Queue | Redis |
| Auth | Steam Web API ticket validation → custom JWT |
| Orchestration | Docker (headless Godot containers) |
| Infrastructure | AWS EC2 |

### Client ↔ Server API Endpoints

```
POST /auth/steam          Steam ticket → JWT
POST /matchmaking/join    Enter matchmaking queue
GET  /matchmaking/status  Poll for match result (0.5s interval)
GET  /match/history       Fetch player match history

POST /match/result        Headless server → backend ONLY
                          The client must NEVER call this endpoint.
```

### Multiplayer Architecture

The game uses a **backend-managed matchmaking + headless Godot server** model:

1. **Auth**: The client obtains a Steam auth ticket → POSTs it to `/auth/steam` → receives a JWT.
2. **Matchmaking**: The client POSTs `/matchmaking/join` and polls `GET /matchmaking/status` every 0.5s. The **Go backend** handles the matchmaking queue and decides when a group of players forms a match.
3. **Game session**: When a match is formed, the backend spins up a **headless Godot 4 Docker container** on AWS EC2. This container acts as the authoritative game server (ENet UDP).
4. **Client connection**: The backend returns the headless server's IP and port in the matchmaking status response. The client connects via Godot's built-in `ENetMultiplayerPeer`.
5. **During gameplay**: All security-critical logic (helmet color assignment, gun fire resolution, timer sync) runs on the **headless server only**. The client is a thin input sender and render receiver.
6. **Post-game**: The headless server POSTs the match result to `/match/result` on the backend, then shuts down.

> **Important**: Steam Matchmaking / Steam Lobbies are **NOT used**. Matchmaking is entirely handled by the Go backend. GodotSteam is used **only** for Steam authentication (ticket generation) and future Steam features (achievements, overlay, etc.).

---

## 4. Security & Architecture Rules

These rules must **never** be violated.
If generating code would break any of these, the agent must stop and ask the user for direction.

### Rule 1 — Helmet color is SERVER AUTHORITY ONLY

The client must **never** store, receive, read, or branch on its own helmet color.
Helmet color data lives on the server only. The server sends render commands to other players only.

```gdscript
# CORRECT
@rpc("authority", "call_remote", "reliable")
func receive_helmet_color(peer_id: int, color: Color) -> void:
    if peer_id == multiplayer.get_unique_id():
        return  # own color → silently drop, never render
    _apply_helmet_color(peer_id, color)

# WRONG — never write this
var my_color = helmet_color         # client must not hold own color
if my_helmet_color == Color.RED:    # never branch on own color
```

### Rule 2 — Always guard with is_multiplayer_authority()

Every `_physics_process`, `_process`, and `_input` on player-owned nodes must start with this guard:

```gdscript
func _physics_process(delta: float) -> void:
    if not is_multiplayer_authority():
        return
    # ... rest of logic
```

### Rule 3 — Gun logic is SERVER AUTHORITY ONLY

The client requests a fire action. The server decides whether it's a real bullet or self-detonation, because only the server secretly holds the shooter's role.

```gdscript
# CORRECT
@rpc("any_peer", "call_remote", "reliable")
func request_fire() -> void:
    if not multiplayer.is_server():
        return
    _server_resolve_fire(multiplayer.get_remote_sender_id())

# WRONG
if my_role == "hunter":   # client never knows own role
    fire_real_bullet()
```

### Rule 4 — All HTTP calls go through api_client.gd only

Never use `HTTPRequest` directly in scene scripts or game logic files.
Always call `ApiClient.post()` or `ApiClient.get()`.

```gdscript
# CORRECT
var result: Dictionary = await ApiClient.post("/matchmaking/join", {})

# WRONG — never instantiate HTTPRequest outside api_client.gd
var http = HTTPRequest.new()
```

### Rule 5 — Gunshot and explosion share one audio file

This is a **core game mechanic** — players must not be able to distinguish gunshots from self-detonations by sound.

```gdscript
# CORRECT — both events use the same file
func on_gunshot() -> void:
    $AudioPlayer.stream = SHOT_OR_EXPLOSION_SFX

func on_helmet_detonation() -> void:
    $AudioPlayer.stream = SHOT_OR_EXPLOSION_SFX  # intentional, same asset

# WRONG — separate files break the mechanic
func on_gunshot() -> void:
    $AudioPlayer.stream = GUNSHOT_SFX
```

### Rule 6 — Client never calls POST /match/result

Only the headless server (Docker container) calls this endpoint after a game ends.
Generating client code that calls `/match/result` is a bug.

### Rule 7 — Every deliverable must be immediately user-verifiable

Every code change must be testable by the user directly — via mouse, keyboard, or visible on screen in the Godot editor — at the moment it is delivered. Do not deliver "invisible" code that requires manual wiring before anything can be observed.

- **Visual feature** → provide or modify a `.tscn` scene so the user can press F5 and see it running.
- **Non-visual feature** (e.g. network logic, state machine) → add `print()` or `push_warning()` calls so results appear in the Godot **Output** panel.
- **Never** deliver a standalone script with no runnable entry point and say "you can wire it up later."

Remove all debug `print()` statements in the commit step (Step 0 of the next task).

---

## 5. GDScript Coding Conventions

### File Header (required on every new script)

```gdscript
## script_name.gd
## One-line description of what this script does.
##
## Authority: [server | client | both]
## AutoLoad: [yes | no]
extends Node
```

### Naming

```gdscript
# Constants — UPPER_SNAKE_CASE
const MAX_PLAYERS: int = 8
const BASE_TIMER: float = 600.0

# Variables — snake_case, always with explicit type
var detonation_timer: float = BASE_TIMER
var is_dark_phase_active: bool = false

# Functions — snake_case, verb-first
func start_dark_phase() -> void: pass
func get_nearby_players(radius: float) -> Array[Node]: pass
func _on_timer_timeout() -> void: pass   # signal handlers prefix: _on_

# Signals — past-tense or noun phrase
signal player_eliminated(peer_id: int)
signal task_completed(task_id: String)
```

### Type Hints — Always Explicit

```gdscript
# CORRECT
var players: Dictionary = {}
func set_color(peer_id: int, color: Color) -> void: pass

# WRONG
var players = {}
func set_color(peer_id, color): pass
```

### Signal Connection — Callable Syntax Only (Godot 4)

```gdscript
# CORRECT
timer.timeout.connect(_on_timer_timeout)

# WRONG — Godot 3 style, never use
timer.connect("timeout", self, "_on_timer_timeout")
```

### Await Pattern for HTTP

```gdscript
func _do_login() -> void:
    var result: Dictionary = await ApiClient.post("/auth/steam", {
        "ticket": ticket,
        "steam_id": str(Steam.getSteamID())
    })
    if result.has("error"):
        push_error("Login failed: " + str(result["error"]))
        return
    GameState.jwt = result["token"]
```

### Comments

- Write all code comments in **English**.
- Add comments **only for complex or non-obvious logic**. Do not comment self-explanatory lines.
- Use `##` for doc-style comments on functions and classes, `#` for inline logic notes.

```gdscript
## Resolves a fire request from a client peer.
## Runs on the server only. Decides real bullet vs self-detonation
## based on the secretly stored role of the requesting peer.
func _server_resolve_fire(peer_id: int) -> void:
    # roles dict is server-only — never exposed to clients
    var role: String = _roles[peer_id]
    if role == "hunter":
        _fire_real_bullet(peer_id)
    else:
        _detonate_helmet(peer_id)
```

---

## 6. Git Convention

### Branch Naming

- Format: `Type/short-description-with-hyphens`
- The type segment starts with a **capital letter**.
- Use hyphens `-` to separate words. No underscores, no spaces.

```
Feat/steam-auth
Feat/matchmaking-polling
Feat/helmet-led-system
Fix/proximity-timer-double-drain
Fix/cctv-camera-lag
Refactor/api-client-error-handling
Chore/godotsteam-plugin-setup
```

### Commit Message Format

- Format: `type: short-description`
- All **lowercase**, written in **English**, kept **concise**.
- No period at the end.

```
feat: steam auth login flow
feat: matchmaking polling loop
feat: helmet led color rpc sync
fix: proximity timer double drain
fix: gun sfx same as explosion sfx
refactor: centralize http error handling
chore: add godotsteam plugin
docs: update copilot instructions
style: format player_controller
```

### Allowed Commit Types

| Type | When to use |
|------|-------------|
| `feat` | New feature or behavior |
| `fix` | Bug fix |
| `refactor` | Restructure without behavior change |
| `style` | Formatting, whitespace, no logic change |
| `docs` | Documentation only |
| `chore` | Build, config, tooling |
| `perf` | Performance improvement |

### Branch Rules

```
- Always branch from main
  git checkout main && git pull
  git checkout -b Feat/my-feature

- Never push directly to main

- Rebase before PR
  git fetch origin main
  git rebase origin/main
```

---

## 7. Development Workflow Loop (Delayed Commit Protocol)

For every task or feature request, the agent must strictly follow this loop.
**Do NOT propose `git commit` commands for the current task until the user explicitly confirms it works, or implicitly confirms by requesting the next feature.**

### Step 0 — Commit Previous Task (Before Starting New Work)

Before starting a new task, commit the previously verified work.

If multiple changes have accumulated without being committed, **do NOT lump them into one massive commit.** Instead:
1. Analyze all uncommitted changes.
2. Categorize them logically by type and scope.
3. **Brief the user in Korean** with a proposed split strategy.
   > 예: "3개의 커밋으로 분리하는 것을 제안합니다: `feat: matchmaking-polling`, `fix: proximity-timer`, `style: hud-layout`. 진행할까요?"
4. **Wait for the user's approval** before providing the actual `git add` and `git commit` commands.

### Step 1 — Analyze

Understand the exact requirements of the new request.
If anything is unclear or information is missing — **stop and ask the user a specific question** before writing any code.

### Step 2 — Create Issue & Branch

Before writing any code, the agent must guide the user to create a GitHub issue.

1. **Determine the correct issue template** based on the work type:
   - New game feature or mechanic → `✨ Feature Request`
   - Bug fix → `🐛 Bug Report`
   - Dev task, refactor, config → `📋 Task`

2. **Provide a filled-out issue draft** in Korean that the user can paste directly into GitHub:

> 예 (Feature Request):
> ```
> 제목: [Feat] Steam 인증 로그인 흐름 구현
>
> GodotSteam에서 auth ticket을 발급받아 /auth/steam으로 POST한 뒤
> 응답으로 받은 JWT를 GameState에 저장하는 흐름을 구현합니다.
>
> 완료 기준:
> - [ ] Steam 클라이언트 실행 중일 때 티켓 발급 성공 로그 확인
> - [ ] /auth/steam 응답으로 JWT 수신 및 GameState.jwt 저장 확인
> ```

3. After the user confirms the issue is created and provides the issue number, **propose the branch**:

> 예:
> ```
> git checkout -b Feat/steam-auth
> ```

**Do NOT start writing code until the user provides the issue number.**

### Step 3 — Code

Write clean, modular GDScript code following all conventions in this document.
Add comments only on complex or non-obvious logic sections.

### Step 4 — Verify & Wait

Suggest exactly how the user can test the new code in Godot (which scene to run, what to click, what to observe).

**STOP HERE.** End the response by asking the user:
- Does it work as expected?
- Are there any errors?

**Do NOT generate `git commit` commands for this newly written code yet.**

### Step 5 — Push & PR Guide

After the user confirms the code works and the agent has committed, when `git push` is executed the agent must **immediately provide a filled-out PR template** based on the work done.

The current PR template has exactly **3 sections**. Provide a copyable markdown block filling in all three:

```markdown
Closes #[이슈 번호]

## 무엇을 왜 변경했나요?
[변경 내용 및 이유를 자유롭게 — 스크린샷/GIF 첨부 환영]

## 메모 (선택)
[리뷰어에게 전달할 내용, 미완성 사항, 후속 작업 등]
```

Accompany the template with a brief Korean note explaining what the user should review before copy-pasting.

### Error Handling

If the user reports an error during Step 4:
1. Debug and provide the fixed code.
2. Explain what caused the error **in Korean**.
3. Return to Step 4 and wait for confirmation again.

---

## 8. Agent Communication Rules

- **All explanations, plans, status updates, and briefings** must be written **in Korean**.
- **All code and inline comments** inside `.gd` files must be written in **English**.
- The user is **new to Godot and GDScript**. When writing new code or implementing a new feature, the agent must provide a **friendly and detailed briefing in Korean** that explains:
  - What was created and why.
  - How the code works at a high level.
  - Any Godot-specific concepts introduced (nodes, signals, RPC, AutoLoad, scenes, etc.).
  - How to test it step by step.
- When the agent is uncertain about a requirement, missing information, or facing an ambiguous situation — **stop immediately, do not guess, and ask the user a specific and clear question.**

---

## 9. What the Agent Must Never Do

- Modify or reference anything inside `apps/server/`
- Store, read, compare, or render the local player's own helmet color
- Call `POST /match/result` from client code
- Use `HTTPRequest` directly outside of `api_client.gd`
- Use Godot 3 connect syntax: `connect("signal", self, "method")`
- Omit type hints on variables or function signatures
- Create separate audio files for gunshot vs helmet explosion
- Push commits to `main` directly
- Make assumptions when something is unclear — always ask first
- Lump multiple unrelated changes into a single commit without briefing the user
- Self-modify this document without an explicit instruction from the user
- Start writing code before the user has created a GitHub issue and provided the issue number
- Deliver code with no runnable or visible entry point — every deliverable must be immediately testable by the user (Rule 7)

---

## 10. Implementation Status

> **This section is maintained in a separate file.**
> Before starting any new task, read [`.github/IMPLEMENTATION_STATUS.md`](.github/IMPLEMENTATION_STATUS.md) to check what has been built, what is a stub, and what is blocked by the backend.
> Update that file immediately whenever a script or scene changes status.

---

## 11. Scene Node Specifications

> **This section is maintained in a separate file.**
> When creating or modifying any scene, read [`.github/SCENE_SPECS.md`](.github/SCENE_SPECS.md) and follow the node tree exactly.
> `@onready` variable names in scripts must match the node names defined there.
> Update that file whenever a scene's node tree is intentionally changed.

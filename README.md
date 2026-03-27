# SOMNIUM

> **2D pixel-art social deduction extraction game for PC (Steam)**  
> Developed by **everydaystudio**

---

## What is SOMNIUM?

SONMIUM is a reverse-mafia game set inside a shared dream. Players each carry an **Anima Orb** — a glowing soul fragment visible to everyone except the bearer. A shared countdown timer ticks toward detonation. Complete tasks to extend it. Escape through the portal before time runs out.

The twist: the Hunter's gun fires real bullets. The Civilian's gun detonates their own orb. And the **gunshot sound is identical to the detonation sound**.

### Roles

| Role | Anima Orb | Win Condition |
|------|-----------|---------------|
| Hunter | CRIMSON | Eliminate the SILVER before they escape |
| Target | SILVER | Escape through the portal |
| Civilian | ASHEN | SILVER escapes → shared win |

---

## Why a Dedicated Go Backend? Why Not Steam Lobbies?

### The core mechanic requires a server that no player can own

SONMIUM's central rule is that **players cannot know their own Anima Orb color (role)**.  
In a standard Steam Lobby or P2P setup, one player acts as the host — a **listen server** running on someone's machine. Every player's role must exist somewhere in memory for the game to function, which means the host always has direct access to every player's role, including their own.

```
Steam P2P (listen server):
  Player A — HOST  ← all roles exist here in memory
  Player B, C, D   ← synced from host
  → Host can read their own role with a memory viewer. Game over.
```

A **headless dedicated server** running on AWS EC2 is the only architecture that makes this impossible:

```
Headless server (EC2):
  AWS EC2 — authoritative server  ← roles only exist here, no player can access
  Player A, B, C, D               ← receive render commands only
  → No player ever has a process they own that holds role data.
```

This means:
- Role assignment runs on the server only. Clients never receive their own Anima Orb color.
- Gun resolution (real bullet vs. orb detonation) runs on the server only. The client only sends a "fire requested" input.
- The detonation timer, task completion state, and portal logic are all server-authoritative.

### Why not Steam Matchmaking / Steam Lobbies at all?

Steam Lobbies are used for **player discovery and session setup** before handing off to a listen server. Since we are not using a listen server, Steam Lobbies provide no architectural benefit here.  
GodotSteam is still used exclusively for **Steam authentication** (ticket generation) and future Steam platform features (achievements, overlay, etc.).

Matchmaking is handled entirely by the Go backend: it manages the queue, decides when a group forms a match, spins up a headless Godot container on EC2, and returns the server's IP/port to each client.

---

## Tech Stack

| | Technology |
|---|---|
| Engine | Godot 4 (GDScript) |
| Multiplayer | Godot built-in ENet (UDP) |
| Platform | PC — Steam (GodotSteam) |
| Backend | Go, PostgreSQL, Redis (AWS EC2) |
| Auth | Steam Web API → JWT |

---

## Repository Structure

```
red-hat/
├── apps/
│   ├── client/   ← Godot 4 project
│   └── server/   ← Go backend
└── .github/      ← Copilot instructions, issue/PR templates
```

---

© 2026 everydaystudio. All rights reserved. See [LICENSE](./LICENSE).

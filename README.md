# RED HAT

> **2D pixel-art social deduction extraction game for PC (Steam)**  
> Developed by **everydaystudio**

---

## What is RED HAT?

RED HAT is a reverse-mafia game where players wear LED detonation helmets — but cannot see their own helmet color. Only other players can see it. A shared countdown timer ticks toward detonation. Complete tasks to extend it. Escape through the gate before time runs out.

The twist: the Hunter's gun fires real bullets. The Civilian's gun self-detonates. And the **gunshot sound is identical to the explosion sound**.

### Roles

| Role | Helmet | Win Condition |
|------|--------|---------------|
| Hunter | RED | Kill the WHITE helmet before escape |
| Target | WHITE | Escape through the extraction gate |
| Civilian | BLACK | WHITE escapes → shared win |

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

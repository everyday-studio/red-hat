# RED HAT — Game Design Specification

> **This document defines the game design and product requirements.**
> It is the source of truth for gameplay rules, UX flows, mechanics, and art direction.
> The agent must read this document when implementing any gameplay feature.
> Update this file whenever game design decisions are made or changed.

---

## 1. Concept

**RED HAT** is a **2D top-down survival social deduction** game for PC (Steam).
Genre: "Don't Starve × Among Us" — cooperative survival with a hidden traitor mechanic.

Players are trapped in a dark, surreal world filled with roaming monsters.
They must cooperate to complete rituals and open an escape portal — but one among them secretly carries a killing role.
Everyone wears an LED detonation helmet. You can **never** see your own color. Only others can see yours.

Core loop:
- Survive monsters together while completing tasks.
- The detonation timer counts down. Completing tasks extends it.
- When all tasks are done, the escape portal opens.
- The WHITE helmet must reach the portal — but the RED helmet is trying to kill them first.
- If WHITE is killed, the WHITE role **secretly transfers** to a random surviving BLACK player. The game never ends early.

**Helmet visibility rule**: Every player's helmet color is **always visible to everyone else**. Players cannot hide or control their own helmet's visibility. This is intentional — the core tension is the asymmetry of "others always see what you cannot see about yourself."

---

## 2. Story & Setting

### World & Lore

This place is **purgatory** — a liminal space between life and rebirth.
Every player has arrived here because of sins committed in the living world.
Upon crossing into this dimension, **all memories of the living world are erased**. No one knows who they were.

The rules of this world:
- Those who escape through the portal are **forgiven** and allowed to be reborn.
- Those who fail to escape are **condemned** — to remain here, or worse.
- The **RED hat** is a mark of unforgivable sin. That player **cannot be reborn** regardless of outcome. They are permanently damned. Their only purpose is to ensure no one else escapes either.

Monsters inhabit this purgatory. They are not just obstacles — they are part of the world's fabric, interacting with both the living and the dead.

The aesthetic reference is **Don't Starve**: hand-crafted, gothic, darkly whimsical. Supernatural and fantastical. Not realistic.

### Characters

Players choose a **character** before the match (or are assigned one). Each character has slightly different stats.
Every character arrives in purgatory wearing an LED helmet whose color marks their role — visible to all others, invisible to themselves.
If a helmet detonates — misfire or bullet — it **explodes**.

---

## 3. Players & Roles

### Player Count

- **Minimum**: 4 players. **Maximum**: 8 players.
- Future roles with new helmet colors may be added. Avoid hardcoding role counts.

### Roles

| Role | Helmet | Count | Win Condition |
|------|--------|-------|---------------|
| Hunter | RED | 1 | WHITE never reaches the portal (kill or timer) |
| Target | WHITE | 1 | Escape through the portal |
| Civilian | BLACK | N-2 | WHITE escapes → group win; die → ghost mode |

> Role assignment is **server-only**. The client never receives or stores its own role.

### WHITE Transfer Mechanic

When the current WHITE player is killed:
1. The WHITE role **secretly transfers** to a random surviving BLACK player on the server.
2. The new WHITE player receives **no notification** — they do not know they are now WHITE.
3. The Hunter does **not** know WHITE has transferred — they must deduce from behavior.
4. The original WHITE's corpse remains on the map; their ghost roams freely.
5. **Transfer is unlimited** — WHITE keeps transferring as long as any BLACK is alive.
6. The ghost of the dead WHITE can **observe the full map** as they roam and thus will naturally see who the new WHITE is.

**Game ends when:**
- The current WHITE escapes the portal, or
- There are no surviving BLACK players left to receive the WHITE transfer (Hunter wins), or
- The timer hits 0 (Hunter wins).

---

## 4. Match Flow

```
Login (Steam Auth)
  → Main Menu
    → Enter Lobby (host selects map)
      → All players Ready
        → Countdown
          → Game starts (distributed spawn)
            → [Survive monsters + complete tasks]
              → Win/Lose condition triggered
                → Result Screen (role reveal)
                  → Same party returns to Lobby
```

---

## 5. Lobby

The pre-game waiting room before a match begins.

- **Space**: A small, bounded area — players can walk around freely but cannot leave it.
- **Visibility**: All players' characters and nicknames are visible to each other.
- **Communication**: Voice chat (proximity-based) and text chat are both available.
- **Text chat**: Messages appear as speech bubbles above the character.
- **Ready system**: Each player has a Ready button. The match starts when all players are ready.
- **Ready cancel**: Players can cancel their ready status, but there is a **1-second delay** after confirming ready before cancellation becomes possible.
- **Combat**: Disabled in the lobby — no shooting or detonation.
- **Monsters**: No monsters in the lobby.

> 🔲 TBD: Whether host can kick players, lobby capacity display.

---

## 6. Maps

- Total number of maps: **not yet finalized**.
- Maps are dark, open-ish environments inspired by Don't Starve — wilderness, ruins, cursed forests, etc.
- Maps contain **roaming monsters** that players must avoid while completing tasks.
- **Map selection**: The host (방장) selects the map directly in the lobby. No random/vote — host decision only.

> 🔲 TBD: Map themes, size, number of tasks per map, number of CCTV cameras per map, number of escape portals per map.

---

## 7. Monsters

Monsters roam the map and are part of purgatory's fabric.
They interact with both **living players** and **ghosts** — both can affect them.

### Day/Night Cycle

The match world cycles between **day** and **night** phases.

| Phase | Duration | Monsters | Player Vision |
|-------|----------|----------|---------------|
| Day | **2 min** | Most monsters **sleeping** — dormant and non-threatening unless disturbed | Normal range |
| Night | **1 min** | All monsters **awake and roaming** — primary danger window | **Narrowed** — visible range decreases significantly |

- The cycle repeats continuously throughout the match (Day → Night → Day → Night → …).
- Daytime is the primary window for task completion and coordination.
- Nighttime forces players to prioritize survival over tasks.
- **Ghosts can wake sleeping monsters during the day** (see Ghost System), extending danger outside the normal night window.

### Monster Properties

- **Damage to living**: Monster contact reduces player HP by half their max HP.
- **No instant death**: Monsters do not cause immediate death (unless HP is already at half).
- **Shootable**: Players can shoot and kill monsters with their revolver.
- **Bullets**: Bullets do not penetrate — stop on first hit (monster or player).
- **Ghost interaction**: Dead players (ghosts) can wake sleeping monsters, or disturb and redirect active ones (see Ghost System).
- **RED possession**: The RED ghost can possess and directly control a monster (see Ghost System).

> 🔲 TBD: Monster behavior AI (patrol, chase, aggro range), monster HP, respawn rules, number of monsters per map, visual/audio cue for day/night transition.

---

## 8. Tasks

- **Dev reference**: Among Us, Goose Goose Duck.
- Completing a task **extends the shared detonation timer**.
- Task completion state is **server-authoritative** and synced to all clients.
- Tasks are thematically framed as rituals, repairs, or arcane interactions with the world.
- **Each match draws from a preset task pool randomly** — exact tasks vary per match.
- Task interaction styles are **mixed**: some are button holds, some are object delivery, some are minigames.

### Task Breakdown

| Difficulty | Count | Coop Mechanic | Timer Extension |
|------------|-------|---------------|-----------------|
| Easy | 5 | Solo only | +0.5 min |
| Medium | 3 | **2 players separated** | +1.5 min |
| Hard | 2 | **2 players co-located** | +3.0 min |

**Total tasks per match: 10+**
Maximum combined extension: (5 × 0.5) + (3 × 1.5) + (2 × 3.0) = **+13 minutes**

#### Medium Task — Separated Mechanic
Two players must each operate a task node at **physically distant locations simultaneously**.
Neither node progresses until both are being operated at the same time.
The two node locations are **well beyond screen range** from each other.

> Design intent: Forces players far enough apart that helmet colors become invisible to each other. RED can exploit this isolation window to eliminate the separated target.

#### Hard Task — Co-located Mechanic
Two players must stand **at the same task node simultaneously** and both contribute to progress.

> 🔲 TBD: Specific task names and designs, minimum separation distance for medium tasks (pixels), whether tasks can be interrupted or reset by monsters, whether task progress persists if a player steps away mid-task.

---

## 9. Detonation Timer

- Base duration: **600 seconds** (10 minutes).
- **Proximity drain**: If more than 2 players are within proximity radius (96px), timer drains **2× faster**.
- Completing a task extends the timer.
- When the timer hits 0: **Hunter wins** (see Section 11 — Win Conditions)

---

## 10. Key Mechanics

### Gun System

- Only the **RED** helmet's gun fires real bullets.
- **WHITE / BLACK** guns instantly **self-detonate the shooter's own helmet** — a clean explosion.
- **RED cannot self-detonate** — their gun fires real projectiles outward. The self-detonate mechanic does not apply to RED.
- The client only requests a fire action. The server resolves the outcome.
- Neither gunshot nor explosion has a distinguishing sound — both events share one audio file.

### Gun Behavior

- **Default weapon**: Revolver — no rapid fire, has a cooldown between shots.
- **Aiming**: Mouse cursor direction.
- Projectile is near-instantaneous (effectively hitscan in feel, implemented as very fast projectile).
- Nearly impossible to dodge by movement alone.
- **No penetration**: Bullet stops on first hit (player or monster). Does not pass through.
- Projectile disappears on wall contact, or if no target is in its path.
- Bullets can hit monsters as well as players.

### Health System

- **Base HP**: 2 (all characters unless otherwise specified).
- **Bullet hit**: Instant death regardless of current HP.
- **Monster hit**: Reduces HP by half of the character's max HP (e.g. base 2 HP → -1 per hit).
- **Natural regeneration**: HP recovers over time. Rate varies by character.
- No healing items (unless specified later).

### Character System

Players choose and customize a character before the match.
All characters share the same role system (RED / WHITE / BLACK).

**Customization:**
- Players mix and match presets across categories: **hair, body type, face shape**, etc.
- Multiple players can select the same base character — customization sets them apart visually.
- **Target roster at launch: 4 characters**, each with distinct stat profiles.

| Stat | Base | Notes |
|------|------|-------|
| HP | 2 | Some characters may have 3 HP with reduced speed |
| Move Speed | Normal | Some characters faster or slower |
| Regen Rate | Normal | Varies by character |

> 🔲 TBD: Final character roster and names, exact stat values per character, unique passives/actives per character (currently undecided — may stay stat-only), whether character selection locks after lobby ready.

### Sound Deception

- Gunshot SFX and helmet explosion SFX are the **same audio file** (`shot_or_explosion.wav`).
- Players cannot distinguish kills from self-detonations by sound alone.
- This is an **intentional core mechanic** — never use separate files.

### Proximity Interference

- If more than 2 players gather within `96px`, the detonation timer drains **2× faster**.
- Proximity check authority: **TBD** (server-authoritative or client-calculated).

### Proximity Voice Chat

- Players hear each other's voice based on in-game distance.
- Volume attenuates with distance.
- **Audible range**: Slightly wider than the visible screen viewport — players just off-screen can still hear and be heard.

### Text Chat

- Messages appear as speech bubbles above the sender's character.
- **Visibility**: Speech bubbles are only visible when the character is within the screen viewport — the bubble is attached to the character node and naturally hidden when off-screen.

> 🔲 TBD: Can ghosts communicate with living players via text?

### Items

Confirmed items. Do not implement until item system is scoped, but these are the final designs:

**Acquisition:**
- Most items **spawn on the map** at game start — first player to reach them picks them up.
- Some items are awarded as **task completion rewards** (TBD: which tasks reward which items).
- Maximum carry limit: **TBD**.
- Whether held items are visible to other players: **TBD**.

| Item | Type | Effect |
|------|------|--------|
| 망원경 (Telescope) | Passive (held) | Extends view range. View center shifts to 1/3 point between character and mouse cursor. Viewport size also increases. |
| 도청장치 (Wiretap) | Passive (held) | Greatly extends proximity voice chat hearing range. Displays the name of who is currently speaking. |
| 레이더 (Radar) | Passive (held) | Detects nearby players. Shown on minimap and as edge indicators on screen. |
| 사이클롭스의 눈 (Cyclops Eye) | Passive (held) | Removes all CCTV blind spots — full camera visibility. |
| 헤르메스의 신발 (Hermes' Shoes) | Passive (held) | Increases movement speed. Eliminates footstep sounds. |
| 공포의 호루라기 (Horror Whistle) | Active (one-use) | Emits a gunshot sound at the user's location. Single use. |

### Ghost System

When a player dies, their corpse stays and they become a **ghost** that roams the map freely.
Ghosts are invisible to living players and cannot communicate with them.

**All ghosts can:**
- Freely roam the entire map (no collision, no restrictions).
- Observe the full game state, including who the current WHITE is.
- **Interact with tasks** (assist in progressing tasks).
- **Wake sleeping monsters** — disturb dormant monsters during the day, triggering them to become active.
- **Distract active monsters** — slow or redirect them away from (or toward) living players.
- **Proximity chat with other ghosts only** (living players cannot hear ghost chat).

**RED ghost only:**
- Can **possess a monster**, taking **direct movement control** of it.
- While possessing, RED drives the monster to hunt or disrupt living players.
- Possession is **exclusive**: one monster at a time.
- Possession has a **time limit** (exact duration: TBD).
- After possession ends, there is a **cooldown** before RED can possess again (exact duration: TBD).

> **Design note**: Possession is retained because in ghost form RED is no longer threatened by monsters — the thematic tension dissolves. Lore-wise, the condemned soul merging with the darkness of purgatory is intentional. If playtesting shows it is overpowered, the fallback is to replace it with "wake multiple monsters simultaneously" (stronger than the single-wake ability of regular ghosts).

> 🔲 TBD: Possession time limit and cooldown values, whether ghosts can switch between free-roam and CCTV spectator view.

### CCTV Spectator

- Ghost players may additionally spectate via fixed cameras placed around the map.
- Cameras have **blind spots** — not all areas are visible.

> 🔲 TBD: Number of cameras per map, CCTV switching UI.

---

## 11. Win & Lose Conditions

### Confirmed Rules

- **WHITE escapes the portal** → WHITE + all Civilians win (group victory).
- **Timer hits 0** → Hunter wins.
- **No BLACK left to receive WHITE transfer** → Hunter wins.
- **Hunter (RED) dies** → Game continues. WHITE must still reach the portal to win.
- **Entering the portal as non-WHITE (while WHITE is alive)** → that player dies instantly (counts as death, not escape).
- After WHITE escapes, remaining Civilians may individually escape through the portal — this is personal survival (stat recorded), not a separate win condition.

### Escape Flow Summary

```
[WHITE alive]
  → Non-WHITE entering portal = instant death
  → WHITE escapes → WHITE + all Civilians win
  → WHITE dies → WHITE transfers to random BLACK → repeat

[No BLACK remaining]
  → Hunter wins (WHITE role has nowhere to go)

[Timer = 0]
  → Hunter wins

---

## 12. Escape Portal

- The map contains **multiple portal locations**, but only **one portal unlocks** when all tasks are completed.
- The active portal's location is visible to all players once unlocked.
- **Channel time**: Players must stand in the portal for **10 seconds** to escape (like a Minecraft Nether portal).
- **Interrupt**: Being **attacked does not interrupt** channeling — the player keeps channeling through damage.
- **Entry rules**:
  - WHITE can enter at any time after unlocking.
  - Non-WHITE entering while WHITE is alive = **instant death**.
  - After WHITE escapes, remaining players may enter individually.
- WHITE must physically channel the portal to win — no auto-win even if Hunter is dead.

> 🔲 TBD: Visual/audio feedback during channeling, whether channeling progress persists if the player briefly steps out of the portal.

---

## 13. Spawning

- Players spawn at **distributed locations across the map** at game start.
- No two players spawn at the same point.
- Spawn positions are predefined per map (Marker2D nodes).
- Distributed spawn is intentional: if all players spawn together, helmet colors are immediately visible to everyone from the start.

---

## 14. Death

- "KO" terminology is not used — it is always **death** (permanent within the match).
- **Causes of death**:
  - Helmet detonation (self-misfire for WHITE/BLACK, or hit by RED's bullet) → explosion
  - HP reduced to 0 by monster damage → collapse
- **On death**:
  - The corpse **remains on the map** at the location of death.
  - The player enters **ghost mode** (see Ghost System in Section 10).
- Footstep sounds exist normally. **Hermes' Shoes** item suppresses them.

> 🔲 TBD: Death animation, whether monster-kill also produces an explosion or a different visual.

---

## 15. End-Game Screen

Displayed when a match ends:

- **Role reveal**: All players' helmet colors (roles) are shown.
- **Winning team**: The team that achieved their win condition is announced.
- Additional stats (kill log, task completion rate, etc.): **TBD**.

> After the result screen the same party **returns to the Lobby** (see Section 4 — Match Flow).

---

## 16. Art Direction

- **Style**: 2D, top-down, hand-crafted gothic pixel art.
- **Primary reference**: Don't Starve / Don't Starve Together — dark, whimsical, slightly grotesque.
- **Secondary references**: Among Us (social deduction UX), Undertale (top-down feel).
- **Character style**: Stylized, slightly exaggerated proportions. Dark and eerie, not cute.
- **Environment**: Cursed wilderness, dark forests, crumbling ruins, arcane structures. Surreal and fantastical.
- The world doesn't need to be explained or realistic. The darkness and monsters are simply part of it.
- **Target resolution**: QHD (2560×1440) as baseline. Players can change resolution in settings.
- **Nickname**: Steam nickname imported by default. Players may change their in-game display name within the game.

> 🔲 TBD: Pixel font usage, UI color palette.

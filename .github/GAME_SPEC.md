# SOMNIUM — Game Design Specification

> **This document defines the game design and product requirements.**
> It is the source of truth for gameplay rules, UX flows, mechanics, and art direction.
> The agent must read this document when implementing any gameplay feature.
> Update this file whenever game design decisions are made or changed.

---

## 1. Concept

**SOMNIUM** is a **2D top-down survival social deduction** game for PC (Steam).
Genre: "Don't Starve × Among Us" — cooperative survival with a hidden traitor mechanic.

Players wake with no memory in a dark, unknown place, a strange pulsing orb hovering above their heads and a foreboding beeping in their ears.
They must cooperate to complete rituals and open an escape portal — but one among them secretly carries a killing role.
Every player has an **Anima Orb** floating above them — a soul fragment made visible in the dream. You can **never** see your own orb's color. Only others can see yours.

Core loop:
- Survive monsters together while completing tasks.
- The detonation timer counts down. Completing tasks extends it.
- When all tasks are done, the escape portal opens.
- The SILVER orb must reach the portal — but the CRIMSON orb is trying to kill them first.
- If SILVER is killed, the SILVER role **secretly transfers** to a random surviving ASHEN player. The game never ends early.

**Anima Orb visibility rule**: Every player's orb color is **always visible to everyone else**. Players cannot hide or control their own orb's visibility. This is intentional — the core tension is the asymmetry of "others always see what you cannot see about yourself."

---

## 2. Story & Setting

### Hidden Narrative (Designer Reference — Never Stated In-Game)

> The true story is **never communicated directly to players.**
> It is revealed only through environmental props, task flavor text, and map easter eggs.
> Players who pay attention may piece it together. Most will not — and that is intentional.

**The actual story:**
A multi-vehicle pile-up on a highway. Multiple strangers, critically injured at the same moment, all fell into comas simultaneously. Their unconscious minds converged in a shared dream space.

The **RED player** is the drunk driver who caused the crash. Inside the dream, no one retains their memories — not even RED. The others do not know who is responsible. RED does not know what they did.

- The **Anima Orb** represents the thread of life — its pulse is the medical monitor. Its color reveals what others cannot tell from looking at themselves.
- **Orb detonation** (the orb bursting apart) = death on the operating table.
- **Escaping through the portal** = surviving — waking up.

**Environmental storytelling (easter eggs — scattered through the map, never explained):**
- Wrecked and mangled vehicles, half-buried under overgrowth.
- One wreck has a **child's car seat** visible through a shattered window.
- Skid marks burned into the ground. Shattered safety glass scattered nearby.
- A broken dashboard or hospital monitor prop in some corner of the map.
- Task flavor text written from a dissociated, dreamlike perspective — eerily familiar.
- The **Resuscitation** task (H2) mirrors CPR in its rhythm mechanics.
- The **Signal** task (E4) references transmitting a distress call.

---

### World & Lore (Player Perspective — What Characters Know)

**Player scenario:**
You wake up.
Your memories are gone — completely blank. No name, no face, no past.
Above your head floats a small orb — faintly glowing, pulsing with a steady ominous *beep… beep…* that you feel more than hear.
You look around. A dark, overgrown place — somewhere between wilderness and ruin, suffocatingly real but wrong in ways you cannot name.
Others are here too. Each one has an orb above their head as well. All just as lost.
Somewhere in this place, there may be a way out. But you don't know why you're here, or what is hunting you.

**What is this place?**
A shared dream. Though no one inside it knows that.
It feels entirely real — the cold air, the soft glow of the orb, the fear.
The creatures here are shaped like wolves, dogs, and bears — but larger, darker, wrong. Whether they are animals or something the dream invented, nobody knows.
The only certainty is the beeping, and the need to find a way out before it stops.

The aesthetic reference is **Don't Starve**: hand-crafted, gothic, darkly whimsical. Supernatural and fantastical. Not realistic.

### Characters

Players choose a **character** before the match (or are assigned one). Each character has slightly different stats.
Every character has an **Anima Orb** floating above their head — a soul fragment made visible in the dream world. Its color marks their role, visible to all others, invisible to themselves.
If an orb detonates — misfire or bullet — it **bursts apart in a flash of light**.

---

## 3. Players & Roles

### Player Count

- **Minimum**: 4 players. **Maximum**: 8 players.
- Future roles with new orb colors may be added. Avoid hardcoding role counts.

### Roles

| Role | Anima Orb | Count | Win Condition |
|------|-----------|-------|---------------|
| Hunter | CRIMSON (심홍빛) | 1 | SILVER never reaches the portal — kill or let timer expire |
| Target | SILVER (은백빛) | 1 | Escape through the portal |
| Civilian | ASHEN (재빛) | N-2 | SILVER escapes → group win; die → ghost mode |

> Role assignment is **server-only**. The client never receives or stores its own role.

### Anima Orb Visual Design

- The orb floats ~8px above the character's head, bobbing gently.
- It pulses slowly (brightness cycles) in sync with the detonation timer beep.
- **CRIMSON**: deep red-orange glow. Slightly faster pulse than others — subtly unsettling.
- **SILVER**: bright silver-white. The strongest glow — most visible orb on the map.
- **ASHEN**: dull grey with faint inner light. Unremarkable by design.
- **Own orb**: The local player's orb is **not rendered** on their screen. Completely invisible to themselves.
- All orbs are always visible at full opacity regardless of distance, fog, or night phase — the soul cannot be hidden.

### SILVER Transfer Mechanic

When the current SILVER player is killed:
1. The SILVER role **secretly transfers** to a random surviving ASHEN player on the server.
2. The new SILVER player receives **no notification** — they do not know they are now SILVER. Their orb color changes on all *other* players' screens, but they cannot see their own orb.
3. The Hunter does **not** know SILVER has transferred — they must deduce from behavior.
4. The original SILVER's corpse remains on the map; their ghost roams freely.
5. **Transfer is unlimited** — SILVER keeps transferring as long as any ASHEN is alive.
6. The ghost of the dead SILVER can **observe the full map** as they roam and thus will naturally see who the new SILVER is.

**Game ends when:**
- The current SILVER escapes the portal, or
- There are no surviving ASHEN players left to receive the SILVER transfer (Hunter wins).

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

- Maps are dark, open-ish environments inspired by Don't Starve — wilderness, ruins, cursed forests, etc.
- Maps contain **roaming monsters** that players must avoid while completing tasks.
- **Map selection**: The host (방장) selects the map directly in the lobby. No random/vote — host decision only.
- **Initial release**: One map theme. Multi-theme support is a future expansion.

> 🔲 TBD: Total number of distinct maps at launch, number of CCTV cameras per map.

---

### 6.1 Map Generation Overview

Every match generates a **new map procedurally** from a random seed.
The generation is fully server-side and deterministic given the seed.

**Generation approach**: Fully procedural — no hand-crafted TileMap files.
**Structural approach**: Zone-based — the map is composed of thematic sub-zones connected by organic passages.
**Theme**: Layout generation is theme-agnostic. The chosen tileset is applied externally after layout is complete.

**Visual reference**: The minimap in the game depicts an organic, irregular landmass with several open clearings connected by winding corridors — similar to a cave system or cursed forest clearing. Not rectangular. Not symmetric. Alive-feeling.

---

### 6.2 Map Size

| Dimension | Value |
|-----------|-------|
| Map width | ~200 tiles |
| Map height | ~150 tiles |
| Tile size | 16×16 px |
| World size | ~3200×2400 px |
| Visible viewport (QHD) | ~160×90 tiles |

The map is roughly **2–3 screen-widths** across in each dimension.
A player at the center of the map cannot see the edge — exploration is required.

---

### 6.3 Generation Algorithm

**Algorithm: Cellular Automata + Zone Graph (CA-ZG)**

This produces organic, irregular shapes that match the minimap reference image:
- Soft, non-rectangular outer boundaries (like a cave or landmass)
- Distinct open clearings (zones) naturally formed within the mass
- Winding corridors connecting zones without hard lines

#### Step 1 — Task-First Placement

Task nodes are placed **before** any terrain is generated.
All terrain is then shaped around guaranteeing navigable paths to those nodes.

1. Choose the number of tasks and their difficulties.
2. Place **Medium Task node pairs** first, each pair separated by at minimum **400 tiles** (guaranteed well beyond screen range).
3. Place **Hard Task nodes** at mutually reachable positions.
4. Place **Easy Task nodes** distributed across remaining space.
5. Mark all task positions as **required floor anchors** — the CA step must produce floor tiles at and around every anchor.

#### Step 2 — Zone Graph Construction

1. Use task anchor positions as seed points for zone centers.
2. Add 2–4 additional "common area" zone centers (open plazas, meeting points) using Poisson disk sampling to ensure spacing.
3. Build a **minimum spanning tree (MST)** over all zone centers — this defines the required corridor connections.
4. Zone centers and MST edges become the skeleton the CA will expand from.

#### Step 3 — Cellular Automata Expansion

1. Initialize a 2D boolean grid (all walls).
2. Mark zone centers and MST corridor paths as initial floor seeds.
3. Run CA for **5–8 iterations** using a standard birth/survival rule:
   - A wall cell becomes floor if it has ≥ 5 floor neighbors (3×3 Moore neighborhood)
   - A floor cell stays floor if it has ≥ 4 floor neighbors
4. This organically grows the seeds into irregular, cave-like open zones with winding corridors.
5. After CA, run a **flood-fill check** from every required floor anchor. Any disconnected anchor triggers a targeted corridor carve to reconnect it.

#### Step 4 — Boundary & Cleanup

1. Perform a flood-fill from the map edge. Any floor tiles reachable from the edge are too close to the boundary — convert them to walls.
2. Fill small isolated floor pockets (< 50 tiles) to prevent unreachable dead-end zones.
3. Apply a **wall-thinning pass** on corridors that are only 1 tile wide — widen to minimum 2 tiles for player traversal.

#### Step 5 — Marker Placement

After terrain is finalized, place all game markers onto valid floor tiles:

| Marker | Placement Rule |
|--------|---------------|
| Player spawns (8 max) | Distributed across zones, minimum 200px apart from each other |
| Task nodes | Already anchored in Step 1; validated against final floor grid |
| Portal candidate locations (3–5) | Placed in open zones far from spawn cluster |
| Monster spawns | Distributed across all zones, prefer edges and narrow passages |
| Item spawns | Scattered across floor tiles, minimum 2 per zone |
| CCTV camera positions | Placed at zone junctions with maximum coverage angle |

#### Step 6 — Validation

Before the map is accepted, run:
1. **Connectivity check**: every floor anchor is reachable from every player spawn.
2. **Medium Task distance check**: each Medium Task node pair is ≥ 400 tiles apart.
3. **Portal exclusion**: no portal candidate overlaps a spawn cluster radius.
4. If any check fails, **regenerate with a new seed** (do not patch — restart cleanly).

---

### 6.4 Tileset & Theme

- Layout is **theme-agnostic** — the CA produces only floor/wall boolean data.
- A tileset is applied as a post-processing step: floor tiles → ground variant, wall tiles → wall variant, borders → decorative edges.
- **Initial release tileset**: Cursed Purgatory (dark stone, overgrown ruins, arcane glow).
- Future themed tilesets can be swapped in without changing any generation code.

> 🔲 TBD: Exact tile variant rules (Wang tile autotiling or simple rule-based), decorative props scatter density, whether water/hazard tiles are used as impassable non-wall terrain.

---

## 7. Monsters

Three types of monsters inhabit the dream world. All are shaped like familiar animals but behave according to their own logic — wrong in ways that feel dreamlike.

**Bullets do not kill monsters.** Shooting a monster only **triggers** it: wakes a sleeping one, or aggros/redirects an active one.
All monsters interact with both **living players** and **ghosts** (see Ghost System).

### Day/Night Cycle

| Phase | Duration | Description | Player Vision |
|-------|----------|-------------|---------------|
| Day | **2 min** | Varies by monster type | Normal range |
| Night | **1 min** | All types become aggressive | **Narrowed** significantly |

- The cycle repeats continuously: Day → Night → Day → Night → …
- Daytime is the primary window for task completion and coordination.
- Nighttime forces players to prioritize survival.
- **Ghosts can wake sleeping monsters during the day** (see Ghost System).

**Day/Night Transition UI:**
- A **clock-like dial** in a corner of the screen shows position within the current cycle.
- On transition to Night: a sound effect plays (exact SFX: TBD) and the screen gradually darkens.

---

### Monster Types

#### 들개 (Stray — name TBD)
A gaunt, twitchy creature shaped like a stray dog.

| Phase | Behavior |
|-------|----------|
| Day | **Avoids** players — retreats if a player comes near. |
| Night | Eyes turn **red**. Becomes hostile. Chases any player that enters its aggro radius. **Faster than the player.** Stops chasing once the player exits its territory radius. |

- Has a defined territory radius. Pursuit ends at that boundary.
- Bullet hit: immediately triggers aggro regardless of phase.

---

#### 곰 (Bear — name TBD)
A massive, hulking creature. Aggressive at all times. Extremely wide detection range.

| Phase | Behavior |
|-------|----------|
| Day | **Slow pursuit** — always hostile, but moves **slower than the player**. Wide aggro range. |
| Night | **Fast pursuit** — moves **faster than the player**. |

- Does not respect a territory radius — pursues indefinitely once aggroed.
- Aggro range is **wider** than either of the other types.
- Bullet hit: briefly increases chase speed (triggers rage).

---

#### 늑대 (Wolf — name TBD)
Wolf-shaped creatures that move and attack in **packs of 2–4**.

| Phase | Behavior |
|-------|----------|
| Day | **Sleeping** — dormant and non-threatening unless disturbed. |
| Night | Awake and roaming. When any wolf spots a player, the **entire pack locks onto that one target simultaneously**. Faster than the player. |

- The pack never splits its attention — all wolves attack the same target.
- Bullet hit: wakes the entire pack and triggers immediate group aggro.

---

### Monster Interactions

| | Stray | Bear | Wolf |
|--|-------|------|------|
| **Stray** | — | Avoids | Ignores |
| **Bear** | Chases | — | Chases |
| **Wolf** | Ignores | Avoids | — |

- Bears actively chase both Strays and Wolves.
- Strays and Wolves both avoid Bears.
- Strays and Wolves do not interact with each other.

---

### Common Monster Rules

- **Damage**: Monster contact reduces the player's HP by half their max HP.
- **No instant death**: One hit does not kill unless HP was already at half.
- **Bullets trigger, not kill**: A bullet hit aggros or redirects a monster — it does not deal damage.
- **Ghost interaction**: Ghosts can wake or redirect monsters (see Ghost System).
- **CRIMSON possession**: The CRIMSON ghost can possess any monster type (see Ghost System).

> 🔲 TBD: Exact aggro and territory radius values in pixels, number of each monster type per map, respawn behavior if a monster roams too far from its spawn origin.

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

> Design intent: Forces players far enough apart that orb colors become harder to cross-reference at a glance. CRIMSON can exploit this isolation window to eliminate the separated SILVER target.

#### Hard Task — Co-located Mechanic
Two players must stand **at the same task node simultaneously** and both contribute to progress.

### Task Interruption Rules

- **Monster contact** while doing a task → task is **immediately interrupted** (progress reset per-task rule below).
- **Leaving the task area voluntarily** → progress behavior varies by task (specified per task).
- **Portal channeling** follows the same rule: monster contact **interrupts** channeling and resets progress to 0.
- **Channeling UI**: All channeling (task or portal) displays a **linear progress bar above the player's head**.

---

### Task Catalog

All 10 tasks are used every match. They are placed procedurally across the map (see Section 6).

---

#### Easy Tasks — Solo, +0.5 min each

**E1 — 봉화 점화 (Ignite the Beacon)**
Hold interact on an unlit brazier for 3 seconds to light it.
- Interaction: button hold
- Monster interruption: resets to 0
- Abandonment: resets to 0
- Flavor text on completion: *"The light goes on. Someone will see it."*
- Easter egg: the brazier is shaped like a road flare.

**E2 — 파편 수거 (Collect the Shards)**
Pick up 3 glass shards scattered within a radius and deliver them to a container node.
- Interaction: object pickup + delivery
- Monster interruption: drops held shard; already-delivered shards remain
- Abandonment: already-delivered shards remain (partial progress maintained)
- Flavor text on completion: *"Tiny pieces of something that was once whole."*
- Easter egg: the shards are windshield glass.

**E3 — 혈흔 정화 (Cleanse the Stain)**
Hold interact over a dark stain on the ground for 5 seconds.
- Interaction: button hold
- Monster interruption: resets to 0
- Abandonment: resets to 0
- Flavor text on completion: *"It comes clean. As if it was never there."*
- Easter egg: the stain is shaped like tire tracks.

**E4 — 신호 송신 (Transmit the Signal)**
Interact with a broken device; align a dial to a target frequency using directional inputs (simple minigame).
- Interaction: dial alignment minigame
- Monster interruption: resets dial to 0
- Abandonment: resets dial to 0
- Flavor text on completion: *"Static. Then, for a moment — a voice."*

**E5 — 인형 반환 (Return the Doll)**
Find a small cloth doll somewhere on the map and carry it to a marked cradle location.
- Interaction: object pickup + delivery
- Monster interruption: doll is dropped at current position; remains there until re-picked up
- Abandonment: doll stays at last position (partial progress maintained)
- Flavor text on completion: *"It belongs somewhere safe."*
- Easter egg: a child's toy.

---

#### Medium Tasks — 2 Players Separated, +1.5 min each

**M1 — 제단 동기화 (Simultaneous Altar Activation)**
Two players each stand at a distant ritual altar and hold interact at the same time for 8 seconds. Neither altar progresses unless both are active simultaneously.
- Interaction: synchronized button hold
- Monster interruption: full reset (both nodes)
- Abandonment: full reset
- Flavor text on completion: *"It responds. As if waiting for two voices at once."*

**M2 — 주파수 동기화 (Frequency Sync)**
Two players at distant signal towers each tune a dial to a matching frequency, then both confirm simultaneously.
- Interaction: individual dial minigame + synchronized confirm
- Monster interruption: resets only the affected player's dial; the other's tuning is maintained
- Abandonment: individual tuning maintained; synchronized confirm step resets if either player leaves
- Flavor text on completion: *"Two signals. One place."*

**M3 — 균형추 조정 (Calibrate the Scales)**
Two players at opposite ends of a mechanism each deliver one weight item to their own node. Both items must be placed for the task to complete.
- Interaction: object delivery to respective nodes
- Monster interruption: drops held item; already-placed items remain
- Abandonment: placed items remain (each node independently maintains progress)
- Flavor text on completion: *"Equal weight. Equal measure."*

---

#### Hard Tasks — 2 Players Co-located, +3 min each

**H1 — 대의식 채널링 (Grand Ritual Channeling)**
Two players stand inside a large glowing ritual circle together and both hold interact for 15 seconds continuously.
- Interaction: synchronized sustained hold
- Monster interruption: full reset
- Abandonment: full reset
- Flavor text on completion: *"Hold fast. Don't let go."*

**H2 — 심폐 소생 (Resuscitation)**
Two players alternate pressing interact at a shared ritual object in a rhythm pattern (one presses, then the other, timed to a visual pulse).
- Interaction: alternating timed button press (rhythm minigame)
- Monster interruption: full reset
- Abandonment: full reset (rhythm breaks immediately on departure)
- Flavor text on completion: *"One. Two. One. Two. Come back."*
- Easter egg: the motion mirrors CPR.

---

## 9. Detonation Timer

- Base duration: **480 seconds** (8 minutes) — this is the default starting time given to each player.
- **Per-player timer**: Each player has their own individual detonation timer. Timer values can differ between players (e.g. via character stats or future mechanics).
- Completing a task extends the timer for **all players simultaneously**.
- When a task is completed, a brief on-screen notification appears for all players: **"[PlayerName]이(가) X분을 연장했습니다"** — displayed for approximately 3 seconds.
- The timer reaching 0 does **not** trigger a win condition — the match continues until SILVER escapes or no ASHEN players remain.

> 🔲 TBD: How per-player timer differences are surfaced in the HUD.

---

## 10. Key Mechanics

### Gun System

- Only the **CRIMSON** role's gun fires real bullets.
- **SILVER / ASHEN** guns instantly **detonate the shooter's own Anima Orb** — a burst of light and the orb shatters.
- **CRIMSON cannot self-detonate** — their gun fires real projectiles outward. The self-detonate mechanic does not apply to CRIMSON.
- The client only requests a fire action. The server resolves the outcome.
- Neither gunshot nor orb detonation has a distinguishing sound — both events share one audio file.

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
All characters share the same role system (CRIMSON / SILVER / ASHEN).

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

- Gunshot SFX and orb detonation SFX are the **same audio file** (`shot_or_explosion.wav`).
- Players cannot distinguish kills from self-detonations by sound alone.
- This is an **intentional core mechanic** — never use separate files.

### Proximity Voice Chat

- Players hear each other's voice based on in-game distance.
- Volume attenuates with distance.
- **Audible range**: Slightly wider than the visible screen viewport — players just off-screen can still hear and be heard.

### Text Chat

- Messages appear as speech bubbles above the sender's character.
- **Visibility**: Speech bubbles are only visible when the character is within the screen viewport — the bubble is attached to the character node and naturally hidden when off-screen.

- Ghosts **cannot** send text messages to living players. Ghost text chat is visible to other ghosts only.

### Items

> 🔲 TBD: Item designs are not yet finalized. Do not implement until the item system is fully scoped.

**Acquisition (confirmed):**
- Most items **spawn on the map** at game start — first player to reach them picks them up.
- Some items are awarded as **task completion rewards** (TBD: which tasks reward which items).
- Maximum carry limit: **TBD**.
- Whether held items are visible to other players: **TBD**.

### Ghost System

When a player dies, their corpse stays and they become a **ghost** that roams the map freely.
Ghosts are invisible to living players and cannot communicate with them.

**All ghosts can:**
- Freely roam the entire map (no collision, no restrictions).
- Observe the full game state, including who the current SILVER is.
- **Interact with tasks** (assist in progressing tasks).
- **Wake sleeping monsters** — disturb dormant monsters during the day, triggering them to become active.
- **Distract active monsters** — slow or redirect them away from (or toward) living players.
- **Proximity chat with other ghosts only** (living players cannot hear ghost chat).

**CRIMSON ghost only:**
- Can **possess a monster**, taking **direct movement control** of it.
- While possessing, CRIMSON drives the monster to hunt or disrupt living players.
- Possession is **exclusive**: one monster at a time.
- Possession duration: **30 seconds**.
- Cooldown after possession ends: **60 seconds**.

> **Design note**: Possession is retained because in ghost form CRIMSON is no longer threatened by monsters — the thematic tension dissolves. If playtesting shows it is overpowered, the fallback is to replace it with "wake multiple monsters simultaneously" (stronger than the single-wake ability of regular ghosts).

> 🔲 TBD: Whether ghosts can switch between free-roam and CCTV spectator view.

### CCTV Spectator

- Ghost players may additionally spectate via fixed cameras placed around the map.
- Cameras have **blind spots** — not all areas are visible.

> 🔲 TBD: Number of cameras per map, CCTV switching UI.

---

## 11. Win & Lose Conditions

### Confirmed Rules

- **SILVER escapes the portal** → SILVER + all Civilians win (group victory).
- **No ASHEN left to receive SILVER transfer** → Hunter wins.
- **Hunter (CRIMSON) dies** → Game continues. SILVER must still reach the portal to win.
- **Entering the portal as non-SILVER (while SILVER is alive)** → that player dies instantly (counts as death, not escape).
- After SILVER escapes, remaining Civilians may individually escape through the portal — this is personal survival (stat recorded), not a separate win condition.

### Escape Flow Summary

```
[SILVER alive]
  → Non-SILVER entering portal = instant death
  → SILVER escapes → SILVER + all Civilians win
  → SILVER dies → SILVER transfers to random ASHEN → repeat

[No ASHEN remaining]
  → Hunter wins (SILVER role has nowhere to go)

---

## 12. Escape Portal

- The map contains **multiple portal locations**, but only **one portal unlocks** when all tasks are completed.
- The active portal's location is visible to all players once unlocked.
- **Channel time**: Players must stand in the portal for **10 seconds** to escape (like a Minecraft Nether portal).
- **Interrupt**: Monster contact **interrupts channeling** immediately — progress resets to 0.
- **Channeling UI**: A linear progress bar is displayed above the player's head during channeling.
- **Entry rules**:
  - SILVER can enter at any time after unlocking.
  - Non-SILVER entering while SILVER is alive = **instant death**.
  - After SILVER escapes, remaining players may enter individually.
- SILVER must physically channel the portal to win — no auto-win even if Hunter is dead.

> 🔲 TBD: Whether portal channeling progress persists if a player briefly steps out without monster contact.

---

## 13. Spawning

- Players spawn at **distributed locations across the map** at game start.
- No two players spawn at the same point.
- Spawn positions are predefined per map (Marker2D nodes).
- Distributed spawn is intentional: if all players spawn together, all orb colors are immediately visible to everyone from the start — the CRIMSON orb would be identified instantly.

---

## 14. Death

- "KO" terminology is not used — it is always **death** (permanent within the match).
- **Causes of death**:
  - Anima Orb detonation (self-misfire for SILVER/ASHEN, or hit by CRIMSON's bullet) → orb bursts in a flash, body collapses
  - HP reduced to 0 by monster damage → collapse
- **On death**:
  - The corpse **remains on the map** at the location of death.
  - The player enters **ghost mode** (see Ghost System in Section 10).
- Footstep sounds exist normally.

> 🔲 TBD: Death animation for orb burst vs monster collapse (likely different visuals), exact orb burst particle effect.

---

## 15. End-Game Screen

Displayed when a match ends:

- **Role reveal**: All players' Anima Orb colors (roles) are shown.
- **Winning team**: The team that achieved their win condition is announced.
**Per-role stats:**
- **RED (Hunter)**: number of players killed.
- **WHITE / BLACK (Civilians)**: number of tasks completed.

**Common stats (all roles):**
- Match outcome (win / lose).
- Survival time (time alive before death, or full match duration if survived).

> After the result screen the same party **returns to the Lobby** (see Section 4 — Match Flow).

---

## 16. Art Direction

- **Style**: **2.5D perspective pixel art** — top-down with a subtle 3/4 perspective angle, giving a sense of depth and height while remaining tile-based. Primary reference: Don't Starve / Don't Starve Together.
- **Perspective rules**:
  - The camera angle is a fixed oblique/isometric-leaning top-down view — **not flat overhead, not full isometric**.
  - Environment objects (trees, rocks, ruins, walls) have visible **side faces** or **height shadow** to convey volume.
  - Characters are drawn facing the camera at a slight downward angle — visible head, torso, and feet.
  - Wall tiles have a top face + front face; floor tiles are flat. This is standard Don't Starve / Stardew Valley tile convention.
  - Sprites are **not rotated** at runtime — directional walking uses separate sprite frames per direction (typically: down, up, left, right = 4 directions, or 8 if budget allows).
- **Primary reference**: Don't Starve / Don't Starve Together — dark, whimsical, slightly grotesque.
- **Secondary references**: Among Us (social deduction UX), Stardew Valley (2.5D tile convention).
- **Character style**: Stylized, slightly exaggerated proportions. Dark and eerie, not cute.
- **Environment**: A surreal overgrown wilderness — dark forests, crumbling structures, arcane formations. Environmental details (wrecked vehicles, shattered glass, skid marks) are present as easter eggs but never explained.
- The world feels real but wrong. Nothing needs explanation. The darkness and monsters are simply part of it.
- **Day/Night cycle UI**: A clock-like dial in a corner of the screen shows the current phase. When the cycle transitions to Night, the screen darkens with a gradual fade and a sound effect plays (SFX: TBD).
- **Target resolution**: QHD (2560×1440) as baseline. Players can change resolution in settings.
- **Nickname**: Steam nickname imported by default. Players may change their in-game display name within the game.

> 🔲 TBD: Pixel font usage, UI color palette, exact tile size (16×16 or 32×32 — to be confirmed against chosen asset pack).

# RED HAT — Scene Node Specifications

> **This document changes as scene designs evolve.**
> The agent must follow these node structures exactly when creating or modifying scenes.
> `@onready` variable names in scripts must match the node names defined here.
> Update this file whenever a scene's node tree is intentionally changed.

---

## `scenes/game/player.tscn`

```
CharacterBody2D  (player_controller.gd)
├── Sprite2D            "Sprite"        ← character pixel sprite
├── AnimationPlayer     "AnimationPlayer"
├── CollisionShape2D    "Collision"     ← capsule or rect
├── Node2D              "HelmetAnchor"  ← helmet position offset above head
│   └── [helmet.tscn instance]  "Helmet"
├── Label               "NameLabel"     ← displays player nickname above head
├── Label               "SpeechBubble"  ← shows text chat above head; hidden when off-screen (node visibility)
└── AudioStreamPlayer2D "AudioPlayer"   ← shot_or_explosion.wav
```

---

## `scenes/game/helmet.tscn`

```
Node2D  (no script — controlled by helmet_system.gd from parent)
└── Sprite2D  "HelmetSprite"  ← LED tint applied here via modulate
```

---

## `scenes/game/world.tscn`

```
Node2D  (game_manager.gd)
├── TileMapLayer        "TileMap"         ← floor and wall tiles
├── Node2D              "SpawnPoints"     ← children are Marker2D per spawn slot
│   ├── Marker2D        "Spawn0"
│   ├── Marker2D        "Spawn1"
│   └── ...             (up to Spawn7 for 8 players)
├── Node2D              "PortalLocations" ← all possible portal spots (children are Marker2D); one unlocks at runtime
├── Node2D              "ActivePortal"    ← Area2D + Sprite2D instantiated here when portal unlocks (portal_system.gd)
├── Node2D              "Monsters"        ← monster instances added/removed here at runtime
├── Node2D              "MonsterSpawns"   ← Marker2D nodes defining monster initial positions
├── Node2D              "ItemSpawns"      ← Marker2D nodes where items can appear at game start
├── Node2D              "Players"         ← player instances added here at runtime
├── Node2D              "CCTVCameras"     ← fixed Camera2D nodes for spectator mode
└── CanvasLayer         "HUD"
    └── [hud.tscn instance]  "HudRoot"
```

---

## `scenes/ui/lobby.tscn`

```
Control  (lobby_controller.gd)
├── ColorRect           "Background"
├── ScrollContainer     "PlayerListScroll"
│   └── VBoxContainer   "PlayerList"      ← player rows added here at runtime
├── HBoxContainer       "BottomBar"
│   ├── Button          "ReadyButton"
│   └── Label           "StatusLabel"
├── PanelContainer      "MapPanel"
│   └── VBoxContainer   "MapPanelContent"
│       ├── Label       "MapNameLabel"    ← shows currently selected map name
│       └── Button      "MapSelectButton" ← visible to host only; opens map picker overlay
└── CanvasLayer         "ChatLayer"
    ├── RichTextLabel   "ChatLog"
    └── LineEdit        "ChatInput"
```

---

## `scenes/ui/hud.tscn`

```
CanvasLayer  (hud_controller.gd)
└── Control  "HudRoot"
    ├── Label         "TimerLabel"       ← detonation countdown (MM:SS)
    ├── Label         "TaskCountLabel"   ← "X / Y tasks complete"
    ├── Label         "DayNightLabel"    ← shows current phase ("Day" / "Night") and phase timer
    └── Label         "ProximityWarning" ← shown when crowd drain penalty is active

> Note: Speech bubbles are **NOT** a HUD element. They are Label nodes on each player.tscn instance, hidden automatically when off-screen.
```

---

## `scenes/result/result_screen.tscn`

```
Control  (no script initially — populated from GameState at scene load)
├── Label         "WinnerLabel"   ← "WHITE escaped" / "Hunter wins"
├── GridContainer "RoleGrid"      ← one cell per player, reveals helmet color
└── Button        "ReturnButton"  ← returns to lobby / main menu
```

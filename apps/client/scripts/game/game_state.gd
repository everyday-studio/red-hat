## game_state.gd
## Global singleton that stores persistent state shared across all scenes.
## Holds authentication credentials, match context, and player identity.
##
## Authority: both
## AutoLoad: yes
extends Node

## JWT token received after Steam authentication.
var jwt: String = ""

## Steam ID of the local player (string to preserve 64-bit precision).
var steam_id: String = ""

## Display name of the local player.
var player_name: String = ""

## Current match ID assigned by the backend.
var match_id: String = ""

## Peer ID assigned by Godot's multiplayer layer for this client.
var local_peer_id: int = 0

## Emitted when the player has been authenticated with the backend.
signal authenticated

## Emitted when the matchmaker has assigned a match to this player.
signal match_found(match_id: String)


func _ready() -> void:
	if not Engine.has_singleton("Steam"):
		push_error("[GameState] Steam singleton not found — make sure you are using the GodotSteam editor")
		return
	var steam: Object = Engine.get_singleton("Steam")
	var init: Dictionary = steam.steamInitEx()
	if init["status"] != 0:
		push_error("[GameState] Steam failed to initialize: " + str(init["verbal"]))
	else:
		print("[GameState] Steam initialized — app_id: %d, verbal: %s" % [steam.getAppID(), init["verbal"]])


func _process(_delta: float) -> void:
	if Engine.has_singleton("Steam"):
		Engine.get_singleton("Steam").run_callbacks()


## Returns true if the player holds a valid JWT (i.e. is logged in).
func is_authenticated() -> bool:
	return jwt != ""


## Clears all per-session state. Call on disconnect or logout.
func clear_session() -> void:
	jwt = ""
	match_id = ""
	local_peer_id = 0

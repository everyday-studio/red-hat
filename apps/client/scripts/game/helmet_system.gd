## helmet_system.gd
## Manages helmet LED color synchronization across clients.
##
## SECURITY: The local player's own helmet color is NEVER stored, rendered,
## or branched on. The server sends render commands to other peers only.
##
## Authority: server (color resolution) / client (rendering for other peers)
## AutoLoad: no
extends Node

## Maps peer_id → Color for all OTHER players visible to this client.
## The local player's entry is never inserted here.
var _peer_colors: Dictionary = {}

@onready var _helmet_sprite = $HelmetSprite


## Called by the server to deliver another player's helmet color.
## If peer_id matches the local player, the packet is silently dropped.
@rpc("authority", "call_remote", "reliable")
func receive_helmet_color(peer_id: int, color: Color) -> void:
	if peer_id == multiplayer.get_unique_id():
		return  # own color → silently drop, never store or render
	_apply_helmet_color(peer_id, color)


## Applies the visual color to the helmet sprite of the given peer.
func _apply_helmet_color(peer_id: int, color: Color) -> void:
	_peer_colors[peer_id] = color
	_helmet_sprite.color = color


## Removes a peer's color entry when they disconnect.
func remove_peer(peer_id: int) -> void:
	_peer_colors.erase(peer_id)

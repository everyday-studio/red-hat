## gun_system.gd
## Handles fire-request routing and server-side bullet/detonation resolution.
##
## SECURITY: The client only requests a fire action. The server decides whether
## it results in a real bullet or a self-detonation, because only the server
## holds the shooter's role. The client never branches on its own role.
##
## Authority: server (resolution) / client (fire request only)
## AutoLoad: no
extends Node

## Secretly held role map: peer_id → "hunter" | "white" | "civilian".
## Populated on the server only. Never sent to clients.
var _roles: Dictionary = {}

## Emitted when a real bullet is fired (visible to all clients).
signal bullet_fired(shooter_peer_id: int, origin: Vector2, direction: Vector2)

## Emitted when a helmet detonation is triggered (visible to all clients).
signal helmet_detonated(peer_id: int)


## SERVER ONLY: Resolves a fire request from the given peer.
## Fires a real bullet if the peer is the Hunter; otherwise self-detonates.
func _server_resolve_fire(peer_id: int) -> void:
	if not multiplayer.is_server():
		return
	# roles dict is server-only — never exposed to clients
	var role: String = _roles.get(peer_id, "")
	if role == "hunter":
		_fire_real_bullet(peer_id)
	else:
		_detonate_helmet(peer_id)


func _fire_real_bullet(_peer_id: int) -> void:
	pass


func _detonate_helmet(_peer_id: int) -> void:
	pass


## SERVER ONLY: Assigns a role to a peer. Called during match setup.
func assign_role(peer_id: int, role: String) -> void:
	if not multiplayer.is_server():
		return
	_roles[peer_id] = role

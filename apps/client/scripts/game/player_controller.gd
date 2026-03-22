## player_controller.gd
## Handles local player input, movement, and RPC calls for actions.
## All physics/input processing is guarded to the authority peer only.
##
## Authority: client (local peer only)
## AutoLoad: no
extends CharacterBody2D

const MOVE_SPEED: float = 120.0

## Peer ID this node belongs to.
var peer_id: int = 0


func _ready() -> void:
	set_multiplayer_authority(peer_id)


func _physics_process(delta: float) -> void:
	if not is_multiplayer_authority():
		return
	_handle_movement(delta)


func _input(event: InputEvent) -> void:
	if not is_multiplayer_authority():
		return
	if event.is_action_pressed("fire"):
		_request_fire()


func _handle_movement(_delta: float) -> void:
	var direction: Vector2 = Input.get_vector("move_left", "move_right", "move_up", "move_down")
	velocity = direction * MOVE_SPEED
	move_and_slide()


## Notifies the server that this player wants to fire.
## The server determines the outcome (real bullet or self-detonation).
@rpc("any_peer", "call_local", "reliable")
func _request_fire() -> void:
	if not multiplayer.is_server():
		return
	# Server-side resolution is handled in gun_system.gd
	pass

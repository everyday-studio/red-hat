## proximity_system.gd
## Detects when more than 2 players are within a defined radius.
## When the threshold is exceeded, the detonation timer drains 2x faster.
##
## Authority: both
## AutoLoad: no
extends Node

## Radius in pixels within which players are considered "in proximity".
const PROXIMITY_RADIUS: float = 96.0

## Minimum number of players in proximity to trigger the drain penalty.
const CROWD_THRESHOLD: int = 3

## Emitted when the crowd threshold is first exceeded.
signal crowding_started

## Emitted when the number of nearby players drops back below the threshold.
signal crowding_ended

var _is_crowded: bool = false


func _physics_process(_delta: float) -> void:
	pass


## Returns all player nodes within PROXIMITY_RADIUS of the given position.
func get_nearby_players(_position: Vector2) -> Array[Node]:
	return []

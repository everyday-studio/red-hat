## cctv_controller.gd
## Manages the CCTV spectator view for players who have been KO'd.
## Cycles through fixed camera positions that have blind spots by design.
##
## Authority: client
## AutoLoad: no
extends Control

## List of camera positions available for the spectator to switch between.
var _camera_positions: Array[Vector2] = []

## Index of the currently displayed camera feed.
var _current_camera_index: int = 0

## Emitted when the player switches to a different camera.
signal camera_switched(index: int)


func _ready() -> void:
	pass


## Switches to the next available CCTV camera.
func next_camera() -> void:
	if _camera_positions.is_empty():
		return
	_current_camera_index = (_current_camera_index + 1) % _camera_positions.size()
	camera_switched.emit(_current_camera_index)


## Switches to the previous available CCTV camera.
func previous_camera() -> void:
	if _camera_positions.is_empty():
		return
	_current_camera_index = (_current_camera_index - 1 + _camera_positions.size()) % _camera_positions.size()
	camera_switched.emit(_current_camera_index)

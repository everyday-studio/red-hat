## game_manager.gd
## Coordinates the overall game lifecycle: round start, task completion,
## timer management, and end-game conditions.
##
## Authority: both
## AutoLoad: no
extends Node

## Base duration of the shared detonation timer in seconds.
const BASE_TIMER: float = 600.0

## Timer drain multiplier when more than 2 players are in proximity.
const PROXIMITY_DRAIN_MULTIPLIER: float = 2.0

## Remaining seconds on the shared detonation timer.
var detonation_timer: float = BASE_TIMER

## Whether the extraction gate is currently open.
var is_gate_open: bool = false

## Emitted each second with the updated timer value.
signal timer_updated(remaining: float)

## Emitted when the timer has reached zero.
signal timer_expired

## Emitted when the extraction gate opens (all tasks complete).
signal gate_opened

## Emitted when the match ends.
signal match_ended(result: Dictionary)


func _process(_delta: float) -> void:
	pass


## Extends the detonation timer by the given number of seconds.
func extend_timer(_seconds: float) -> void:
	pass


## Called when all tasks are completed; opens the extraction gate.
func _on_all_tasks_completed() -> void:
	pass

## matchmaking.gd
## Joins the matchmaking queue and polls GET /matchmaking/status every 0.5s
## until a match is found or the process is cancelled.
##
## Authority: client
## AutoLoad: no
extends Node

const POLL_INTERVAL: float = 0.5

## Emitted when the server has found a match and returned a match_id.
signal match_found(match_id: String)

## Emitted when matchmaking fails or the request is rejected.
signal matchmaking_failed(reason: String)

var _poll_timer: Timer


## Sends POST /matchmaking/join, then starts polling for a result.
func join_queue() -> void:
	pass


## Cancels an in-progress matchmaking attempt.
func cancel() -> void:
	_stop_polling()


func _poll_status() -> void:
	pass


func _stop_polling() -> void:
	if is_instance_valid(_poll_timer):
		_poll_timer.stop()
		_poll_timer.queue_free()
		_poll_timer = null

## task_system.gd
## Tracks task completion state and notifies the game manager when all
## tasks are done so the extraction gate can open.
##
## Authority: both
## AutoLoad: no
extends Node

## Total number of tasks in the current match.
var total_tasks: int = 0

## Number of tasks completed so far.
var completed_tasks: int = 0

## Emitted when a single task is completed.
signal task_completed(task_id: String)

## Emitted when every task in the match has been completed.
signal all_tasks_completed


## Marks the given task as completed and checks for full completion.
@rpc("any_peer", "call_local", "reliable")
func complete_task(task_id: String) -> void:
	completed_tasks += 1
	task_completed.emit(task_id)
	if completed_tasks >= total_tasks:
		all_tasks_completed.emit()


## Returns a value between 0.0 and 1.0 representing overall task progress.
func get_progress() -> float:
	if total_tasks == 0:
		return 0.0
	return float(completed_tasks) / float(total_tasks)

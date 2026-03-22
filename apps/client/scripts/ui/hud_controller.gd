## hud_controller.gd
## Updates the in-game HUD: detonation timer, task progress bar,
## and proximity warning indicator.
##
## Authority: client
## AutoLoad: no
extends Control

@onready var _label_timer: Label = $LabelTimer
@onready var _bar_tasks: ProgressBar = $BarTasks
@onready var _icon_proximity: TextureRect = $IconProximity


func _ready() -> void:
	pass


## Updates the displayed timer value. Called every second by GameManager.
func update_timer(remaining: float) -> void:
	var minutes: int = floori(remaining / 60.0)
	var seconds: int = int(remaining) % 60
	_label_timer.text = "%02d:%02d" % [minutes, seconds]


## Updates the task progress bar (value: 0.0 – 1.0).
func update_task_progress(value: float) -> void:
	_bar_tasks.value = value * 100.0


## Shows or hides the crowd proximity warning icon.
func set_crowding_warning(is_visible: bool) -> void:
	_icon_proximity.visible = is_visible

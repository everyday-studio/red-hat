## lobby_controller.gd
## Controls the lobby scene: displays the player list, ready states,
## and triggers matchmaking when the local player clicks "Find Match".
##
## Authority: client
## AutoLoad: no
extends Control

@onready var _btn_find_match: Button = $BtnFindMatch
@onready var _label_status: Label = $LabelStatus

var _matchmaking: Node = null


func _ready() -> void:
	_btn_find_match.pressed.connect(_on_find_match_pressed)


func _on_find_match_pressed() -> void:
	pass

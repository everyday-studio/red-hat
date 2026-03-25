## steam_auth_test.gd
## Manual test scene for the Steam authentication flow.
## Run this scene with F5 and press the button to verify auth end-to-end.
##
## Authority: client
## AutoLoad: no
extends Control

@onready var _status_label: Label = $VBoxContainer/StatusLabel
@onready var _auth_button: Button = $VBoxContainer/AuthButton

var _steam_auth: Node = null


func _ready() -> void:
	_steam_auth = preload("res://scripts/network/steam_auth.gd").new()
	add_child(_steam_auth)
	_steam_auth.auth_succeeded.connect(_on_auth_succeeded)
	_steam_auth.auth_failed.connect(_on_auth_failed)
	_auth_button.pressed.connect(_on_auth_button_pressed)
	print("[SteamAuthTest] Ready — click the button to begin Steam authentication")


func _on_auth_button_pressed() -> void:
	_status_label.text = "인증 중..."
	_auth_button.disabled = true
	print("[SteamAuthTest] Button pressed — starting authentication")
	_steam_auth.authenticate()


func _on_auth_succeeded() -> void:
	_status_label.text = "✅ 인증 성공!\n이름: %s\nSteam ID: %s" % [GameState.player_name, GameState.steam_id]
	_auth_button.disabled = false
	print("[SteamAuthTest] SUCCESS — player_name='%s'  steam_id=%s" % [GameState.player_name, GameState.steam_id])


func _on_auth_failed(reason: String) -> void:
	_status_label.text = "❌ 인증 실패\n" + reason
	_auth_button.disabled = false
	push_warning("[SteamAuthTest] FAILED — " + reason)

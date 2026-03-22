## steam_auth.gd
## Handles Steam ticket generation and backend authentication.
## On success, writes JWT, steam_id, and player_name into GameState.
##
## Authority: client
## AutoLoad: no
extends Node

## Emitted when authentication completes successfully.
signal auth_succeeded

## Emitted when authentication fails for any reason.
signal auth_failed(reason: String)


## Begins the Steam authentication flow.
## 1. Requests an auth ticket from GodotSteam.
## 2. POSTs the ticket to /auth/steam.
## 3. Stores the returned JWT in GameState.
func authenticate() -> void:
	pass

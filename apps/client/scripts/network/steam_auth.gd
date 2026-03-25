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

var _pending_ticket_handle: int = 0


## Begins the Steam authentication flow.
## 1. Requests an auth ticket from GodotSteam.
## 2. POSTs the ticket to /auth/steam.
## 3. Stores the returned JWT in GameState.
func authenticate() -> void:
	if not Engine.has_singleton("Steam"):
		push_error("[SteamAuth] Steam singleton not found — make sure you are using the GodotSteam editor")
		auth_failed.emit("Steam singleton not found")
		return

	var steam: Object = Engine.get_singleton("Steam")

	if not steam.isSteamRunning():
		push_warning("[SteamAuth] Steam client is not running")
		auth_failed.emit("Steam client is not running")
		return

	# getSteamID() returns 0 if Steam is running but not logged in / not fully initialized
	var steam_id_check: int = steam.getSteamID()
	print("[SteamAuth] Steam state — isSteamRunning: true, getSteamID: %d, personaName: '%s'" % [
		steam_id_check, steam.getPersonaName()
	])

	if steam_id_check == 0:
		push_warning("[SteamAuth] getSteamID() returned 0 — Steam is running but not fully initialized or not logged in")
		auth_failed.emit("Steam not fully initialized (getSteamID = 0)")
		return

	steam.get_ticket_for_web_api.connect(_on_ticket_received, CONNECT_ONE_SHOT)
	# Pass empty string as identity — using a custom identity requires server-side registration
	_pending_ticket_handle = steam.getAuthTicketForWebApi("")
	print("[SteamAuth] Requesting auth ticket (handle: %d)" % _pending_ticket_handle)

	# handle 0 = k_HAuthTicketInvalid — Steam rejected the request immediately
	if _pending_ticket_handle == 0:
		steam.get_ticket_for_web_api.disconnect(_on_ticket_received)
		push_warning("[SteamAuth] getAuthTicketForWebApi returned invalid handle (0) — is steam_appid.txt present and Steam running?")
		auth_failed.emit("Failed to request auth ticket (invalid handle)")


# GodotSteam signal: get_ticket_for_web_api(auth_ticket, result, ticket_size, ticket_buffer)
func _on_ticket_received(auth_ticket: int, result: int, _ticket_size: int, ticket: PackedByteArray) -> void:
	# k_EResultOK = 1
	if result != 1:
		push_warning("[SteamAuth] Ticket generation failed (result: %d)" % result)
		auth_failed.emit("Ticket generation failed (result: %d)" % result)
		return

	print("[SteamAuth] Ticket received — handle: %d, size: %d bytes" % [auth_ticket, ticket.size()])

	var steam: Object = Engine.get_singleton("Steam")
	var steam_id: String = str(steam.getSteamID())
	var ticket_hex: String = _bytes_to_hex(ticket)

	var response: Dictionary = await ApiClient.post("/auth/steam", {
		"ticket": ticket_hex,
		"steam_id": steam_id,
	})

	# Always cancel the ticket after use regardless of result
	steam.cancelAuthTicket(auth_ticket)

	if response.has("error"):
		push_warning("[SteamAuth] Auth failed: " + str(response["error"]))
		auth_failed.emit(str(response["error"]))
		return

	if not response.has("token"):
		push_warning("[SteamAuth] Server response missing 'token' field")
		auth_failed.emit("Invalid server response")
		return

	GameState.jwt = response["token"]
	GameState.steam_id = steam_id
	GameState.player_name = steam.getPersonaName()
	GameState.authenticated.emit()

	print("[SteamAuth] Authenticated! name='%s'  steam_id=%s" % [GameState.player_name, GameState.steam_id])
	auth_succeeded.emit()


func _bytes_to_hex(bytes: PackedByteArray) -> String:
	var hex: String = ""
	for b: int in bytes:
		hex += "%02x" % b
	return hex

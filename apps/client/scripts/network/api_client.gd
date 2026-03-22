## api_client.gd
## Centralized HTTP client. ALL requests to the backend must go through here.
## Direct instantiation of HTTPRequest outside this file is strictly forbidden.
##
## Authority: client
## AutoLoad: no
extends Node

const BASE_URL: String = "http://localhost:8080"
const REQUEST_TIMEOUT: float = 10.0


## Sends a POST request to the given endpoint with a JSON body.
## Returns the decoded JSON response as a Dictionary.
## On failure, returns {"error": <reason: String>}.
func post(endpoint: String, body: Dictionary) -> Dictionary:
	return await _request(HTTPClient.METHOD_POST, endpoint, body)


## Sends a GET request to the given endpoint.
## Returns the decoded JSON response as a Dictionary.
## On failure, returns {"error": <reason: String>}.
func get_request(endpoint: String) -> Dictionary:
	return await _request(HTTPClient.METHOD_GET, endpoint, {})


func _request(method: HTTPClient.Method, endpoint: String, body: Dictionary) -> Dictionary:
	var http: HTTPRequest = HTTPRequest.new()
	add_child(http)
	http.timeout = REQUEST_TIMEOUT

	var headers: PackedStringArray = ["Content-Type: application/json"]
	if GameState.is_authenticated():
		headers.append("Authorization: Bearer " + GameState.jwt)

	var body_string: String = JSON.stringify(body) if not body.is_empty() else ""
	var send_error: int = http.request(BASE_URL + endpoint, headers, method, body_string)

	if send_error != OK:
		http.queue_free()
		return {"error": "Failed to send request (code: %d)" % send_error}

	# [result, response_code, headers, body]
	var response: Array = await http.request_completed
	http.queue_free()

	var result: int = response[0]
	var response_code: int = response[1]
	var raw_body: PackedByteArray = response[3]

	if result != HTTPRequest.RESULT_SUCCESS:
		return {"error": "HTTP request failed (result: %d)" % result}

	var json: JSON = JSON.new()
	if json.parse(raw_body.get_string_from_utf8()) != OK:
		return {"error": "Failed to parse response as JSON"}

	if response_code >= 400:
		var data = json.data
		if data is Dictionary and data.has("error"):
			return {"error": data["error"]}
		return {"error": "Server returned status %d" % response_code}

	if json.data is Dictionary:
		return json.data
	return {"data": json.data}

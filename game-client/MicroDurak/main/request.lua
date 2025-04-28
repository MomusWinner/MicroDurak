local save = require("main.save_manager")

local HOST_AUTH = "http://localhost:8080/auth/"
local HOST_GAME_SERVER = "ws://localhost:7070/game-manager"
local HOST_MATCHMAKER = "ws://localhost:3000/matchmaker/find-match"

local M = {}

---@param user table
---@param handler function
function M.register(user, handler)
	assert(user.name)
	assert(user.age)
	assert(user.email)
	assert(user.password)

	print("REQUEST register")
	local url = HOST_AUTH .. "register"
	local method = "POST"
	local header = { ["Content-Type"] = "application/json" }
	local body = json.encode(user)

	http.request(url, method, handler, header, body)
end

---@param user table
---@param handler function
function M.login(user, handler)
	assert(user.email)
	assert(user.password)

	print("REQUEST login")
	local url = HOST_AUTH .. "login"
	local method = "POST"
	local header = { ["Content-Type"] = "application/json" }
	local body = json.encode(user)
	print(body)

	http.request(url, method, handler, header, body)
end

---@param game_id string
---@param callback function(self:object, connection:object, data:table)
---@return connection
function M.connect_game_server(game_id, callback)
	local url = HOST_GAME_SERVER .. "/" .. game_id
	local token = save.get().token

	local params = {
		timeout = 3000,
		headers = "Authorization: " .. token,
	}

	return websocket.connect(url, params, callback)
end

---@param callback function(self:object, connection:object, data:table)
---@return connection
function M.connect_matchmaker(callback)
	local url = HOST_MATCHMAKER
	local token = save.get().token

	local params = {
		timeout = 3000,
		headers = "Authorization: " .. token,
	}

	return websocket.connect(url, params, callback)
end

return M

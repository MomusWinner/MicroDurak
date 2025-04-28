local requests = require("main.request")
local s = require("main.game_state")

local M = {}

M.on_disconnect = nil

local connection = nil
local observers = {}

local LOG_PREFIX = "[FIND_MATCH] "

local function call_observers(game_id)
	print(#observers)
	for _, callback in ipairs(observers) do
		callback(game_id)
	end
	observers = {}
end

local function websocket_callback(self, conn, data)
	if data.event == websocket.EVENT_DISCONNECTED then
		print(
			LOG_PREFIX
				.. "Disconnected: "
				.. tostring(conn)
				.. " Code: "
				.. data.code
				.. " Message: "
				.. tostring(data.message)
		)
		connection = nil
		if M.on_disconnect then
			M.on_disconnect()
		end
	elseif data.event == websocket.EVENT_CONNECTED then
		print(LOG_PREFIX .. "Connected: " .. tostring(conn))
	elseif data.event == websocket.EVENT_ERROR then
		print(LOG_PREFIX .. "Error: '" .. tostring(data.message) .. "'")
		if data.handshake_response then
			print(LOG_PREFIX .. "Handshake response status: '" .. tostring(data.handshake_response.status) .. "'")
			for key, value in pairs(data.handshake_response.headers) do
				print(LOG_PREFIX .. "Handshake response header: '" .. key .. ": " .. value .. "'")
			end
			print(LOG_PREFIX .. "Handshake response body: '" .. tostring(data.handshake_response.response) .. "'")
		end
	elseif data.event == websocket.EVENT_MESSAGE then
		print(LOG_PREFIX .. "Receiving: '" .. tostring(data.message) .. "'")
		local status = json.decode(data.message)
		pprint(status)
		if status.status == "pending" or status.status == "found_group" then
			websoket.send(connection, "ok")
		elseif status.status == "created" then
			s.game_id = status.game_id
			call_observers(s.game_id)
		end
	end
end

function M.connect()
	connection = requests.connect_matchmaker(websocket_callback)
	print(LOG_PREFIX .. "Connecting ...")
end

---@param callback function(game_id:string)
function M.subscribeOnFound(callback)
	table.insert(observers, callback)
end

return M

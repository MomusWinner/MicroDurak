local requests = require("main.request")
local s = require("main.game_state")
local save = require("main.save_manager")

local M = {}

M.ACTION_READY = "ACTION_READY"
M.ACTION_ATTACK = "ACTION_ATTACK"
M.ACTION_DEFEND = "ACTION_DEFEND"
M.ACTION_END_ATTACK = "ACTION_END_ATTACK"
M.ACTION_TAKE_ALL_CARDS = "ACTION_TAKE_ALL_CARDS"
M.ACTION_CHECK_ATTACK_TIMER = "ACTION_CHECK_ATTACK_TIMER"
M.ACTION_CHECK_DEFEND_TIMER = "ACTION_CHECK_DEFEND_TIMER"

M.ERROR_EMPTY = ""
M.ERROR_SERVER = "SERVER_ERROR"
M.ERROR_BAD_REQUEST = "BAD_REQUEST"
M.ERROR_USER_ALREADY_READY = "USER_ALREADY_READY"
M.ERROR_NOT_YOUR_TURN = "NOT_YOUR_TURN"
M.ERROR_INCORRECT_CARD = "INCORRECT_CARD"
M.ERROR_USER_NO_HAS_CARD = "USER_NO_HAS_CARD"
M.ERROR_ATTACK_TIME_OVER = "ATTACK_TIME_OVER"
M.ERROR_DEFEND_TIME_OVER = "DEFEND_TIME_OVER"
M.ERROR_NO_SAME_RANK_CARD_IN_TABLE = "NO_SAME_RANK_CARD_IN_TABLE"
M.ERROR_NOT_FOUND_CART_ON_TABLE = "NOT_FOUND_CART_ON_TABLE"
M.ERROR_TARGET_CARD_GREATER_THEN_YOUR = "TARGET_CARD_GREATER_THEN_YOUR"
M.ERROR_GAME_SHOULD_BE_STARTED = "GAME_SHOULD_BE_STARTED"
M.ERROR_CANNOT_END_ATTACK_IN_FIRST_TURN = "CANNOT_END_ATTACK_IN_FIRST_TURN"
M.ERROR_ALL_CARD_SHOULD_BE_BEAT_OFF_BEFORE_END_ATTACK = "ALL_CARD_SHOULD_BE_BEAT_OFF_BEFORE_END_ATTACK"
M.ERROR_TABLE_HOLDS_ONLY_SIX_CARDS = "TABLE_HOLDS_ONLY_SIX_CARDS"
M.ERROR_DEFENDER_NO_CARDS = "DEFENDER_NO_CARDS"
M.ERROR_ALREADY_END_ATTACK = "ALREADY_END_ATTACK" -- TODO: implement
M.ERROR_UNREGISTERED_ACTION = "UNREGISTERED_ACTION"

M.EVENT_START = "EVENT_START"
M.EVENT_READY = "EVENT_READY"
M.EVENT_ATTACK = "EVENT_ATTACK"
M.EVENT_DEFEND = "EVENT_DEFEND"
M.EVENT_END_ATTACK = "EVENT_END_ATTACK"
M.EVENT_TAKE_ALL_CARDS = "EVENT_TAKE_ALL_CARDS"
M.EVENT_ATTACK_TIMER_NOT_COMPLETED = "ATTACK_TIMER_NOT_COMPLETED"
M.EVENT_DEFEND_TIMER_NOT_COMPLETED = "DEFEND_TIMER_NOT_COMPLETED"
M.EVENT_ATTACK_TIMER_COMPLETED = "ATTACK_TIMER_COMPLETED"
M.EVENT_DEFEND_TIMER_COMPLETED = "DEFEND_TIMER_COMPLETED"
M.EVENT_END_GAME = "END_GAME"

M.on_connect = nil
M.on_disconnect = nil

local LOG_PREFIX = "[GAME WS] "

local connection = nil
local subscribers = {}
local sentCommands = {}

local function log_v(message)
	msg.post("game:/debug_gui", "message", { text = message })
end

local function log(message)
	if type(message) == "table" then
		print(LOG_PREFIX)
		pprint(message)
	else
		print(LOG_PREFIX .. message)
	end
end

local function deepEqual(t1, t2)
	if type(t1) ~= type(t2) then
		return false
	end
	if type(t1) ~= "table" then
		return t1 == t2
	end
	if #t1 ~= #t2 then
		return false
	end

	for k, v in pairs(t1) do
		if type(v) == "table" then
			if not deepEqual(v, t2[k]) then
				return false
			end
		else
			if v ~= t2[k] then
				return false
			end
		end
	end

	return true
end

local function init_command(command)
	command.user_id = save.get().user_id
	command.game_id = s.game_id
	return command
end

function M.new_ready_command()
	return init_command({ action = M.ACTION_READY })
end

---@param card Card
function M.new_attack_command(card)
	return init_command({
		action = M.ACTION_ATTACK,
		card = {
			suit = card.suit,
			rank = card.rank,
		},
	})
end

---@param user_card Card
---@param target_card Card
function M.new_defend_command(user_card, target_card)
	return init_command({
		action = M.ACTION_DEFEND,
		user_card = {
			suit = user_card.suit,
			rank = user_card.rank,
		},
		target_card = {
			suit = target_card.suit,
			rank = target_card.rank,
		},
	})
end

function M.new_take_all_cards_command()
	return init_command({ action = M.ACTION_TAKE_ALL_CARDS })
end

function M.new_end_attack_command()
	return init_command({ action = M.ACTION_END_ATTACK })
end

function M.new_check_attack_timer_command()
	return init_command({ action = M.ACTION_CHECK_ATTACK_TIMER })
end

function M.new_check_defend_timer_command()
	return init_command({ action = M.ACTION_CHECK_DEFEND_TIMER })
end

function M.subscribe(event, callback)
	if not subscribers[event] then
		subscribers[event] = {}
	end

	table.insert(subscribers[event], callback)
end

local function invoke_event(event, response)
	if subscribers[event] then
		for _, callback in ipairs(subscribers[event]) do
			callback(response)
		end
	end
end

local function handle_event(event)
	log_v("Receive EVENT: " .. event.event)
	log(event)
	if event.event ~= nil then
		invoke_event(event.event, event)
	end
end

local function handle_command_response(command_response)
	for i, item in ipairs(sentCommands) do
		if deepEqual(item.command, command_response.command) then
			local error_text = ""
			if command_response.error ~= M.ERROR_EMPTY then
				error_text = "\nError: " .. command_response.error
			end
			log_v("Receive COMMAND RESPONSE: " .. command_response.command.action .. error_text)
			item.callback(command_response)
			table.remove(sentCommands, i)
			return
		end
	end
end

function M.send_command(command, callback)
	if connection == nil then
		error("Websocket not initialized or closed")
	end
	log("Send command: ")
	log_v("Send COMMAND: " .. command.action)
	log(command)
	local msg = json.encode(command)
	local item = { command = command, callback = callback }
	pprint(item)
	table.insert(sentCommands, item)
	websocket.send(connection, msg)
end

local function callback(_, conn, data)
	if data.event == websocket.EVENT_DISCONNECTED then
		log("Disconnected: " .. tostring(conn) .. " Code: " .. data.code .. " Message: " .. tostring(data.message))
		log_v("Disconnected")
		connection = nil
		if M.on_disconnect then
			M.on_disconnect()
		end
	elseif data.event == websocket.EVENT_CONNECTED then
		log("Connected: " .. tostring(conn))
		log_v("Connected")
		if M.on_connect then
			M.on_connect()
		end
	elseif data.event == websocket.EVENT_ERROR then
		log("Error: '" .. tostring(data.message) .. "'")
		log_v("Error:" .. tostring(data.message))
		if data.handshake_response then
			log("Handshake response status: '" .. tostring(data.handshake_response.status) .. "'")
			for key, value in pairs(data.handshake_response.headers) do
				log("Handshake response header: '" .. key .. ": " .. value .. "'")
			end
			log("Handshake response body: '" .. tostring(data.handshake_response.response) .. "'")
		end
	elseif data.event == websocket.EVENT_MESSAGE then
		log("Receiving: '" .. tostring(data.message) .. "'")
		local msg = json.decode(data.message)
		if msg == nil then
			log("Receive unregistered text message : " .. data.message)
			log_v("Recieve unregisterd text message")
			return
		end
		if msg.messages then
			print("Messages is NOT nil")
			for _, message in ipairs(msg.messages) do
				if message.event ~= nil then
					log("handle event")
					handle_event(message)
				elseif message.command ~= nil then
					log("handle command response")
					handle_command_response(message)
				else
					log("Receive unregistered json message : ")
					log(message)
				end
			end
		else
			log("Messages is nil")
		end
	end
end

function M.connect()
	log("Connecting ...")
	connection = requests.connect_game_server(s.game_id, callback)
end

-- function final(self) end
--
-- function on_input(self, action_id, action)
-- end

return M

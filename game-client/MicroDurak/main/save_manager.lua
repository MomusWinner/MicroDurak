local defsave = require("defsave.defsave")
local s = require("main.game_state")

local M = {}

local default_settings = {
	token = nil,
	user_id = nil,
}

local settings = default_settings

function M.init()
	defsave.appname = "micro_durak"
	defsave.load("config")
	defsave.default_data = {
		config = {
			settings = default_settings,
		},
	}
	defsave.autosave = true
	settings = defsave.get("config", "settings")
end

function M.get()
	return settings
end

function M.set(new_settings)
	settings = new_settings
	if not s.simulate_save then
		defsave.set("config", "settings", settings)
		defsave.save("config")
	end
end

function M.update(dt) end

function M.on_input(action_id, action)
	if action_id == hash("key_c") then
		print("REMOVE SAVE")
		defsave.reset_to_default("config")
		defsave.save("config")
	end
end

return M

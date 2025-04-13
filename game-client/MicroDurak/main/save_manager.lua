local defsave = require("defsave.defsave")

local M = {}

function M.init()
	defsave.appname = "micro_durak"
	defsave.load("config")
	defsave.default_data = {
		config = {
			token = nil,
		},
	}
	defsave.autosave = true
end

function M.get(key)
	return defsave.get("config", key)
end

function M.set(key, value)
	if value == nil then
		print("reset KEY: " .. key)
	else
		print("save KEY: " .. key .. " VALUE: " .. value)
	end
	defsave.set("config", key, value)
	defsave.save("config")
end

function M.update(dt)
	defsave.update(dt)
end

function M.on_input(action_id, action)
	if action_id == hash("key_c") then
		print("REMOVE SAVE")
		defsave.reset_to_default("config")
		defsave.save("config")
	end
end

return M

local M = {}

function M.get_opponent(game_state)
	for _, user in ipairs(game_state.users) do
		if user.id ~= game_state.me.id then
			return user
		end
	end
end

return M

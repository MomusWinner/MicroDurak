local HOST = "http://localhost:8080/auth/"

local M = {}

---@param user table
---@param handler function
function M.register(user, handler)
	assert(user.name)
	assert(user.age)
	assert(user.email)
	assert(user.password)

	print("REQUEST register")
	local url = HOST .. "register"
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
	local url = HOST .. "login"
	local method = "POST"
	local header = { ["Content-Type"] = "application/json" }
	local body = json.encode(user)
	print(boyd)

	http.request(url, method, handler, header, body)
end

return M

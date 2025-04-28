local Card = {}

Card.Suit = {
	HEARTS = 0,
	SPADES = 1,
	CLUBS = 2,
	DIAMONDS = 3,
}

Card.width = 75
Card.height = 100

---@class Card
---@field url string
---@field suit number
---@field rank number
---@field width number
---@field owner_id string
---@param url string
---@param suit number
---@param rank number
---@param owner_id string
function Card:new(url, suit, rank, owner_id)
	local o = {
		url = url,
		suit = suit,
		rank = rank,
		owner_id = owner_id,
	}
	setmetatable(o, self)
	self.__index = self
	return o
end

function Card:init(suit, rank)
	self.suit = suit
	self.rank = rank
end

function Card:equal(other)
	if other.rank == self.ranks and other.soit == self.suit then
		return true
	end
	return false
end

function Card:move(pos)
	go.animate(self.url, "position", go.PLAYBACK_ONCE_FORWARD, pos, go.EASING_LINEAR, 0.4)
end

---@return Card
function Card.create(suit, rank, owner_id)
	local p = vmath.vector3()
	local url = factory.create("core#card_factory", p, nil, 1)
	local card = Card:new(url, suit, rank, owner_id)
	return card
end

---@param x number
---@param y number
function Card:pick(x, y)
	local pos = go.get_position(self.url)
	local x_pick = x >= pos.x - Card.width / 2 and x <= pos.x + Card.width / 2
	local y_pick = y >= pos.y - Card.height / 2 and y <= pos.y + Card.height / 2

	return x_pick and y_pick
end

function Card:show()
	msg.post(self.url, "show", { suit = self.suit, rank = self.rank })
end

function Card:hide()
	msg.post(self.url, "hide")
end

function Card:to_string()
	return "Card S: " .. self.suit .. " R: " .. self.rank .. " Owner: " .. self.owner_id
end

return Card

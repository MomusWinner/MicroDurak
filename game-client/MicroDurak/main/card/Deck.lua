local Card = require("main.card.Card")

local Deck = {}

---@class CardHolder
---@field position any
---@param position any
function Deck:new(position)
	local o = {
		position = position,
		cards = {},
	}
	setmetatable(o, self)
	self.__index = self
	return o
end

function Deck:init(length, trump_suit, trump_rank)
	local trump_card = Card.create(trump_suit, trump_rank)
	local trump_pos = vmath.vector3(self.position)
	trump_pos.x = trump_pos.x - Card.height * 0.3
	go.set_position(trump_pos, trump_card.url)
	go.set_rotation(vmath.quat_rotation_z(math.rad(90)), trump_card.url)
	trump_card:show()
	table.insert(self.cards, trump_card)

	for _ = 1, length - 1 do
		local card = Card.create(1, 6)
		go.set_position(self.position, card.url)
		table.insert(self.cards, card)
	end
end

function Deck:get_card()
	if #self.cards == 0 then
		return
	end

	return table.remove(self.cards, #self.cards)
end

return Deck

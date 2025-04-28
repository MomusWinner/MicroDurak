local Card = require("main.card.Card")
local Table = {}

---@class Table
---@field position any
---@field w_card_offset number
---@field h_card_offset number
---@param position any
---@param w_card_offset number
---@param h_card_offset number
function Table:new(position, w_card_offset, h_card_offset)
	local o = {
		position = position,
		w_card_offset = w_card_offset,
		h_card_offset = h_card_offset,
		cards = {},
		beat_off_cards = {},
	}
	setmetatable(o, self)
	self.__index = self
	return o
end

---@param card Card
---@param i number
function Table:set_card_position(card, i)
	local pos = vmath.vector3(self.position)
	if i == 1 then
		pos.x = pos.x - Card.width - self.w_card_offset
		pos.y = pos.y + Card.height + self.h_card_offset
	elseif i == 2 then
		pos.y = pos.y + Card.height + self.h_card_offset
	elseif i == 3 then
		pos.x = pos.x + Card.width + self.w_card_offset
		pos.y = pos.y + Card.height + self.h_card_offset
	elseif i == 4 then
		pos.x = pos.x - Card.width - self.w_card_offset
		pos.y = pos.y - Card.height - self.h_card_offset
	elseif i == 5 then
		pos.y = pos.y - Card.height - self.h_card_offset
	elseif i == 6 then
		pos.x = pos.x + Card.width + self.w_card_offset
		pos.y = pos.y - Card.height - self.h_card_offset
	end

	card:move(pos)
	return pos
end

function Table:set_beat_off_card_position(card, i)
	local pos = self:set_card_position(card, i)
	pos.x = pos.x + 10
	pos.y = pos.y - 10
	pos.z = 1
	card:move(pos)
end

---@param suit number
---@param rank number
function Table:contain_card(suit, rank)
	for _, card in ipairs(self.cards) do
		if card.suit == suit and card.rank == rank then
			return true
		end
	end

	return false
end

---@param suit number
---@param rank number
function Table:contain_beat_off_card(suit, rank)
	for _, card in ipairs(self.beat_off_cards) do
		if card.suit == suit and card.rank == rank then
			return true
		end
	end

	return false
end

---@param new_card Card
function Table:add_card(new_card)
	for i = 1, 6, 1 do
		if self.cards[i] == nil then
			self.cards[i] = new_card
			self:set_card_position(new_card, i)
			return
		end
	end

	-- self.cards[i] = new_card
	-- table.insert(self.cards, new_card)
	-- self:set_card_position(new_card, #self.cards)
end

function Table:beat_off(card, i)
	if self.beat_off_cards[i] ~= nil then
		return
	end

	self.beat_off_cards[i] = card
	self:set_beat_off_card_position(card, i)
end

---@return boolean
function Table:all_card_beat_off()
	for i, card in ipairs(self.cards) do
		if card ~= nil and self.beat_off_cards[i] == nil then
			return false
		end
	end

	return true
end

---@param suit number
---@param rank number
function Table:get_card(suit, rank)
	for i, card in ipairs(self.cards) do
		if card.rank == rank and card.suit == suit then
			return table.remove(self.cards, i)
		end
	end
end

---@param suit number
---@param rank number
function Table:get_beat_off_card(suit, rank)
	for i, card in ipairs(self.beat_off_cards) do
		if card.rank == rank and card.suit == suit then
			return table.remove(self.beat_off_cards, i)
		end
	end
end

---@param suit number
---@param rank number
---@return number
function Table:find_card(suit, rank)
	for i, card in ipairs(self.cards) do
		if card.rank == rank and card.suit == suit then
			return i
		end
	end

	return -1
end

function Table:pick_card(x, y)
	for i, card in ipairs(self.cards) do
		if card ~= nil and card:pick(x, y) and self.beat_off_cards[i] == nil then
			return i
		end
	end
	return -1
end

function Table:get_all_cards()
	local all_cards = {}
	for i = 1, 6, 1 do
		local pos = vmath.vector3(self.position)
		pos.x = pos.x - 300
		if self.cards[i] ~= nil then
			self.cards[i]:move(pos)
			local card = table.remove(self.cards, i)
			table.insert(all_cards, card)
		end
		if self.beat_off_cards[i] ~= nil then
			self.beat_off_cards[i]:move(pos)
			local card = table.remove(self.beat_off_cards, i)
			table.insert(all_cards, card)
		end
	end

	return all_cards
end

function Table:clean()
	for i = 1, 6, 1 do
		local pos = vmath.vector3(self.position)
		pos.x = pos.x - 300
		if self.cards[i] ~= nil then
			self.cards[i]:move(pos)
			table.remove(self.cards, i)
		end
		if self.beat_off_cards[i] ~= nil then
			self.beat_off_cards[i]:move(pos)
			table.remove(self.beat_off_cards, i)
		end
	end
end

---@param x number
---@param y number
function Table:pick(x, y)
	return y > 200
end

return Table

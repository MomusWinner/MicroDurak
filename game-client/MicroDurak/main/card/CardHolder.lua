local Card = require("main.card.Card")

local CardHolder = {}

---@class CardHolder
---@field druggable boolean
---@field height number
---@field width number
---@field half_width number

---@param height number
---@param width number
---@param druggable boolean
function CardHolder:new(height, width, druggable, on_take_card)
	local o = {
		druggable = druggable,
		height = height,
		width = width,
		on_take_card = on_take_card,
		half_width = sys.get_config_number("display.width") / 2,
		cards = {},
	}
	setmetatable(o, self)
	self.__index = self
	return o
end

function CardHolder:block()
	self.druggable = false
end

function CardHolder:unblock()
	self.druggable = true
end

function CardHolder:update_card_position()
	local length = #self.cards
	if length == 1 then
		local p = vmath.vector3(self.half_width, self.height, 0)
		self.cards[1]:move(p)
		return
	end
	for i, card in ipairs(self.cards) do
		local pos = vmath.vector3()
		pos.y = self.height
		pos.x = ((i - 1) / (length - 1)) * self.width - (self.width / 2) + self.half_width
		card:move(pos)
	end
end

---@param card Card
function CardHolder:add_card(card)
	table.insert(self.cards, card)
	self:update_card_position()
end

---@param suit number
---@param rank number
---@return Card | nil
function CardHolder:get_card(suit, rank)
	for i, card in ipairs(self.cards) do
		if card.rank == rank and card.suit == suit then
			table.remove(self.cards, i)
			self:update_card_position()
			return card
		end
	end
end

---@return Card | nil
function CardHolder:get_first_card()
	if #self.cards == 0 then
		return nil
	end

	return table.remove(self.cards, 1)
end

function CardHolder:on_input(action_id, action)
	if not self.druggable then
		return
	end

	if action_id == hash("touch") then
		if action.pressed then
			for _, card in ipairs(self.cards) do
				if card:pick(action.x, action.y) then
					self.on_take_card(card)
					return
				end
			end
		end
	end
end

function CardHolder:on_update(dt)
	-- for _, card in ipairs(self.cards) do
	--     local pos = go.get_position(card.url)
	--     print(pos)
	--     local w_pos_s = pos
	--     local w_pos_e = pos
	--     w_pos_s.x = pos.x - Card.width /2
	--     w_pos_e.x = pos.x + Card.width /2
	--     msg.post("@render:", "draw_line", {
	--         start_point = w_pos_s,
	--         end_point = w_pos_e,
	--         color = vmath.vector4(1, 1, 1, 1)
	--     })
	-- end
end

return CardHolder

local Card = {}

Card.Suit = {
	HEARTS = 0,
	SPADES = 1,
	CLUBS = 2,
	DIAMONDS = 3,
}

function Card:new(o)
	o = o or {}
	assert(o.rank)
	assert(o.suit)
	assert(o.url)
	setmetatable(o, self)
	self.__index = self
	return o
end

function Card:equal(other)
	if other.rank == self.ranks and other.soit == self.suit then
		return true
	end
	return false
end

function Card:show() end
function Card:hide() end

return Card

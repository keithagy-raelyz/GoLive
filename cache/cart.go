package cache

import (
	"GoLive/db"
)

type Cart struct {
	contents          []CartItem
	paymentProcessing bool // true = processing; false = active/user still shopping
}

type CartItem struct {
	Product db.Product
	Count   int
}

func NewCart() Cart {
	return Cart{make([]CartItem, 0), false}
}

func (c CartItem) Total(quantity int, price float64) float64 {
	return float64(quantity) * price
}
func (c CartItem) Value() float64 {
	return float64(c.Count) * c.Product.Price
}

func (c Cart) Contents() []CartItem {
	return c.contents
}

func (c Cart) GrandTotal() float64 {
	var total float64
	for _, cartItem := range c.contents {
		total += cartItem.Value()
	}
	return total
}

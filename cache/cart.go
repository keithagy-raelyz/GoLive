package cache

import (
	"GoLive/db"
)

type Cart []CartItem

type CartItem struct {
	Product db.Product
	Count   int
}

func (c CartItem) Total(quantity int, price float64) float64 {
	return float64(quantity) * price
}
func (c CartItem) Value() float64 {
	return float64(c.Count) * c.Product.Price
}

func (c Cart) GrandTotal() float64 {
	var total float64
	for _, cartItem := range c {
		total += cartItem.Value()
	}
	return total
}

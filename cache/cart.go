package cache

import (
	"GoLive/db"
)

//Cart type stores an array of CartItems for processing. The paymentProcessing boolean tracks the status of the cart.
type Cart struct {
	Contents          CartContents
	paymentProcessing bool // true = processing; false = active/user still shopping
}

type CartContents struct {
	Contents []CartItem
}

//CartItem stores a product and the amount added to the cart.
type CartItem struct {
	Product db.Product
	Count   int
}

//NewCart initializes an empty cart.
func NewCart() Cart {
	return Cart{CartContents{make([]CartItem, 0)}, false}
}

//Total calculates the total amount paid for a particular type of cartitem.
func (c CartItem) Total(quantity int, price float64) float64 {
	return float64(quantity) * price
}

//Value reflects the price of a cartItem.
func (c CartItem) Value() float64 {
	return float64(c.Count) * c.Product.Price
}

// //Contents provides a copy of the cart.
// func (c Cart) Contents() CartContents {
// 	return c.Contents
// }

//GrandTotal is used in templates to provide an indication of the total amount.
func (c CartContents) GrandTotal() float64 {
	var total float64
	for _, cartItem := range c.Contents {
		total += cartItem.Value()
	}
	return total
}

package cache

import (
	"GoLive/db"
	"time"
)

//ActiveSession interface contains methods to manipulate session Data.
type ActiveSession interface {
	// monitor checks the expiry time.
	monitor()
	// updateExpiryTime updates the expiry time when the cache data gets accessed.
	updateExpiryTime(time.Time)
	// getSessionID returns the session's underlying ID.
	getKey() string
	// GetSessionOwner returns the owner of a given session
	GetSessionOwner() (db.MerchantUser, []CartItem)
}

//activeUserCache is a wrapper for the default cache type for each Cache to be distinguishable and have internal methods if required.
type activeUserCache struct {
	cache
}

// UserSession implements the ActiveSession interface.
// Stores information about the logged in user and his cart data.
type UserSession struct { // cart CRUD tied to methods on this type.
	cart  Cart
	owner db.User // owner can be User or a Merchant.
	session
}

//GetSessionOwner returns the user and his cart
func (u *UserSession) GetSessionOwner() (db.MerchantUser, Cart) {
	return db.MerchantUser{
		User:      u.owner,
		MerchDesc: ""}, u.cart
}

// UpdateCart updates the session's cart.
func (u *UserSession) updateCart(productID string, operator string, product *db.Product, items *itemCache) {
	switch operator {
	case "+":
		for i := range u.cart.Contents.Contents {
			if u.cart.Contents.Contents[i].Product.Id == productID {
				u.cart.Contents.Contents[i].Count++
				items.blockStock(productID)
			}
		}
	case "-":
		for i := range u.cart.Contents.Contents {
			if u.cart.Contents.Contents[i].Product.Id == productID {
				u.cart.Contents.Contents[i].Count--
				items.blockStock(productID)
				if u.cart.Contents.Contents[i].Count == 0 {
					if i == 0 {
						u.cart.Contents = CartContents{u.cart.Contents.Contents[i+1:]}
					} else if i == len(u.cart.Contents.Contents)-1 {
						u.cart.Contents = CartContents{u.cart.Contents.Contents[0:i]}
					} else {
						firsthalf := u.cart.Contents.Contents[0:i]
						secondhalf := u.cart.Contents.Contents[i+1:]
						u.cart.Contents = CartContents{append(firsthalf, secondhalf...)}
					}
				}
			}
		}

	case "append":
		cartItem := CartItem{
			*product,
			1,
		}
		if alreadyInCart(u.cart.Contents.Contents, cartItem) {
			u.updateCart(productID, "+", product, items)
			return
		} else {
			u.cart.Contents = CartContents{append(u.cart.Contents.Contents, cartItem)}
		}
	}
}

func alreadyInCart(cartContents []CartItem, item CartItem) bool {
	for _, existing := range cartContents {
		if existing.Product.Id == item.Product.Id {
			return true
		}
	}
	return false
}

//clearCart empties the user cart.
func (u *UserSession) clearCart() {
	u.cart = Cart{}
}

// //merchantCache is a wrapper for the default cache type for each Cache to be distinguishable and have internal methods if required.
// type merchantCache struct {
// 	cache
// }

//MerchantSession implements the ActiveSession interface by embedding the session struct.
//stores information about the logged in merchant user.
type MerchantSession struct {
	owner db.MerchantUser
	session
}

//GetSessionOwner returns the merchant user
func (m *MerchantSession) GetSessionOwner() (db.MerchantUser, []CartItem) {
	return m.owner, nil
}

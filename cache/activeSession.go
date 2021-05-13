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
	// GetSessionOwner returns the owner (type db.User or db.MerchantUser) of a given session
	GetSessionOwner()
}

// UserSession implements the ActiveSession interface.
// Stores information about the logged in user and his cart data.
type UserSession struct { // cart CRUD tied to methods on this type.
	cart  Cart
	owner db.User // owner can be User or a Merchant.
	session
}

func (u *UserSession) GetSessionOwner() (db.User, []CartItem) {
	return u.owner, u.cart
}

// UpdateCart updates the session's cart.
func (u *UserSession) updateCart(productID string, operator string, product *db.Product) {
	switch operator {
	case "+":
		for i := range u.cart {
			if u.cart[i].Product.Id == productID {
				u.cart[i].Count++
				a.cacheManager.BlockStock(productID)
			}
		}
	case "-":
		for i := range u.cart {
			if u.cart[i].Product.Id == productID {
				u.cart[i].Count--
				a.cacheManager.BlockStock(productID)
				if u.cart[i].Count == 0 {
					if i == 0 {
						u.cart = u.cart[i+1:]
					} else if i == len(u.cart)-1 {
						u.cart = u.cart[0:i]
					} else {
						firsthalf := u.cart[0:i]
						secondhalf := u.cart[i+1:]
						u.cart = append(firsthalf, secondhalf...)
					}
				}
			}
		}

	case "append":
		cartItem := CartItem{
			*product,
			1,
		}
		u.cart = append(u.cart, cartItem)
	}
}

func (u *UserSession) clearCart() {
	u.cart = Cart{}
}

//MerchantSession implements the ActiveSession interface by embedding the session struct.
//stores information about the logged in merchant user.
type MerchantSession struct {
	owner db.MerchantUser
	session
}

func (m *MerchantSession) GetSessionOwner() db.MerchantUser {
	return m.owner
}

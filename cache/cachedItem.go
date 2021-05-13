package cache

import (
	"GoLive/db"
	"time"
)

//cachedItem stores the information of a product and its session.
type cachedItem struct {
	item   *db.Product
	expiry time.Time
}

func (c *cachedItem) getKey() string {
	return c.item.Id
}

func (c *cachedItem) addQty(amt int) {
	c.item.Quantity += amt
}

func (c *cachedItem) reduceQty(amt int) {
	c.item.Quantity -= amt
}

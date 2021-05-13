package cache

import (
	"GoLive/db"
	"time"
)

//itemCache is a wrapper for the default cache type for each Cache to be distinguishable and have internal methods if required.
type itemCache struct {
	cache
	// itemsCache map[string]cachedItem
	// sorted     []db.Product
}

func (i *itemCache) blockStock(productID string) {
	i.cache[productID].(*cachedItem).reduceQty(1)
	i.cache[productID].(*cachedItem).addSales(1)
}

func (i *itemCache) releaseStock(productID string) {
	i.cache[productID].(*cachedItem).addQty(1)
	i.cache[productID].(*cachedItem).reduceSales(1)
}

func (c *itemCache) rollback(user *UserSession) {
	cart := user.cart
	for _, item := range cart.contents {
		c.cache[item.Product.Id].(*cachedItem).addQty(item.Count)
		c.cache[item.Product.Id].(*cachedItem).reduceQty(item.Count)
		c.cache[item.Product.Id].updateExpiryTime(time.Now().Add(SessionLife * time.Minute))
	}
}

func (i *itemCache) increase(prodID string, amt int, database *db.Database) error {
	retrievedItem, err := i.get(prodID, "", database)
	if err != nil {
		return err
	}
	retrievedItem.(*cachedItem).addQty(amt)
	//update DB return DB err if any
	return database.UpdateProduct(*retrievedItem.(*cachedItem).item)
}

func (i *itemCache) decrease(prodID string, amt int, database *db.Database) error {
	retrievedItem, err := i.get(prodID, "", database)
	if err != nil {
		return err
	}
	retrievedItem.(*cachedItem).reduceQty(amt)
	//update DB return DB err if any
	return database.UpdateProduct(*retrievedItem.(*cachedItem).item)
}

// func (i *itemCache) getItem(prodID string, database *db.Database) (*cachedItem, error) {
// 	if itemCached, ok := ((*i).cache[prodID]); !ok {

// 		//Additem gets from the DB
// 		itemToCache, err := i.AddItemFromDB(prodID, database)
// 		if err != nil {
// 			return &itemToCache, err
// 		}

// 		//additemtocache
// 		(*i).cache[itemToCache.item.Id] = &itemToCache
// 		go i.tidy(itemToCache.item.Id, &itemToCache)
// 		return &itemToCache, nil
// 	} else {
// 		itemCached.updateExpiryTime(time.Now().Add(SessionLife * time.Minute))
// 		((*i).cache[prodID]) = itemCached
// 		return itemCached.(*cachedItem), nil
// 	}
// }

// func (i *itemCache) CheckItem(prodID string) (cachedItem, bool) {
// 	itemCached, ok := (*i).cache[prodID]
// 	return *(itemCached.(*cachedItem)), ok
// }

// func (i *itemCache) tidy(prodID string) {
// 	item := (*i).cache[prodID]
// 	item.monitor()
// 	delete((*i).cache, prodID)

// }

//cachedItem stores the information of a product and its session.
type cachedItem struct {
	item *db.Product
	expiry
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

func (c *cachedItem) addSales(amt int) {
	c.item.Sales += amt
}

func (c *cachedItem) reduceSales(amt int) {
	c.item.Sales -= amt
}

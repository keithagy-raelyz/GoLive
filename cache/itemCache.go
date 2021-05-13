package cache

import (
	"GoLive/db"
	"time"
)

//itemCache is a wrapper for the default cache type for each Cache to be distinguishable and have internal methods if required.
type itemCache struct {
	cache
	// itemsCache map[string]CachedItem
	// sorted     []db.Product
}

//blockStock blocks the stock out preventing others from checking out when added to a cart.
func (i *itemCache) blockStock(productID string) {
	i.cache[productID].(*CachedItem).reduceQty(1)
	i.cache[productID].(*CachedItem).addSales(1)
}

//releaseStock releases the stock if checkout fails.
func (i *itemCache) releaseStock(productID string) {
	i.cache[productID].(*CachedItem).addQty(1)
	i.cache[productID].(*CachedItem).reduceSales(1)
}

func (c *itemCache) rollback(user *UserSession) {
	cart := user.cart
	for _, item := range cart.contents {
		c.cache[item.Product.Id].(*CachedItem).addQty(item.Count)
		c.cache[item.Product.Id].(*CachedItem).reduceQty(item.Count)
		c.cache[item.Product.Id].updateExpiryTime(time.Now().Add(SessionLife * time.Minute))
	}
}

//increase increases the quantity of an item in the event the item is updated.
func (i *itemCache) increase(prodID string, amt int, database *db.Database) error {
	retrievedItem, err := i.get(prodID, "", database)
	if err != nil {
		return err
	}
	retrievedItem.(*CachedItem).addQty(amt)
	//update DB return DB err if any
	return database.UpdateProduct(*retrievedItem.(*CachedItem).item)
}

//decrease decreases the quantity of an item in the event the item is updated.
func (i *itemCache) decrease(prodID string, amt int, database *db.Database) error {
	retrievedItem, err := i.get(prodID, "", database)
	if err != nil {
		return err
	}
	retrievedItem.(*CachedItem).reduceQty(amt)
	//update DB return DB err if any
	return database.UpdateProduct(*retrievedItem.(*CachedItem).item)
}

//getAllProducts returns all the products stored in the cache. If the qty is less than 5, ping the DB to fetch a full row of products.
func (i *itemCache) getAllProducts(database *db.Database) ([]*CachedItem, error) {
	var count int
	var allProducts []*CachedItem
	for _, v := range (*i).cache {
		allProducts = append(allProducts, v.(*CachedItem))
		count++
	}
	if count < 5 {
		var allProductsFromDB []*CachedItem
		allProds, err := database.GetAllProducts()
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(allProds); j++ {
			itemToCache := &CachedItem{
				&allProds[j], expiry{time.Now().Add(SessionLife * time.Minute)},
			}
			(*i).cache[allProds[j].Id] = itemToCache
			allProductsFromDB = append(allProductsFromDB, itemToCache)
		}
		return allProductsFromDB, nil
	} else {
		return allProducts, nil
	}
}

//updateProduct updates the product if its found in the cache and updates the copy stored in the DB.
func (i *itemCache) updateProduct(product db.Product, database *db.Database) error {
	if cachedItem, ok := i.cache[product.Id]; !ok {
		itemToCache := &CachedItem{
			item:   &product,
			expiry: expiry{},
		}
		i.cache[product.Id] = itemToCache
		return database.UpdateProduct(product)
	} else {
		cachedItem.(*CachedItem).update(product)
		return database.UpdateProduct(product)
	}
}

// func (i *itemCache) getItem(prodID string, database *db.Database) (*CachedItem, error) {
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
// 		return itemCached.(*CachedItem), nil
// 	}
// }

// func (i *itemCache) CheckItem(prodID string) (CachedItem, bool) {
// 	itemCached, ok := (*i).cache[prodID]
// 	return *(itemCached.(*CachedItem)), ok
// }

// func (i *itemCache) tidy(prodID string) {
// 	item := (*i).cache[prodID]
// 	item.monitor()
// 	delete((*i).cache, prodID)

// }

//CachedItem stores the information of a product and its session.
type CachedItem struct {
	item *db.Product
	expiry
}

//getKey returns the item's ID.
func (c *CachedItem) getKey() string {
	return c.item.Id
}

//addQty increases the quantity of the item.
func (c *CachedItem) addQty(amt int) {
	c.item.Quantity += amt
}

//reduceQty reduces the quantity of the item.
func (c *CachedItem) reduceQty(amt int) {
	c.item.Quantity -= amt
}

//addSales adds a sale to the item.
func (c *CachedItem) addSales(amt int) {
	c.item.Sales += amt
}

//reduceSales reduce the sales of the item.
func (c *CachedItem) reduceSales(amt int) {
	c.item.Sales -= amt
}

//Copy returns a copy of the stored product.
func (c *CachedItem) Copy() db.Product {
	return *c.item
}

//update updates the stored product with the new values.
func (c *CachedItem) update(product db.Product) {
	c.item.Price = product.Price
	c.item.Name = product.Name
	c.item.Quantity = product.Quantity
	c.item.Thumbnail = product.Thumbnail
	c.item.Price = product.Price
	c.item.ProdDesc = product.ProdDesc
}

//deleteProduct deletes the product found in the cache and in the database.
func (i *itemCache) deleteProduct(prodID string, merchID string, database *db.Database) error {
	delete(i.cache, prodID)
	return database.DeleteProduct(prodID, merchID)
}

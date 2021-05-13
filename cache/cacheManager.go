package cache

import (
	"GoLive/db"
	"database/sql"
	"errors"
)

//CacheManager stores the various cache types and has CRUD functionality to access the underlying cache methods.
type CacheManager struct {
	activeUserCache activeUserCache // key: session value; value: logged in users and merchants.
	itemCache       itemCache       // Populates home page.
	database        *db.Database    // Pointer to DB
}

// NewCacheManager returns a pointer to an initialized CacheManager.
func NewCacheManager(database *sql.DB) *CacheManager {
	db := &db.Database{}
	db.InitializeDB(database)

	newAuC := activeUserCache{make(map[string]CacheObject)}
	newIC := itemCache{make(map[string]CacheObject)}
	newCM := &CacheManager{
		activeUserCache: newAuC,
		itemCache:       newIC,
		database:        db,
	}

	return newCM
}

func (c *CacheManager) GetAllProducts() ([]db.Product, error) {
	return c.database.GetAllProducts()
}

func (c *CacheManager) GetAllMerchants() ([]db.Merchant, error) {
	return c.database.GetAllMerchants()
}

//AddtoCache identifies the type of the payload before adding it into the respective cache by calling on the respective cache.add() method.
func (c *CacheManager) AddtoCache(payLoad CacheObject) {
	switch v := payLoad.(type) {
	case ActiveSession:
		c.activeUserCache.add(v)
	case *cachedItem:
		c.itemCache.add(v)
	}
}

//UpdateCache identifies the type of the payload before adding it into the respective cache by calling on the respective cache.updateExpiryTime() method.
// In the case of a user updating their cart, UpdateCache also updates the respective product cart.
func (c *CacheManager) UpdateCache(payLoad CacheObject, key string, productID string, operator string) {
	switch v := payLoad.(type) {
	case *UserSession:
		retrievedItem, _ := c.itemCache.get(key, productID, c.database)
		c.activeUserCache.update(v.getKey(), productID, operator, *retrievedItem.(*cachedItem), &c.itemCache)
	case *MerchantSession:
		c.activeUserCache.update(v.getKey(), "", "", cachedItem{}, &c.itemCache)
	case *cachedItem:
		c.activeUserCache.update(v.getKey(), "", "", cachedItem{}, &c.itemCache)
	}
}

//GetFromCache identifies the type of the payload before retrieving it from the respective cache by calling on the respective cache.get() method.
func (c *CacheManager) GetFromCache(key string, userID string, cacheType string) (CacheObject, error) {
	switch cacheType {
	case "activeSession":
		return c.activeUserCache.get(key, userID, c.database)
	case "cachedItems":
		return c.itemCache.get(key, "", c.database)
	}
	return nil, errors.New("invalid cacheType")
}

//RemoveFromCache identifies the type of hte payload before removing it from the respective cache by calling on the respective cache.remove() method.
func (c *CacheManager) RemoveFromCache(key string, cacheType string) {
	switch cacheType {
	case "activeSession":
		c.activeUserCache.remove(key)
	case "cachedItems":
		c.itemCache.remove(key)
	}
}

func (c *CacheManager) ReleaseStock(productID string) {
	c.itemCache.releaseStock(productID)
}

func (c *CacheManager) AddCartProcessing(userID string) {
	c.activeUserCache.cache[userID].(*UserSession).cart.paymentProcessing = true
}

func (c *CacheManager) CartSuccess(userID string) {
	delete(c.activeUserCache.cache, userID)
}

func (c *CacheManager) CartFailure(userID string) {
	c.itemCache.rollback(c.activeUserCache.cache[userID].(*UserSession))
}

func (c *CacheManager) ClearActiveUserCart(sessionID string) {
	c.activeUserCache.cache[sessionID].(*UserSession).clearCart()
}

func (c *CacheManager) UpdateItemInCache(prodID string, operator string, amt int) error {
	switch operator {
	case "+":
		return c.itemCache.increase(prodID, amt, c.database)
	case "-":
		return c.itemCache.decrease(prodID, amt, c.database)
	default:
		return errors.New("invalid operator supplied")
	}
}

// func (c *CacheManager) GetItemFromCache(prodID string) (*cachedItem, error) {
// 	return c.itemCache.get(prodID, c.database)
// }

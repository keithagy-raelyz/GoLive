package cache

import (
	"GoLive/db"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

//CacheManager stores the various cache types and has CRUD functionality to access the underlying cache methods.
type CacheManager struct {
	activeUserCache activeUserCache // key: session value; value: logged in users and merchants.
	merchantCache   merchantCache   // cached Merchant pages
	itemCache       itemCache       // Populates home page.
	database        *db.Database    // Pointer to DB
}

// NewCacheManager returns a pointer to an initialized CacheManager.
func NewCacheManager(database *sql.DB) *CacheManager {
	db := &db.Database{}
	db.InitializeDB(database)

	newAuC := activeUserCache{make(map[string]CacheObject)}
	newIC := itemCache{make(map[string]CacheObject)}
	newMc := merchantCache{make(map[string]CacheObject)}
	newCM := &CacheManager{
		activeUserCache: newAuC,
		itemCache:       newIC,
		merchantCache:   newMc,
		database:        db,
	}

	return newCM
}

//GetAllProducts returns all the products from the cache.
func (c *CacheManager) GetAllProducts() ([]db.Product, error) {
	allProds, err := c.itemCache.getAllProducts(c.database)
	if err != nil {
		return nil, err
	}

	var products []db.Product
	for _, prod := range allProds {
		products = append(products, prod.Copy())

	}
	return products, nil
}

//.GetAllMerchants calls on cache method to obtain the cached merchant objects before obtaining the stored data and returning
func (c *CacheManager) GetAllMerchants() ([]db.Merchant, error) {
	allMerch, err := c.merchantCache.getAllMerchants(c.database)
	if err != nil {
		return nil, err
	}

	var merchants []db.Merchant
	for _, merch := range allMerch {
		merchants = append(merchants, *merch.merchantPage)

	}
	return merchants, nil
}

//AddtoCache identifies the type of the payload before adding it into the respective cache by calling on the respective cache.add() method.
func (c *CacheManager) AddtoCache(payLoad CacheObject) {
	switch v := payLoad.(type) {
	case *UserSession:
		c.activeUserCache.add(v)
	case *MerchantSession:
		c.activeUserCache.add(v)
	case *CachedMerchant:
		c.merchantCache.add(v)
	case *CachedItem:
		c.itemCache.add(v)
	}
}

//UpdateCache identifies the type of the payload before adding it into the respective cache by calling on the respective cache.updateExpiryTime() method.
// In the case of a user updating their cart, UpdateCache also updates the respective product cart.
func (c *CacheManager) UpdateCache(payLoad CacheObject, productID string, operator string) {
	switch v := payLoad.(type) {
	case *UserSession:
		retrievedItem, _ := c.itemCache.get(productID, "product", c.database)
		c.activeUserCache.update(v.getKey(), productID, operator, *retrievedItem.(*CachedItem), &c.itemCache)
	case *MerchantSession:
		c.activeUserCache.update(v.getKey(), "", "", CachedItem{}, &c.itemCache)
	case *CachedItem:
		c.activeUserCache.update(v.getKey(), "", "", CachedItem{}, &c.itemCache)
	}
}

//GetFromCache identifies the type of the payload before retrieving it from the respective cache by calling on the respective cache.get() method.
func (c *CacheManager) GetFromCache(key string, cacheType string) (CacheObject, error) {
	switch cacheType {
	case "activeSession":
		return c.activeUserCache.get(key, "", c.database)
	case "cachedMerchants":
		return c.merchantCache.get(key, "merchant", c.database)
	case "cachedItems":
		return c.itemCache.get(key, "product", c.database)
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

//ReleaseStock releases the stock for processing
func (c *CacheManager) ReleaseStock(productID string) {
	c.itemCache.releaseStock(productID)
}

//AddCartProcessing triggers the boolean to identify that the cart is being processed.
func (c *CacheManager) AddCartProcessing(userID string) {
	c.activeUserCache.cache[userID].(*UserSession).cart.paymentProcessing = true
}

//CartFailure triggers a rollback on the inventory
func (c *CacheManager) CartFailure(userID string) {
	c.itemCache.rollback(c.activeUserCache.cache[userID].(*UserSession))
}

//ClearActiveUserCart clears the cart of the user after checkout
func (c *CacheManager) ClearActiveUserCart(sessionID string) {
	c.activeUserCache.cache[sessionID].(*UserSession).clearCart()
}

//UpdateItemInCache updates the item count in the cache if someone adds it into his cart.
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

// func (c *CacheManager) GetItemFromCache(prodID string) (*CachedItem, error) {
// 	return c.itemCache.get(prodID, c.database)
// }

//UpdateProduct updates the product present in the cache.
func (c *CacheManager) UpdateProduct(product db.Product) error {
	return c.itemCache.updateProduct(product, c.database)
}

//DeleteProduct deletes the product present in the cache and the database.
func (c *CacheManager) DeleteProduct(prodID string, merchID string) error {
	c.merchantCache.deleteProduct(prodID, merchID, c.database)
	return c.itemCache.deleteProduct(prodID, merchID, c.database)
}

//CreateMerchant calls on the createmerchant method from the database.
func (c *CacheManager) CreateMerchant(merchant db.MerchantUser, pw string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	if err != nil {
		return err
	}
	return c.database.CreateMerchant(merchant, string(hashed))
}

//CheckMerchant checks if the merchant exists in the database.
func (c *CacheManager) CheckMerchant(merchant db.MerchantUser) error {
	return c.database.CheckMerchant(merchant)
}

//CheckUser checks if the user exists in the database.
func (c *CacheManager) CheckUser(user db.User) error {
	return c.database.CheckUser(user)
}

//CreateUser creates a user in the database.
func (c *CacheManager) CreateUser(user db.User, pw string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	if err != nil {
		return err
	}
	return c.database.CreateUser(user, string(hashed))
}

//UpdateUser updates the user in the database.
func (c *CacheManager) UpdateUser(user db.User) error {
	return c.database.UpdateUser(user)
}

//DeleteUser deletes the User from the database.
func (c *CacheManager) DeleteUser(userId string) error {
	return c.database.DeleteUser(userId)
}

//UpdateMerchantCache updates the merchant stored in the cache if found and the database.
func (c *CacheManager) UpdateMerchantCache(merchant db.MerchantUser) error {
	return c.merchantCache.updateMerchant(merchant, c.database)
}

//DeleteMerchantFromCache removes the merchant from the cache and the database.
func (c *CacheManager) DeleteMerchantFromCache(merchID string) error {
	return c.merchantCache.deleteMerchant(merchID, c.database)
}

//CreateProduct creates a new products in the database.
func (c *CacheManager) CreateProduct(product db.Product) error {
	return c.database.CreateProduct(product)
}

//GetUser gets the user profile from the database.
func (c *CacheManager) GetUser(userID string) (db.User, error) {
	return c.database.GetUserFromID(userID)
}

//GetUsers is an admin function to view all users present in the database.
func (c *CacheManager) GetUsers() ([]db.User, error) {
	return c.database.GetUsers()
}

//UserLogin gets the user credentials after the pw has been authenticated.
func (c *CacheManager) UserLogin(username string) (db.User, error) {
	return c.database.GetUser(username)
}

//MerchantLogin gets the user credentials after the pw has been authenticated.
func (c *CacheManager) MerchantLogin(username string) (db.MerchantUser, error) {
	return c.database.GetMerchant(username)
}

package cache

import (
	"GoLive/db"
	"errors"
	"fmt"
	"sort"
	"time"
)

const (
	SessionLife = 30 // default life of active session in minutes.
)

//CacheManager stores the various cache types and has CRUD functionality to access the underlying cache methods.
type CacheManager struct {
	activeUserCache activeUserCache // key: session value; value: see session definition below.
	merchantCache   merchantCache   // key: Pointer to cache; value: access frequency.
	itemsCache      itemCache       // Populates home page.
	database        *db.Database
}

// NewCacheManager returns a pointer to an initialized CacheManager.
func NewCacheManager(database *db.Database) *CacheManager {
	newAuC := activeUserCache{make(map[string]ActiveSession)}
	newMC := merchantCache{make(map[string]ActiveSession)}
	newIC := itemCache{make(map[string]cachedItem), make([]db.Product, 0)}

	newCM := &CacheManager{
		activeUserCache: newAuC,
		merchantCache:   newMC,
		itemsCache:      newIC,
		database:        database,
	}

	return newCM
}

// NewUserSession takes in required inputs and returns a new UserSession.
func NewUserSession(sessionKey string, expiry time.Time, user db.User, cart Cart) *UserSession {
	return &UserSession{
		session: session{sessionKey, expiry},
		owner:   user,
		cart:    cart,
	}
}

// NewMerchantSession takes in required inputs and returns a new MerchantSession
func NewMerchSession(sessionKey string, expiry time.Time, user db.MerchantUser) *MerchantSession {
	return &MerchantSession{
		session: session{sessionKey, expiry},
		owner:   user,
	}
}

//AddtoCache identifies the type of the payload before adding it into the respective cache by calling on the respective cache.add() method.
func (c *CacheManager) AddtoCache(payLoad ActiveSession) {
	switch v := payLoad.(type) {
	case *UserSession:
		c.activeUserCache.add(v)
	case *MerchantSession:
		c.activeUserCache.add(v)
	case *cachedMerchant:
		c.merchantCache.add(v)

	}
}

//UpdateCache identifies the type of the payload before adding it into the respective cache by calling on the respective cache.updateExpiryTime() method.
// In the case of a user updating their cart, UpdateCache also updates the respective product cart.
func (c *CacheManager) UpdateCart(payLoad ActiveSession, productID string, operator string) {
	switch v := payLoad.(type) {
	case *UserSession:
		cachedItem, _ := c.itemsCache.get(productID, c.database)
		c.activeUserCache.update(v.getSessionID(), productID, operator, cachedItem)
	case *MerchantSession:
		c.activeUserCache.update(v.getSessionID(), "", "", cachedItem{})
	case *cachedMerchant:
		c.activeUserCache.update(v.getSessionID(), "", "", cachedItem{})
	}
}

type SessionCopy interface {
	Data()
}

//GetFromCache identifies the type of the payload before retrieving it from the respective cache by calling on the respective cache.get() method.
func (c *CacheManager) GetFromCache(key string, cacheType string) (ActiveSession, bool) {
	switch cacheType {
	case "activeUsers":
		return c.activeUserCache.get(key)
	case "activeMerchants":
		return c.activeUserCache.get(key)
	case "cachedMerchants":
		return c.merchantCache.get(key)

	}
	return nil, false
}

//RemoveFromCache identifies the type of hte payload before removing it from the respective cache by calling on the respective cache.remove() method.
func (c *CacheManager) RemoveFromCache(key string, cacheType string) {
	switch cacheType {
	case "activeUsers":
		c.activeUserCache.remove(key)
	case "activeMerchants":
		c.activeUserCache.remove(key)
	case "cachedMerchants":
		c.merchantCache.remove(key)

	}

}

//UpdateCacheFromDB feeds an array of the payload which is obtained from the latest DB query and calls on the respective cache's cache.update() method.
func (c *CacheManager) UpdateCacheFromDB(payLoad ...ActiveSession) {
	switch v := payLoad[0].(type) {
	case *UserSession:
		c.activeUserCache.refresh(v)
	case *MerchantSession:
		c.activeUserCache.refresh(v)
	case *cachedMerchant:
		c.merchantCache.refresh(v)

	}
}

//ActiveSession interface contains methods to manipulate session Data.
type ActiveSession interface {
	//monitor checks the expiry time.
	monitor()
	//getSessionID returns the session's underlying ID.
	getSessionID() string
	//updateExpiryTime updates the expiry time when the cache data gets accessed.
	updateExpiryTime(time.Time)
}

//session implements the ActiveSession interface.
//similarly embedding in any struct type allows the struct to implement the ActiveSession interface.
type session struct {
	key    string
	expiry time.Time
}

//monitor sleeps for the difference between the expiration time and the current time.
//eg expires at 4pm, current time is 4.30pm. it sleeps for 30mins.
//after sleep ends, it checks if the expiration time has exceeded the current time it returns.
//or else it loops.
func (s *session) monitor() {
	for {
		sleeptime := time.Until(s.expiry)
		time.Sleep(sleeptime)
		if s.expiry.Before(time.Now()) {
			return
		}
	}
}

//getSessionID returns the stored key.
func (s *session) getSessionID() string {
	return s.key
}

//updateExpiryTime updates the session's expiry time.
func (s *session) updateExpiryTime(updatedTime time.Time) {
	s.expiry = updatedTime
}

//MerchantSession implements the ActiveSession interface by embedding the session struct.
//stores information about the logged in merchant user.
type MerchantSession struct {
	session
	owner db.MerchantUser
}

func (m *MerchantSession) GetSessionOwner() db.MerchantUser {
	return m.owner
}

//UserSession implements the ActiveSession interface by embedding the session struct.
//Stores information about the logged in user and his cart data.
type UserSession struct { // cart CRUD tied to methods on this type.
	session
	owner db.User // owner can be User or a Merchant.
	cart  Cart
}

func (u *UserSession) GetSessionOwner() (db.User, []CartItem) {
	return u.owner, u.cart
}

// UpdateCart updates the session's cart.
func (u *UserSession) updateCart(productID string, operator string, product *db.Product) {
	switch operator {
	case "+":
		for i, _ := range u.cart {
			if u.cart[i].Product.Id == productID {
				u.cart[i].Count++
			}
		}
	case "-":
		for i, _ := range u.cart {
			if u.cart[i].Product.Id == productID {
				u.cart[i].Count--
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

//activeUserCache is a wrapper for the default cache type for each Cache to be distinguishable and have internal methods if required.
type activeUserCache struct {
	cache
}

//cache type is a map that maps the key to an active session
type cache map[string]ActiveSession // key: session value; value: see session definition below.

//add takes in an active session, checks if the active session already exists in the cache.
//if it doesn't exist in the cache, add it into the cache and fire off a go tidy() go routine which regulates the session and removes it on expiry.
func (c *cache) add(payLoad ActiveSession) {
	key := payLoad.getSessionID()
	if ok := c.check(key); !ok {
		(*c)[key] = payLoad
		fmt.Println((*c)[key], "added into the cache")
		fmt.Println(key, "key value during storage")
		go c.tidy(key, payLoad)
	} else {

	}
}

func (c *cache) update(key string, productID string, operator string, cachedItem cachedItem) {
	if c.check(key) {
		activeSession, _ := (*c)[key]
		switch v := activeSession.(type) {
		case *UserSession:
			v.updateCart(productID, operator, cachedItem.item)
		case *MerchantSession:
			//c.activeUserCache.refresh(v)
		case *cachedMerchant:
			//c.merchantCache.refresh(v)

		}

	}

}

//refresh flushes the cache with the data obtained from the DB.
func (c *cache) refresh(payLoad ...ActiveSession) {
	flushedMap := make(map[string]ActiveSession)
	for _, activeSession := range payLoad {
		flushedMap[activeSession.getSessionID()] = activeSession
	}
	(*c) = flushedMap
}

//get returns the ActiveSession stored in the cache given the key.
func (c *cache) get(key string) (ActiveSession, bool) {
	activeSession, ok := (*c)[key]
	fmt.Println((*c)[key], "getting from cache")
	fmt.Println("key value during access", key)
	return activeSession, ok
}

//tidy calls on monitor which checks if the session has expired. If the session has expired, monitor returns and the session is deleted from the cache.
func (c *cache) tidy(key string, session ActiveSession) {
	session.monitor()
	delete(*c, key)
}

// Validate that session is active, and if it is active refresh the timestamp.
func (c *cache) check(key string) bool {
	activeSession, ok := (*c)[key]
	if ok {
		activeSession.updateExpiryTime(time.Now().Add(SessionLife * time.Minute))
	}
	return ok
}

//remove deletes the  session given the ActiveSession.
func (c *cache) remove(key string) {
	delete(*c, key)
}

type Cart []CartItem

//cachedMerchant stores the merchant details and products.
type cachedMerchant struct {
	session
	merchant db.MerchantUser
	products *[]db.Product
	hitRate  int
}

//merchantCache is a wrapper for the default cache type for each Cache to be distinguishable and have internal methods if required.
type merchantCache struct {
	cache
}

////UpdateSorted flushes the sorted array with the new data from the DB. It gets called when the cache is flushed.
//func (i *itemCache) updateSorted() {
//	preSort := make([]db.Product, 30)
//	for _, as := range i.cache {
//		preSort = append(preSort, (as.(*cachedItem)).item)
//	}
//	i.sort()
//}

////update override the default update method for caches.
//func (i *itemCache) update(payLoad ...ActiveSession) {
//	flushedMap := make(map[string]ActiveSession)
//	for _, activeSession := range payLoad {
//		flushedMap[activeSession.getSessionID()] = activeSession
//	}
//	i.cache = flushedMap
//	i.updateSorted()
//}

//sort sorts the sorted item cache.
func (i *itemCache) sort() {
	sort.Slice(i.sorted, func(j, k int) bool { return i.sorted[j].Sales > i.sorted[k].Sales })
}

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

func rollback() {
	//db.exec(add back 1)
	//mutex.unlock()
	//count ++
}

//cachedItem stores the information of a product and its session.
type cachedItem struct {
	item   *db.Product
	expiry time.Time
}

func (c *cachedItem) updateExpiryTime(updatedTime time.Time) {
	c.expiry = updatedTime
}

func (c *cachedItem) monitor() {
	for {
		sleeptime := time.Until(c.expiry)
		time.Sleep(sleeptime)
		if c.expiry.Before(time.Now()) {
			return
		}
	}
}
func (c *CacheManager) UpdateItemInCache(prodID string, operator string, amt int) error {
	switch operator {
	case "+":
		return (*c).itemsCache.increase(prodID, amt, c.database)
	case "-":
		return (*c).itemsCache.decrease(prodID, amt, c.database)
	default:
		return errors.New("invalid operator supplied")
	}
}

func (i *itemCache) increase(prodID string, amt int, database *db.Database) error {
	cachedItem, err := i.get(prodID, database)
	if err != nil {
		return err
	}
	cachedItem.addQty(amt)
	//update DB return DB err if any
	return database.UpdateProduct(*cachedItem.item)
}

func (i *itemCache) decrease(prodID string, amt int, database *db.Database) error {
	cachedItem, err := i.get(prodID, database)
	if err != nil {
		return err
	}
	cachedItem.reduceQty(amt)
	//update DB return DB err if any
	return database.UpdateProduct(*cachedItem.item)
}

//func (c *CacheManager) CheckItemPresence(prodID string) (cachedItem, bool) {
//	cachedItem, ok := (*c).itemsCache.CheckItem(prodID)
//	if !ok {
//		//query db
//
//		(*c).itemsCache.AddItemFromDB(cachedItem)
//	}
//	return cachedItem, ok
//}

func (c *CacheManager) GetItemFromCache(prodID string) (cachedItem, error) {
	return (*c).itemsCache.get(prodID, c.database)
}

func (i *itemCache) get(prodID string, database *db.Database) (cachedItem, error) {
	if cachedItem, ok := ((*i).itemsCache[prodID]); !ok {

		//Additem gets fromt he DB
		cachedItem, err := i.AddItemFromDB(prodID, database)
		if err != nil {
			return cachedItem, err
		}

		//additemtocache
		(*i).itemsCache[cachedItem.item.Id] = cachedItem
		go i.tidy(cachedItem.item.Id)
		return cachedItem, nil
	} else {
		cachedItem.updateExpiryTime(time.Now().Add(SessionLife * time.Minute))
		((*i).itemsCache[prodID]) = cachedItem
		return cachedItem, nil
	}
}
func (i *itemCache) CheckItem(prodID string) (cachedItem, bool) {
	cachedItem, ok := (*i).itemsCache[prodID]
	return cachedItem, ok
}

func (i *itemCache) AddItemFromDB(prodID string, db *db.Database) (cachedItem, error) {
	//DB interaction with error returned return error here
	prod, err := db.GetProduct(prodID)
	if err != nil {
		return cachedItem{}, err
	}
	cachedItem := cachedItem{&prod, time.Now().Add(SessionLife * time.Minute)}
	return cachedItem, nil
}

func (i *itemCache) tidy(prodID string) {
	item := (*i).itemsCache[prodID]
	item.monitor()
	delete((*i).itemsCache, prodID)

}

//itemCache is a wrapper for the default cache type for each Cache to be distinguishable and have internal methods if required.
type itemCache struct {
	itemsCache map[string]cachedItem
	sorted     []db.Product
}

func (c *cachedItem) getID() string {
	return c.item.Id
}

func (c *cachedItem) addQty(amt int) {
	c.item.Quantity += amt
}

func (c *cachedItem) reduceQty(amt int) {
	c.item.Quantity -= amt
}

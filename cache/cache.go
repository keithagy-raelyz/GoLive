package cache

import (
	"GoLive/db"
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
}

// NewCacheManager returns a pointer to an initialized CacheManager.
func NewCacheManager() *CacheManager {
	newAuC := activeUserCache{make(map[string]ActiveSession)}
	newMC := merchantCache{make(map[string]ActiveSession)}
	newIC := itemCache{make(map[string]ActiveSession), make([]db.Product, 0)}

	newCM := &CacheManager{
		activeUserCache: newAuC,
		merchantCache:   newMC,
		itemsCache:      newIC,
	}

	return newCM
}

// NewUserSession takes in required inputs and returns a new UserSession.
func NewUserSession(sessionKey string, expiry time.Time, user db.User, cart *[]db.Product) *UserSession {
	return &UserSession{
		session: session{sessionKey, expiry},
		owner:   user,
		cart:    cart,
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
	case *cachedItem:
		c.itemsCache.add(v)
	}
}

//UpdateCache identifies the type of the payload before adding it into the respective cache by calling on the respective cache.updateExpiryTime() method.
// In the case of a user updating their cart, UpdateCache also updates the respective product cart.
func (c *CacheManager) UpdateCache(activeSession ActiveSession, cart *db.Product) {
	switch v := activeSession.(type) {
	case *UserSession:
		sessionID := activeSession.getSessionID()
		c.activeUserCache.
	case *MerchantSession:
		c.activeUserCache.add(v)
	case *cachedMerchant:
		c.merchantCache.add(v)
	case *cachedItem:
		c.itemsCache.add(v)
	}
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
	case "cachedItem":
		return c.itemsCache.get(key)
	}
	return nil, false
}

//RemoveFromCache identifies the type of hte payload before removing it from the respective cache by calling on the respective cache.remove() method.
func (c *CacheManager) RemoveFromCache(payLoad ActiveSession) {
	switch v := payLoad.(type) {
	case *UserSession:
		c.activeUserCache.remove(v)
	case *MerchantSession:
		c.activeUserCache.remove(v)
	case *cachedMerchant:
		c.merchantCache.remove(v)
	case *cachedItem:
		c.itemsCache.remove(v)
	}
}

//UpdateCacheFromDB feeds an array of the payload which is obtained from the latest DB query and calls on the respective cache's cache.update() method.
func (c *CacheManager) UpdateCacheFromDB(payLoad ...ActiveSession) {
	switch v := payLoad[0].(type) {
	case *UserSession:
		c.activeUserCache.update(v)
	case *MerchantSession:
		c.activeUserCache.update(v)
	case *cachedMerchant:
		c.merchantCache.remove(v)
	case *cachedItem:
		c.itemsCache.remove(v)
	}
}

//ActiveSession interface contains methods to manipulate session Data.
type ActiveSession interface {
	//monitor checks the expiry time.
	monitor()
	//getSessionID returns the session's underlying ID.
	getSessionID() string
	//UpdateExpiryTime updates the expiry time when the cache data gets accessed.
	UpdateExpiryTime(time.Time)
}

//session implements the ActiveSession interface.
//similarly embedding in any struct type allows the struct to implement the ActiveSession interface.
type session struct {
	key    string
	expiry time.Time
}

//MerchantSession implements the ActiveSession interface by embedding the session struct.
//stores information about the logged in merchant user.
type MerchantSession struct {
	session
	owner db.MerchantUser
}

//monitor sleeps for the difference between the expiration time and the current time.
//eg expires at 4pm, current time is 4.30pm. it sleeps for 30mins.
//after sleep ends, it checks if the expiration time has exceeded the current time it returns.
//or else it loops.
func (m *MerchantSession) monitor() {
	for {
		sleeptime := time.Until(m.session.expiry)
		time.Sleep(sleeptime)
		if m.session.expiry.Before(time.Now()) {
			return
		}
	}
}

//getSessionID returns the stored key.
func (m *MerchantSession) getSessionID() string {
	return m.session.key
}

//UpdateExpiryTime updates the session's expiry time.
func (m *MerchantSession) UpdateExpiryTime(updatedTime time.Time) {
	m.session.expiry = updatedTime
}

//UserSession implements the ActiveSession interface by embedding the session struct.
//Stores information about the logged in user and his cart data.
type UserSession struct { // Cart CRUD tied to methods on this type.
	session
	owner db.User // Owner can be User or a Merchant.
	cart  *[]db.Product
}

//monitor sleeps for the difference between the expiration time and the current time.
//eg expires at 4pm, current time is 4.30pm. it sleeps for 30mins.
//after sleep ends, it checks if the expiration time has exceeded the current time it returns.
//or else it loops.
func (u *UserSession) monitor() {
	for {
		sleeptime := time.Until(u.session.expiry)
		time.Sleep(sleeptime)
		if u.session.expiry.Before(time.Now()) {
			return
		}
	}
}

//getSessionID returns the stored key.
func (u *UserSession) getSessionID() string {
	return u.session.key
}

//UpdateExpiryTime updates the session's expiry time.
func (u *UserSession) UpdateExpiryTime(updatedTime time.Time) {
	u.session.expiry = updatedTime
}

// UpdateCart updates the session's cart.
func (u *UserSession) UpdateCart(cart *[]db.Product) {
	u.cart = cart
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
		go c.tidy(key, payLoad)
	}
}

//update flushes the cache with the data obtained from the DB.
func (c *cache) update(payLoad ...ActiveSession) {
	flushedMap := make(map[string]ActiveSession)
	for _, activeSession := range payLoad {
		flushedMap[activeSession.getSessionID()] = activeSession
	}
	(*c) = flushedMap
}

//get returns the ActiveSession stored in the cache given the key.
func (c *cache) get(key string) (ActiveSession, bool) {
	activeSession, ok := (*c)[key]
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
		activeSession.UpdateExpiryTime(time.Now().Add(SessionLife * time.Minute))
	}
	return ok
}

//remove deletes the  session given the ActiveSession.
func (c *cache) remove(session ActiveSession) {
	delete(*c, session.getSessionID())
}

//cachedMerchant stores the merchant details and products.
type cachedMerchant struct {
	session
	merchant db.MerchantUser
	products *[]db.Product
	hitRate  int
}

//monitor sleeps for the difference between the expiration time and the current time.
//eg expires at 4pm, current time is 4.30pm. it sleeps for 30mins.
//after sleep ends, it checks if the expiration time has exceeded the current time it returns.
//or else it loops.
func (c *cachedMerchant) monitor() {
	for {
		sleeptime := time.Until(c.session.expiry)
		time.Sleep(sleeptime)
		if c.session.expiry.Before(time.Now()) {
			return
		}
	}
}

//getSessionID returns the stored key.
func (c *cachedMerchant) getSessionID() string {
	return c.session.key
}

//UpdateExpiryTime updates the session's expiry time.
func (c *cachedMerchant) UpdateExpiryTime(updatedTime time.Time) {
	c.session.expiry = updatedTime
}

//merchantCache is a wrapper for the default cache type for each Cache to be distinguishable and have internal methods if required.
type merchantCache struct {
	cache
}

//cachedItem stores the information of a product and its session.
type cachedItem struct {
	session
	item db.Product
}

//monitor sleeps for the difference between the expiration time and the current time.
//eg expires at 4pm, current time is 4.30pm. it sleeps for 30mins.
//after sleep ends, it checks if the expiration time has exceeded the current time it returns.
//or else it loops.
func (c *cachedItem) monitor() {
	for {
		sleeptime := time.Until(c.session.expiry)
		time.Sleep(sleeptime)
		if c.session.expiry.Before(time.Now()) {
			return
		}
	}
}

//getSessionID returns the stored key.
func (c *cachedItem) getSessionID() string {
	return c.session.key
}

//UpdateExpiryTime updates the session's expiry time.
func (c *cachedItem) UpdateExpiryTime(updatedTime time.Time) {
	c.session.expiry = updatedTime
}

//itemCache is a wrapper for the default cache type for each Cache to be distinguishable and have internal methods if required.
type itemCache struct {
	cache
	sorted []db.Product
}

//UpdateSorted flushes the sorted array with the new data from the DB. It gets called when the cache is flushed.
func (i *itemCache) updateSorted() {
	preSort := make([]db.Product, 30)
	for _, as := range i.cache {
		preSort = append(preSort, (as.(*cachedItem)).item)
	}
	i.sort()
}

//update override the default update method for caches.
func (i *itemCache) update(payLoad ...ActiveSession) {
	flushedMap := make(map[string]ActiveSession)
	for _, activeSession := range payLoad {
		flushedMap[activeSession.getSessionID()] = activeSession
	}
	i.cache = flushedMap
	i.updateSorted()
}

//sort sorts the sorted item cache.
func (i *itemCache) sort() {
	sort.Slice(i.sorted, func(j, k int) bool { return i.sorted[j].Sales > i.sorted[k].Sales })
}

package cache

import (
	"GoLive/db"
	"fmt"
	"time"
)

const (
	SessionLife = 30 // default life of active session in minutes.
)

//cache type is a map that maps the key to an active session
type cache map[string]CacheObject // key: session value; value: see session definition below.

// Defines the values contained in type cache.
// CacheObjects can be of type cachedItem, MerchantSession or UserSession.
type CacheObject interface {
	monitor()
	updateExpiryTime(time.Time)
	getKey() string
}

//add takes in an active session, checks if the active session already exists in the cache.
//if it doesn't exist in the cache, add it into the cache and fire off a go tidy() go routine which regulates the session and removes it on expiry.
func (c *cache) add(payLoad CacheObject) {
	key := payLoad.getKey()
	if ok := c.check(key); !ok {
		(*c)[key] = payLoad
		fmt.Println((*c)[key], "added into the cache")
		fmt.Println(key, "key value during storage")
		go c.tidy(key, payLoad)
	}
}

func (c *cache) update(key string, productID string, operator string, toCache cachedItem, items *itemCache) {
	if c.check(key) {
		cacheObject := (*c)[key]
		switch v := cacheObject.(type) {
		case *UserSession:
			v.updateCart(productID, operator, toCache.item, items)
		case *MerchantSession:
			//c.activeUserCache.refresh(v)
		case *cachedItem:
			v.updateExpiryTime(time.Now().Add(SessionLife * time.Minute))
		}
	}
}

// //refresh flushes the cache with the data obtained from the DB.
// func (c *cache) refresh(payLoad ...ActiveSession) {
// 	flushedMap := make(map[string]ActiveSession)
// 	for _, activeSession := range payLoad {
// 		flushedMap[activeSession.getKey()] = activeSession
// 	}
// 	(*c) = flushedMap
// }

//get returns the CacheObject stored in the cache given the key.
func (c *cache) get(key string, userID string, database *db.Database) (CacheObject, error) {
	if retrieved, ok := ((*c)[key]); !ok {

		// Additem gets from the DB
		itemToCache, err := c.AddItemFromDB(key, userID, database)
		if err != nil {
			return itemToCache, err
		}

		//additemtocache
		(*c)[itemToCache.getKey()] = itemToCache
		go c.tidy(itemToCache.getKey(), itemToCache)
		return itemToCache, nil
	} else {
		retrieved.updateExpiryTime(time.Now().Add(SessionLife * time.Minute))
		((*c)[key]) = retrieved
		return retrieved, nil
	}

	// cacheObject, ok := (*c)[key]
	// fmt.Println((*c)[key], "getting from cache")
	// fmt.Println("key value during access", key)
	// return cacheObject, ok
}

func (c *cache) AddItemFromDB(key string, userID string, db *db.Database) (CacheObject, error) {
	switch string(key[0]) {
	case "U":
		user, err := db.GetUserFromID(userID)
		if err != nil {
			return &UserSession{}, err
		}
		return NewUserSession(key, time.Now().Add(SessionLife*time.Minute), user, NewCart()), nil

	case "M":
		merch, err := db.GetInventory(userID)
		if err != nil {
			return &MerchantSession{}, err
		}
		return NewMerchSession(key, time.Now().Add(SessionLife*time.Minute), merch.MerchantUser), nil

	default:
		//DB interaction with error returned return error here
		prod, err := db.GetProduct(key)
		if err != nil {
			return &cachedItem{}, err
		}
		cachedItem := cachedItem{
			item:   &prod,
			expiry: expiry{time.Now().Add(SessionLife * time.Minute)},
		}
		return &cachedItem, nil
	}
}

//tidy calls on monitor which checks if the session has expired. If the session has expired, monitor returns and the session is deleted from the cache.
func (c *cache) tidy(key string, obj CacheObject) {
	obj.monitor()
	delete(*c, key)
}

// Validate that session is active, and if it is active refresh the timestamp.
func (c *cache) check(key string) bool {
	cacheObject, ok := (*c)[key]
	if ok {
		cacheObject.updateExpiryTime(time.Now().Add(SessionLife * time.Minute))
	}
	return ok
}

//remove deletes the session given the ActiveSession.
func (c *cache) remove(key string) {
	delete(*c, key)
}

// NewUserSession takes in required inputs and returns a new UserSession.
func NewUserSession(uuid string, expiryTime time.Time, user db.User, cart Cart) *UserSession {
	return &UserSession{
		owner:   user,
		cart:    cart,
		session: session{uuid, expiry{expiryTime}},
	}
}

// NewMerchantSession takes in required inputs and returns a new MerchantSession
func NewMerchSession(uuid string, expiryTime time.Time, user db.MerchantUser) *MerchantSession {
	return &MerchantSession{
		owner:   user,
		session: session{uuid, expiry{expiryTime}},
	}
}

type SessionCopy interface {
	Data()
}

// //cachedMerchant stores the merchant details and products.
// type cachedMerchant struct {
// 	cacheObject
// 	// merchant db.MerchantUser
// 	// products *[]db.Product
// 	// hitRate  int
// }

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

// //sort sorts the sorted item cache.
// func (i *itemCache) sort() {
// 	sort.Slice(i.sorted, func(j, k int) bool { return i.sorted[j].Sales > i.sorted[k].Sales })
// }

//func (c *CacheManager) CheckItemPresence(prodID string) (cachedItem, bool) {
//	cachedItem, ok := (*c).itemsCache.CheckItem(prodID)
//	if !ok {
//		//query db
//
//		(*c).itemsCache.AddItemFromDB(cachedItem)
//	}
//	return cachedItem, ok
//}

package cache

import (
	"GoLive/db"
	"time"
)

const (
	sessionLife = 30 // default life of active session in minutes
)

type CacheManager struct {
	activeUserCache cache // key: session value; value: see session definition below
	merchantCache   cache // key: Pointer to cache; value: access frequency
	featuredItems   cache // Populates home page
}

//
func (c *CacheManager) AddtoCache(payLoad interface{}) {
	switch v := payLoad.(type) {
	case userSession:
		c.activeUserCache.add(v)
	case merchantSession:
		c.merchantCache.add(v)
	case db.Product:
		c.featuredItems.add(v)
	}
}

type cache interface {
	add(payload interface{})
	tidy()
	check()
}

type userSession struct { // Cart CRUD tied to methods on this type
	sessionID           string
	expiry              time.Time
	owner               db.User // Owner can be User or a Merchant
	merchantDescription *string // nil means user, some pointer means merchant
	cart                *[]db.Product
}

type userCache map[string]*userSession // key: session value; value: see session definition below
func (u *userCache) add(payLoad *userSession) {
	(*u)[payLoad.sessionID] = payLoad
	go u.tidy(payLoad)
}

func (u *userCache) tidy(session *userSession) {
	session.monitor()
	delete((*u), session.sessionID)
}

// Validate that session is active, and if active refresh the timestamp
func (u *userCache) check(sessionID string) bool {
	userSession, ok := (*u)[sessionID]
	if ok {
		userSession.expiry = time.Now().Add(sessionLife * time.Minute)
	}
	return ok
}

func (m *userSession) monitor() {
	for {
		sleeptime := m.expiry.Sub(time.Now())
		time.Sleep(sleeptime)
		if m.expiry.Before(time.Now()) {
			return
		}
	}
}

type merchantSession struct {
	sessionID string
	expiry    time.Time
	merchant  db.MerchantUser
	products  *[]db.Product
}

type merchantCache map[*merchantSession]int // key: Pointer to cache; value: access frequency
func (m *merchantCache) add() {

}
func (m *merchantCache) tidy() {
	// TODO Logic for tidying

}

func (m *merchantSession) monitor()

func (m *merchantCache) check() {

}

type FeaturedItemsSession struct {
	expiry time.Time
	item   db.Product
}

type featuredItemsCache []FeaturedItemsSession

func (f *featuredItemsCache) add() {

}
func (f *featuredItemsCache) tidy() {
	// TODO Filter most popular recent items based on sales log (we don't have a sales log as yet)
}
func (f *featuredItemsCache) check() {

}

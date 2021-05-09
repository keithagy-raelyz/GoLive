package app

import (
	"GoLive/cache"
	"GoLive/db"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (a *App) InitializeCacheManager() {
	a.cacheManager = cache.NewCacheManager()
	dummyUserSession := cache.NewUserSession(
		"6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
		time.Now().Add(30*time.Minute),
		db.User{
			Id:       "4",
			Name:     "DogShiet",
			Email:    "feecalmatter@hotmail.com",
			Password: "abc"},
		&[]cache.CartItem{})
	a.cacheManager.AddtoCache(dummyUserSession)
}

func (a *App) InitializeBlacklist() {
	a.blacklist = make(map[string]bool)
	//TODO undo blacklist here commenting out to test routes.
	//a.blacklist["/cart"] = true
	//a.blacklist["/checkout"] = true
	//a.blacklist["/users"] = true
}

// Helper function to check if login should be verified for a given request URL
func (a *App) NeedSessionCookie(r *http.Request) bool {
	url := r.URL.String()
	if url == "/users" && r.Method == http.MethodPost {
		return false
	}
	urlSplit := strings.Split(url, "/")
	url = "/" + urlSplit[1]
	fmt.Println("url:", url)

	_, ok := a.blacklist[url]

	return ok
}

// Validate session
func (a *App) Middleware(endPoint http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.NeedSessionCookie(r) {
			endPoint.ServeHTTP(w, r)
			return
		} else {
			if _, found := a.HaveValidSessionCookie(r); found {
				endPoint.ServeHTTP(w, r)
				return
			}
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

func (a *App) HaveValidSessionCookie(r *http.Request) (cache.ActiveSession, bool) {
	// Get session cookie
	sessionValue, err := r.Cookie("sessionCookie")
	if err != nil {
		// No Session Cookie
		return nil, false
	}
	sessionValStr := sessionValue.String()
	session, found := a.cacheManager.GetFromCache(sessionValStr, "activeUsers")
	if !found {
		session, found = a.cacheManager.GetFromCache(sessionValStr, "activeMerchants")
	}
	return session, found
}

// UpdateSession is an App method, to be called by HTTP handlers for the relevant cache manager to refresh sessions / update carts for the active user.

//TODO
// UpdateSession is an App method, to be called by HTTP handlers for the relevant cache manager to refresh sessions / update carts for the active user.
func (a *App) UpdateSession(activeSession cache.ActiveSession, cart *[]db.Product) {
	if cart == nil {
		// MerchantSession or UserSession page navigation, extend expiry by standard session life
		// a.cacheManager.UpdateCache()
		//activeSession.UpdateExpiryTim(time.Now().Add(cache.SessionLife * time.Minute))
	} else {
		// UserSession adding product to cart, update expiry by standard session life and update cart
		//activeSession.(*cache.UserSession).UpdateCart(cart)
		//activeSession.UpdateExpiryTime(time.Now().Add(cache.SessionLife * time.Minute))
	}
}

// DeleteSession is called upon logout
func (a *App) DeleteSession(session cache.ActiveSession) {

}

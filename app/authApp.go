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
		make([]cache.CartItem, 0))
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
		fmt.Println("no session cookie present in haveValidSessionCookie")
		return nil, false
	}
	sessionValStr := sessionValue.Value
	session, found := a.cacheManager.GetFromCache(sessionValStr, "activeUsers")
	if !found {
		fmt.Println("not found in cache manager")
		session, found = a.cacheManager.GetFromCache(sessionValStr, "activeMerchants")
	}
	return session, found
}

// UpdateSession is an App method, to be called by HTTP handlers for the relevant cache manager to refresh sessions / update carts for the active user.
func (a *App) UpdateSession(r *http.Request, cart []cache.CartItem) error {
	sessionValue, err := r.Cookie("sessionCookie")
	if err != nil {
		// No session cookie
		return err
	}
	sessionValStr := sessionValue.String()
	cacheType := ""
	switch string(sessionValStr[0]) {
	case "U":
		cacheType = "activeUsers"
	case "M":
		cacheType = "activeMerchants"
	}
	a.cacheManager.UpdateCache(sessionValStr, cacheType, cart)
	return nil
}

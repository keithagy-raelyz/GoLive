package app

import (
	"GoLive/cache"
	"fmt"
	"net/http"
	"strings"
)

func (a *App) InitializeCacheManager() {
	a.cacheManager = cache.NewCacheManager(a.connectDB())
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

func (a *App) HaveValidSessionCookie(r *http.Request) (cache.CacheObject, bool) {
	// Get session cookie
	sessionValue, err := r.Cookie("sessionCookie")
	if err != nil {
		// No Session Cookie
		fmt.Println("no session cookie present in haveValidSessionCookie")
		return nil, false
	}
	sessionValStr := sessionValue.Value
	fmt.Println(sessionValStr)
	session, err := a.cacheManager.GetFromCache(sessionValStr, "activeSession")
	if err != nil {
		fmt.Println("not found in cache manager:", nil)
		return session, false
	}
	return session, true
}

//// UpdateSession is an App method, to be called by HTTP handlers for the relevant cache manager to refresh sessions / update carts for the active user.
//func (a *App) UpdateSession(r *http.Request, cart []cache.CartItem) error {
//	sessionValue, err := r.Cookie("sessionCookie")
//	if err != nil {
//		// No session cookie
//		return err
//	}
//	sessionValStr := sessionValue.String()
//	cacheType := ""
//	switch string(sessionValStr[0]) {
//	case "U":
//		cacheType = "activeUsers"
//	case "M":
//		cacheType = "activeMerchants"
//	}
//	a.cacheManager.UpdateCart(sessionValStr, cacheType, cart)
//	return nil
//}

func (a *App) UpdateSession(r *http.Request) {

	//TODO uncomment this
	//sessionValue, err := r.Cookie("sessionCookie")
	//if err != nil {
	//	// No session cookie
	//	return err
	//}
	//sessionValStr := sessionValue.String()
	//cacheType := ""
	//switch string(sessionValStr[0]) {
	//case "U":
	//	cacheType = "activeUsers"
	//case "M":
	//	cacheType = "activeMerchants"
	//}
	//
	//return nil
}

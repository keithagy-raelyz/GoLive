package db

import (
	"fmt"
	"net/http"
	"strings"
)

//Auth internal management code
type Auth struct {
	session   map[string]string
	blacklist map[string]bool
}

func (d *Database) InitializeAndGetAuth() *Auth {
	d.a = &Auth{
		session:   make(map[string]string),
		blacklist: make(map[string]bool),
	}
	d.a.session["6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b"] = "74dcfbc208bb6aa08c90fb05bda0f2bc53285713e89611dfdd97ae129b5f6195"
	d.a.blacklist["/cart"] = true
	d.a.blacklist["/checkout"] = true
	d.a.blacklist["/users"] = true

	return d.a
}

func (a *Auth) InBlacklist(r *http.Request) bool {
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
func (a *Auth) Middleware(endPoint http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session cookie
		sessionCookie, err := r.Cookie("sessionCookie")
		if err != nil {
			// No Session Cookie
			if a.InBlacklist(r) {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			endPoint.ServeHTTP(w, r)
			return
		}
		sessionCookieVal := sessionCookie.Value

		secondTest, err := r.Cookie("i184m")
		if err != nil {
			// No Second Test
			if a.InBlacklist(r) {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			endPoint.ServeHTTP(w, r)
			return
		}
		secondTestVal := secondTest.Value

		if secondTestValue, ok := a.session[sessionCookieVal]; ok {
			if secondTestValue == secondTestVal {
				// Validated; continue
				endPoint.ServeHTTP(w, r)
				return
			} else {
				// Failed second check
				w.WriteHeader(http.StatusUnavailableForLegalReasons)
				w.Write([]byte("451 - :)"))
				return
			}
		} else {
			// Session Cookie is present but wrong/expired
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	})
}

func (a *Auth) ManageSession() {

}

func (a *Auth) UpdateSession() {

}

//DB Interaction code
func (d *Database) GetSessions() {

}

func (d *Database) GetSession(userSession string) {

}

func (d *Database) CreateSession() {

}

func (d *Database) UpdateSession() {

}
func (d *Database) DeleteSession(sessionID string) {
	//TODO consider delete to User called by random Curl requests ie no Authentication
	//res, err := d.b.Exec("DELETE FROM sessions where sessionID =? ", sessionID)
	//if err != nil {
	//	//TODO return custom error msg
	//	return err
	//}
	//rowCount, err := res.RowsAffected()
	//if err != nil || rowCount != 1 {
	//	//TODO return custom error msg
	//	return err
	//}
	//return nil
}

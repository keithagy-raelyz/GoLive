package app

import (
	"GoLive/cache"
	"GoLive/db"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type Data struct {
	User         db.User          //to display/edit individual user
	Merchant     db.Merchant      //to display/edit individual merchant profile to himself logged in merchant
	MerchantShop db.Merchant      //display/edit individual merchant shop to consumers public page
	Merchants    []db.Merchant    //to display/edit all the merchants
	Error        Error            //to display/edit an error message IF there is an error msg
	Products     []db.Product     //to display/edit featured items
	Cart         []cache.CartItem //to display/edit checkout cart
	JSON         string           // to display any FINALIZED data which will not undergo further changes (e.g cart at checkout page)
}
type Error struct {
	ErrMsg string
}

const (
	InvalidCredentials = "Invalid Username or Password"
	EmptyUsername      = "Username cannot be empty"
	UserExists         = "Username is in use"
	PasswordMismatch   = "Passwords are not the same"
	EmptyEmail         = "Email is required"
	DBError            = "Internal Server Error Please try again later"
)

func (a *App) displayLogin(w http.ResponseWriter, r *http.Request) {
	// Check session; if already logged in then redirect to home page
	if _, alreadyLoggedIn := a.HaveValidSessionCookie(r); alreadyLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// ExecuteTemplate
	data := Data{}
	parseLoginPage(&w, data)
}

func (a *App) validateUserLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}

	foundUser, err := a.db.GetUser(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		//w.Write([]byte("403 - Invalid Login Credentials"))
		data := Data{Error: Error{InvalidCredentials}}
		parseLoginPage(&w, data)
		return
	}
	if err := bcrypt.CompareHashAndPassword(hashedPW, []byte(foundUser.Password)); err != nil {
		w.WriteHeader(http.StatusForbidden)
		//w.Write([]byte("403 - Invalid Login Credentials"))
		data := Data{Error: Error{InvalidCredentials}}
		parseLoginPage(&w, data)
		return
	}
	newsessionKey := "U" + uuid.NewV4().String()
	newsession := cache.NewUserSession(
		newsessionKey,
		time.Now().Add(cache.SessionLife*time.Minute),
		foundUser,
		cache.NewCart(),
	)
	newCookie := &http.Cookie{
		Name:  "sessionCookie",
		Value: newsessionKey,
		Path:  "/",
	}

	http.SetCookie(w, newCookie)
	a.cacheManager.AddtoCache(newsession)
	p, _ := a.db.GetAllProducts()
	data := Data{User: foundUser, Products: p}
	parseHomePage(&w, data)
}

func (a *App) validateMerchantLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}

	foundMerch, err := a.db.GetMerchant(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		//w.Write([]byte("403 - Invalid Login Credentials"))
		return
	}
	if err := bcrypt.CompareHashAndPassword(hashedPW, []byte(foundMerch.User.Password)); err != nil {
		w.WriteHeader(http.StatusForbidden)
		//w.Write([]byte("403 - Invalid Login Credentials"))
		return
	}

	newsessionKey := "M" + uuid.NewV4().String()
	newsession := cache.NewMerchSession(newsessionKey,
		time.Now().Add(cache.SessionLife),
		foundMerch)
	a.cacheManager.AddtoCache(newsession)
}

func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	sessionKey, err := r.Cookie("sessionCookie")
	sessionKeyVal := sessionKey.String()
	// Verify valid user type
	if err != nil {
		// No Session Cookie
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	switch string(sessionKeyVal[0]) {
	case "U":
		a.cacheManager.RemoveFromCache(sessionKeyVal, "activeUsers")
	case "M":
		a.cacheManager.RemoveFromCache(sessionKeyVal, "activeMerchants")
	}
	sessionKey.Expires = time.Now()
	http.SetCookie(w, sessionKey)
	jData, _ := json.Marshal(Response{true})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)

}

func (a *App) allUser(w http.ResponseWriter, r *http.Request) {
	//TODO display all Users
	users, err := a.db.GetUsers()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(users)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Get All Users OK"))
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// Verify valid ID
	userID := params["userid"]
	var u db.User
	u, err := a.db.GetUser(userID)
	if err != nil {
		// Invalid user ID inputted
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No user for input USERID"))
		return
	}
	fmt.Println(u)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Valid merchant ID, displaying store"))
}

func (a *App) postUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	username := r.FormValue("username")
	email := strings.ToLower(r.FormValue("email"))
	pw1 := r.FormValue("pw1")
	pw2 := r.FormValue("pw2")

	var u db.User

	u.Name = username
	u.Email = email
	if u.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		data := Data{Error: Error{EmptyUsername}}
		parseLoginPage(&w, data)
		return
	}
	if u.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		data := Data{Error: Error{EmptyEmail}}
		parseLoginPage(&w, data)
		return
	}
	err = a.db.CheckUser(u)
	if err != nil {
		//send the err msg back (err = errmsg)
		w.WriteHeader(http.StatusBadRequest)
		data := Data{Error: Error{UserExists}}
		parseLoginPage(&w, data)
		return
	}
	if pw1 != pw2 {
		//t.ParseFiles("./templates/errorRegister.html")
		//data := Data{nil, ErrorMsg{"color:red", "Password entered are different"}, ErrorMsg{"display:none", ""}, ErrorMsg{"display:block", ""}, 0, nil}
		//t.Execute(w, data)
		w.WriteHeader(http.StatusBadRequest)
		data := Data{Error: Error{PasswordMismatch}}
		parseLoginPage(&w, data)
		return
	}
	err = a.db.CreateUser(u, pw1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		data := Data{Error: Error{DBError}}
		parseLoginPage(&w, data)
		return
	}

	//w.Write([]byte("201 - User Creation Successful"))
	newsessionKey := "U" + uuid.NewV4().String()
	newsession := cache.NewUserSession(
		newsessionKey,
		time.Now().Add(cache.SessionLife*time.Minute),
		u,
		cache.NewCart())
	newCookie := &http.Cookie{
		Name:  "sessionCookie",
		Value: newsessionKey,
	}
	http.SetCookie(w, newCookie)
	a.cacheManager.AddtoCache(newsession)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) putUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, ok := params["userid"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Status Not Found"))
		http.Redirect(w, r, "/users", http.StatusNotFound)
		return
	}

	username := r.URL.Query().Get("Username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Username cannot be empty"))
		return
	}

	email := r.URL.Query().Get("Email")
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Email cannot be empty"))
		return
	}
	var u db.User
	u.Name = username
	u.Email = email
	u.Id = userID
	err := a.db.UpdateUser(u)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Unprocessable Entity"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Updated Successfully"))
}

func (a *App) delUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, ok := params["userid"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Status not Found"))
		return
	}
	err := a.db.DeleteUser(userID)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Unprocessable Entity"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Deleted Successfully"))
}

func parseLoginPage(w *http.ResponseWriter, data Data) {
	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/loginBody.html", "templates/error.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(*w, data)
	if err != nil {
		log.Fatal(err)
	}
}

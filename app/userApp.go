package app

import (
	"GoLive/cache"
	"GoLive/db"
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
	User     db.User
	Merchant db.MerchantUser
	Error    Error
	Products []db.Product
}
type Error struct {
	ErrMsg string
}

func (a *App) displayLogin(w http.ResponseWriter, r *http.Request) {
	// Check session; if already logged in then redirect to home page
	if _, alreadyLoggedIn := a.HaveValidSessionCookie(r); alreadyLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// ExecuteTemplate
	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/loginBody.html", "templates/error.html")
	if err != nil {
		log.Fatal(err)
	}
	data := Data{}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
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
		fmt.Println(err, "error in get user")
		w.WriteHeader(http.StatusForbidden)
		//w.Write([]byte("403 - Invalid Login Credentials"))
		t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/loginBody.html", "templates/error.html")
		if err != nil {
			log.Fatal(err)
		}
		data := Data{Error: Error{"Piece of shit who can't even remember your password / username"}}
		err = t.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := bcrypt.CompareHashAndPassword(hashedPW, []byte(foundUser.Password)); err != nil {
		fmt.Println(foundUser.Password, "not correct")
		w.WriteHeader(http.StatusForbidden)
		//w.Write([]byte("403 - Invalid Login Credentials"))
		t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/loginBody.html", "templates/error.html")
		if err != nil {
			log.Fatal(err)
		}
		data := Data{Error: Error{"Piece of shit who can't even remember your password / username"}}
		err = t.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	newsessionKey := "U" + uuid.NewV4().String()
	newsession := cache.NewUserSession(
		newsessionKey,
		time.Now().Add(cache.SessionLife*time.Minute),
		foundUser,
		nil)
	newCookie := &http.Cookie{
		Name:  "sessionCookie",
		Value: newsessionKey,
	}
	r.AddCookie(newCookie)
	a.cacheManager.AddtoCache(newsession)

	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/body.html", "templates/error.html")
	if err != nil {
		log.Fatal(err)
	}
	p, _ := a.db.GetAllProducts()
	data := Data{User: foundUser, Products: p}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) allUser(w http.ResponseWriter, r *http.Request) {

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
	// Verify valid merchant ID
	userID := params["userid"]
	var u db.User
	u, err := a.db.GetUser(userID)
	if err != nil {
		fmt.Println(err)
		// Invalid user ID inputted
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No user for input USERID"))
		return
	}
	fmt.Println(u, "printing user")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Valid merchant ID, displaying store"))
}

func (a *App) postUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	email := strings.ToLower(r.FormValue("email"))
	pw1 := r.FormValue("pw1")
	pw2 := r.FormValue("pw2")

	var u db.User

	u.Name = username
	u.Email = email
	if u.Name == "" {
		fmt.Println("Empty user name")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}
	if u.Email == "" {
		fmt.Println("Empty email")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}
	err := a.db.CheckUser(u)
	if err != nil {
		//send the err msg back (err = errmsg)
		fmt.Println(err, "CheckUser error")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}
	if pw1 != pw2 {
		//t.ParseFiles("./templates/errorRegister.html")
		//data := Data{nil, ErrorMsg{"color:red", "Password entered are different"}, ErrorMsg{"display:none", ""}, ErrorMsg{"display:block", ""}, 0, nil}
		//t.Execute(w, data)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Repeated PW"))
		return
	}
	err = a.db.CreateUser(u, pw1)
	if err != nil {
		//send the err msg back (err = errmsg)
		fmt.Println(err, "createuser error")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	//w.Write([]byte("201 - User Creation Successful"))
	newsessionKey := "U" + uuid.NewV4().String()
	newsession := cache.NewUserSession(
		newsessionKey,
		time.Now().Add(cache.SessionLife*time.Minute),
		u,
		nil)
	newCookie := &http.Cookie{
		Name:  "sessionCookie",
		Value: newsessionKey,
	}
	r.AddCookie(newCookie)
	a.cacheManager.AddtoCache(newsession)
	http.Redirect(w, r, "/", http.StatusCreated)
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
